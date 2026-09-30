package service

import (
	"context"
	"fmt"

	"github.com/oa-portal/erp-cost-service/internal/event"
	"github.com/oa-portal/erp-cost-service/internal/model"
	"github.com/oa-portal/erp-cost-service/internal/repo"
	"gorm.io/gorm"
)

const (
	// 成本中心类型
	CCTypeDept    = 1 // 部门
	CCTypeWorkshop = 2 // 车间
	CCTypeOther   = 3 // 其他
	// 成本要素
	ElemMaterial = 1 // 物料
	ElemLabor    = 2 // 人工
	ElemOverhead = 3 // 制造费用
	// 成本记录状态
	StatusPending   = 1 // 待归集
	StatusCollected = 2 // 已归集
)

type CostService struct {
	db           *gorm.DB
	ccRepo       *repo.CostCenterRepo
	pcRepo       *repo.ProductCostRepo
	recordRepo   *repo.CostRecordRepo
	publisher    event.Publisher
}

func NewCostService(
	db *gorm.DB,
	ccRepo *repo.CostCenterRepo,
	pcRepo *repo.ProductCostRepo,
	recordRepo *repo.CostRecordRepo,
	pub event.Publisher,
) *CostService {
	return &CostService{
		db: db, ccRepo: ccRepo, pcRepo: pcRepo,
		recordRepo: recordRepo, publisher: pub,
	}
}

// ---------- 成本中心 ----------

func (s *CostService) CreateCostCenter(ctx context.Context, c *model.CostCenter) error {
	return s.ccRepo.Create(ctx, c)
}

func (s *CostService) ListCostCenters(ctx context.Context, tenantID string, typ int32, p repo.Page) (*repo.PageResult[model.CostCenter], error) {
	return s.ccRepo.List(ctx, tenantID, typ, p)
}

// ---------- 产品成本 ----------

// CalculateProductCost 手动录入成本并计算
func (s *CostService) CalculateProductCost(ctx context.Context, pc *model.ProductCost) error {
	if pc.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	pc.TotalCost = pc.MaterialCost + pc.LaborCost + pc.OverheadCost
	pc.UnitCost = pc.TotalCost / pc.Quantity
	return s.pcRepo.Create(ctx, pc)
}

func (s *CostService) ListProductCosts(ctx context.Context, tenantID string, productID int64, p repo.Page) (*repo.PageResult[model.ProductCost], error) {
	return s.pcRepo.List(ctx, tenantID, productID, p)
}

// CollectCost 按产品+期间归集：汇总未归集的成本记录，生成产品成本，标记记录为已归集
func (s *CostService) CollectCost(ctx context.Context, tenantID string, productID int64, period string, quantity float64) (*model.ProductCost, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}
	material, labor, overhead, err := s.pcRepo.SumUncollected(ctx, tenantID, productID, period)
	if err != nil {
		return nil, fmt.Errorf("sum uncollected: %w", err)
	}
	if material+labor+overhead <= 0 {
		return nil, fmt.Errorf("no pending cost records for product %d in period %s", productID, period)
	}

	total := material + labor + overhead
	pc := &model.ProductCost{
		TenantID:     tenantID,
		ProductID:    productID,
		MaterialCost: material,
		LaborCost:    labor,
		OverheadCost: overhead,
		TotalCost:    total,
		Quantity:     quantity,
		UnitCost:     total / quantity,
		Period:       period,
	}

	// 事务：创建产品成本 + 标记记录已归集
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(pc).Error; err != nil {
			return err
		}
		return tx.Model(&model.CostRecord{}).
			Where("tenant_id = ? AND product_id = ? AND period = ? AND status = ?", tenantID, productID, period, StatusPending).
			Update("status", StatusCollected).Error
	}); err != nil {
		return nil, err
	}

	s.publisher.Publish("cost-events", "", nil)
	return pc, nil
}

// ---------- 成本记录 ----------

func (s *CostService) AddCostRecord(ctx context.Context, r *model.CostRecord) error {
	if r.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if r.Element < 1 || r.Element > 3 {
		return fmt.Errorf("invalid element: %d", r.Element)
	}
	r.Status = StatusPending
	return s.recordRepo.Create(ctx, r)
}

func (s *CostService) ListCostRecords(ctx context.Context, tenantID string, status int32, p repo.Page) (*repo.PageResult[model.CostRecord], error) {
	return s.recordRepo.List(ctx, tenantID, status, p)
}
