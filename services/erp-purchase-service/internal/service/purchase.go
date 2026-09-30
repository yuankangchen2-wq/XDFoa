package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oa-portal/purchase-service/internal/event"
	"github.com/oa-portal/purchase-service/internal/model"
	"github.com/oa-portal/purchase-service/internal/repo"
	"gorm.io/gorm"
)

type PurchaseService struct {
	db        *gorm.DB
	orderRepo *repo.PurchaseOrderRepo
	rcptRepo  *repo.PurchaseReceiptRepo
	invClient InventoryClient
	publisher event.Publisher
}

// InventoryClient 库存服务客户端接口（便于 mock）
type InventoryClient interface {
	StockIn(tenantID, skuCode string, warehouseID, qty int64, referenceID, batchNo string) error
}

func NewPurchaseService(
	db *gorm.DB,
	orderRepo *repo.PurchaseOrderRepo,
	rcptRepo *repo.PurchaseReceiptRepo,
	invClient InventoryClient,
	pub event.Publisher,
) *PurchaseService {
	return &PurchaseService{
		db: db, orderRepo: orderRepo, rcptRepo: rcptRepo,
		invClient: invClient, publisher: pub,
	}
}

// CreateOrder 创建采购订单
func (s *PurchaseService) CreateOrder(ctx context.Context, o *model.PurchaseOrder) error {
	o.OrderNo = fmt.Sprintf("PO%d", time.Now().UnixNano())
	o.Status = 1 // 草稿

	// 计算总金额
	var total float64
	for i := range o.Items {
		total += float64(o.Items[i].Quantity) * o.Items[i].UnitPrice
	}
	o.TotalAmount = total

	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	// 发布事件
	s.publisher.Publish("purchase-events", o.OrderNo, event.PurchaseOrderCreated{
		TenantID: o.TenantID, OrderNo: o.OrderNo,
		SupplierCode: o.SupplierCode, WarehouseID: o.WarehouseID,
		TotalAmount: o.TotalAmount,
	})
	return nil
}

// ListOrders 查询采购订单
func (s *PurchaseService) ListOrders(ctx context.Context, tenantID string, status int32, p repo.Page) (*repo.PageResult[model.PurchaseOrder], error) {
	return s.orderRepo.List(ctx, tenantID, status, p)
}

// GetOrder 获取订单详情
func (s *PurchaseService) GetOrder(ctx context.Context, id int64) (*model.PurchaseOrder, error) {
	return s.orderRepo.Get(ctx, id)
}

// ApproveOrder 审批订单
func (s *PurchaseService) ApproveOrder(ctx context.Context, id int64) error {
	o, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if o.Status != 1 {
		return fmt.Errorf("order status invalid, current=%d", o.Status)
	}
	return s.orderRepo.UpdateStatus(ctx, id, 2)
}

// ReceiveGoods 到货并入库
func (s *PurchaseService) ReceiveGoods(ctx context.Context, rcp *model.PurchaseReceipt) error {
	// 校验订单
	o, err := s.orderRepo.Get(ctx, rcp.OrderID)
	if err != nil {
		return err
	}
	if o.Status != 2 && o.Status != 3 {
		return fmt.Errorf("order not approved, status=%d", o.Status)
	}

	rcp.TenantID = o.TenantID
	rcp.ReceiptNo = fmt.Sprintf("RC%d", time.Now().UnixNano())
	rcp.Status = 1 // 待入库

	// 事务：创建到货单 + 更新订单明细已收数量
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rcp).Error; err != nil {
			return err
		}
		// 更新订单明细 received_qty
		for _, item := range rcp.Items {
			if err := tx.Model(&model.PurchaseOrderItem{}).
				Where("order_id = ? AND sku_code = ?", rcp.OrderID, item.SkuCode).
				UpdateColumn("received_qty", gorm.Expr("received_qty + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// 调用 inventory-service 入库（同步）
	for _, item := range rcp.Items {
		if err := s.invClient.StockIn(o.TenantID, item.SkuCode, o.WarehouseID, item.Quantity, rcp.ReceiptNo, ""); err != nil {
			// 入库失败不回滚到货单（实际生产中应有补偿机制），记录日志即可
			fmt.Printf("WARN: stockin failed for sku=%s: %v\n", item.SkuCode, err)
		}
	}

	// 更新到货单状态为已入库
	s.rcptRepo.UpdateStatus(ctx, rcp.ID, 2)

	// 更新订单状态：检查是否全部到货
	s.updateOrderStatusAfterReceipt(ctx, o)

	// 发布到货事件
	for _, item := range rcp.Items {
		s.publisher.Publish("purchase-events", item.SkuCode, event.GoodsReceived{
			TenantID: o.TenantID, ReceiptNo: rcp.ReceiptNo, OrderID: o.ID,
			WarehouseID: o.WarehouseID, SkuCode: item.SkuCode, Quantity: item.Quantity,
		})
	}

	return nil
}

// updateOrderStatusAfterReceipt 到货后更新订单状态
func (s *PurchaseService) updateOrderStatusAfterReceipt(ctx context.Context, o *model.PurchaseOrder) {
	// 重新加载订单明细
	updated, err := s.orderRepo.Get(ctx, o.ID)
	if err != nil {
		return
	}
	allReceived := true
	partial := false
	for _, item := range updated.Items {
		if item.ReceivedQty < item.Quantity {
			allReceived = false
		}
		if item.ReceivedQty > 0 {
			partial = true
		}
	}
	if allReceived {
		s.orderRepo.UpdateStatus(ctx, o.ID, 4) // 全部到货
	} else if partial {
		s.orderRepo.UpdateStatus(ctx, o.ID, 3) // 部分到货
	}
}

// SettleOrder 结算订单
func (s *PurchaseService) SettleOrder(ctx context.Context, id int64) error {
	o, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if o.Status != 4 {
		return fmt.Errorf("order not fully received, status=%d", o.Status)
	}
	if err := s.orderRepo.UpdateStatus(ctx, id, 5); err != nil {
		return err
	}
	s.publisher.Publish("purchase-events", o.OrderNo, event.PurchaseOrderCreated{
		TenantID: o.TenantID, OrderNo: o.OrderNo,
	})
	return nil
}

// ListReceipts 查询到货单
func (s *PurchaseService) ListReceipts(ctx context.Context, tenantID string, p repo.Page) (*repo.PageResult[model.PurchaseReceipt], error) {
	return s.rcptRepo.List(ctx, tenantID, p)
}
