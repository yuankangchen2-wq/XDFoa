package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oa-portal/production-service/internal/event"
	"github.com/oa-portal/production-service/internal/model"
	"github.com/oa-portal/production-service/internal/repo"
	"gorm.io/gorm"
)

// InventoryClient 库存服务客户端接口
type InventoryClient interface {
	DeductStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error
	StockIn(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error
}

type ProductionService struct {
	db        *gorm.DB
	bomRepo   *repo.BomRepo
	woRepo    *repo.WorkOrderRepo
	issueRepo *repo.MaterialIssueRepo
	reportRepo *repo.ProductionReportRepo
	invClient InventoryClient
	publisher event.Publisher
}

func NewProductionService(
	db *gorm.DB,
	bomRepo *repo.BomRepo,
	woRepo *repo.WorkOrderRepo,
	issueRepo *repo.MaterialIssueRepo,
	reportRepo *repo.ProductionReportRepo,
	invClient InventoryClient,
	pub event.Publisher,
) *ProductionService {
	return &ProductionService{
		db: db, bomRepo: bomRepo, woRepo: woRepo, issueRepo: issueRepo,
		reportRepo: reportRepo, invClient: invClient, publisher: pub,
	}
}

// ---------- BOM ----------

func (s *ProductionService) CreateBom(ctx context.Context, b *model.Bom) error {
	b.Status = 1
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(b).Error
	})
}

func (s *ProductionService) ListBoms(ctx context.Context, tenantID string, p repo.Page) (*repo.PageResult[model.Bom], error) {
	return s.bomRepo.List(ctx, tenantID, p)
}

func (s *ProductionService) EnableBom(ctx context.Context, id int64) error {
	return s.bomRepo.UpdateStatus(ctx, id, 2)
}

// ---------- MRP 计算（BOM 展开）----------

// MRPCalculate 根据成品 SKU 和数量，展开 BOM 计算所需物料
// 返回：map[子件SKU]所需数量
func (s *ProductionService) MRPCalculate(ctx context.Context, tenantID, productSku string, qty int64) (map[string]int64, error) {
	bom, err := s.bomRepo.GetByProductSku(ctx, tenantID, productSku)
	if err != nil {
		return nil, fmt.Errorf("BOM not found for sku=%s: %w", productSku, err)
	}

	requirements := make(map[string]int64)
	for _, item := range bom.Items {
		requirements[item.ComponentSku] += item.Quantity * qty
	}
	return requirements, nil
}

// ---------- 生产工单 ----------

func (s *ProductionService) CreateWorkOrder(ctx context.Context, w *model.WorkOrder) error {
	w.OrderNo = fmt.Sprintf("WO%d", time.Now().UnixNano())
	w.Status = 1 // 待领料

	if err := s.woRepo.Create(ctx, w); err != nil {
		return err
	}

	s.publisher.Publish("production-events", w.OrderNo, event.WorkOrderCreated{
		TenantID: w.TenantID, OrderNo: w.OrderNo,
		ProductSku: w.ProductSku, Quantity: w.Quantity,
	})
	return nil
}

func (s *ProductionService) ListWorkOrders(ctx context.Context, tenantID string, status int32, p repo.Page) (*repo.PageResult[model.WorkOrder], error) {
	return s.woRepo.List(ctx, tenantID, status, p)
}

func (s *ProductionService) GetWorkOrder(ctx context.Context, id int64) (*model.WorkOrder, error) {
	return s.woRepo.Get(ctx, id)
}

// ---------- 领料 ----------

// IssueMaterial 领料：展开 BOM 计算所需物料，调用 inventory 扣减
func (s *ProductionService) IssueMaterial(ctx context.Context, workOrderID int64) (*model.MaterialIssue, error) {
	wo, err := s.woRepo.Get(ctx, workOrderID)
	if err != nil {
		return nil, err
	}
	if wo.Status != 1 {
		return nil, fmt.Errorf("work order status invalid, current=%d", wo.Status)
	}

	// MRP 展开计算物料需求
	requirements, err := s.MRPCalculate(ctx, wo.TenantID, wo.ProductSku, wo.Quantity)
	if err != nil {
		return nil, err
	}

	// 构建领料单
	issue := &model.MaterialIssue{
		TenantID:    wo.TenantID,
		WorkOrderID: wo.ID,
		IssueNo:     fmt.Sprintf("MI%d", time.Now().UnixNano()),
		Status:      1,
	}
	for sku, qty := range requirements {
		issue.Items = append(issue.Items, model.MaterialIssueItem{
			SkuCode: sku, Quantity: qty,
		})
	}

	// 保存领料单
	if err := s.issueRepo.Create(ctx, issue); err != nil {
		return nil, err
	}

	// 调用 inventory 扣减物料库存
	for _, item := range issue.Items {
		if err := s.invClient.DeductStock(wo.TenantID, item.SkuCode, wo.WarehouseID, item.Quantity, issue.IssueNo); err != nil {
			fmt.Printf("WARN: deduct stock failed for sku=%s: %v\n", item.SkuCode, err)
		}
	}

	// 更新领料单状态和工单状态
	s.issueRepo.UpdateStatus(ctx, issue.ID, 2)
	s.woRepo.UpdateStatus(ctx, wo.ID, 2) // 生产中

	s.publisher.Publish("production-events", issue.IssueNo, event.MaterialIssued{
		TenantID: wo.TenantID, IssueNo: issue.IssueNo, WorkOrderID: wo.ID,
	})

	return issue, nil
}

// ---------- 报工（完工入库）----------

// ReportProduction 报工：完工数量入库，更新工单状态
func (s *ProductionService) ReportProduction(ctx context.Context, workOrderID, producedQty int64) error {
	wo, err := s.woRepo.Get(ctx, workOrderID)
	if err != nil {
		return err
	}
	if wo.Status != 2 {
		return fmt.Errorf("work order not in production, status=%d", wo.Status)
	}
	if producedQty <= 0 {
		return fmt.Errorf("produced qty must be positive")
	}

	// 调用 inventory 完工入库
	if err := s.invClient.StockIn(wo.TenantID, wo.ProductSku, wo.WarehouseID, producedQty, wo.OrderNo); err != nil {
		return fmt.Errorf("stockin failed: %w", err)
	}

	// 更新工单已生产数量
	if err := s.woRepo.AddProducedQty(ctx, wo.ID, producedQty); err != nil {
		return err
	}

	// 保存报工记录
	report := &model.ProductionReport{
		TenantID: wo.TenantID, WorkOrderID: wo.ID,
		ProducedQty: producedQty, ReportNo: fmt.Sprintf("RP%d", time.Now().UnixNano()),
	}
	if err := s.reportRepo.Create(ctx, report); err != nil {
		return err
	}

	// 如果已生产数量 >= 计划数量，标记为完工
	updated, _ := s.woRepo.Get(ctx, wo.ID)
	if updated.ProducedQty >= wo.Quantity {
		s.woRepo.UpdateStatus(ctx, wo.ID, 3) // 已完工
		s.publisher.Publish("production-events", wo.OrderNo, event.ProductionCompleted{
			TenantID: wo.TenantID, OrderNo: wo.OrderNo,
			ProductSku: wo.ProductSku, Quantity: wo.Quantity,
		})
	}

	return nil
}

func (s *ProductionService) ListMaterialIssues(ctx context.Context, tenantID string, p repo.Page) (*repo.PageResult[model.MaterialIssue], error) {
	return s.issueRepo.List(ctx, tenantID, p)
}
