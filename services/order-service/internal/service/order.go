package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oa-portal/order-service/internal/event"
	"github.com/oa-portal/order-service/internal/model"
	"github.com/oa-portal/order-service/internal/repo"
	"gorm.io/gorm"
)

// InventoryClient 库存服务客户端接口
type InventoryClient interface {
	AllocateStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error
	ReleaseStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error
	DeductStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error
}

type OrderService struct {
	db        *gorm.DB
	orderRepo *repo.OrderRepo
	logRepo   *repo.OrderLogRepo
	invClient InventoryClient
	publisher event.Publisher
}

func NewOrderService(
	db *gorm.DB,
	orderRepo *repo.OrderRepo,
	logRepo *repo.OrderLogRepo,
	invClient InventoryClient,
	pub event.Publisher,
) *OrderService {
	return &OrderService{
		db: db, orderRepo: orderRepo, logRepo: logRepo,
		invClient: invClient, publisher: pub,
	}
}

// CreateOrder 下单：预扣库存 + 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, o *model.Order) error {
	o.OrderNo = fmt.Sprintf("ORD%d", time.Now().UnixNano())
	o.Status = 1 // 待支付

	// 计算总金额
	var total float64
	for i := range o.Items {
		total += float64(o.Items[i].Quantity) * o.Items[i].UnitPrice
	}
	o.TotalAmount = total

	// 先预扣所有商品库存
	for _, item := range o.Items {
		if err := s.invClient.AllocateStock(o.TenantID, item.SkuCode, o.WarehouseID, item.Quantity, o.OrderNo); err != nil {
			return fmt.Errorf("allocate stock failed for sku=%s: %w", item.SkuCode, err)
		}
	}

	// 创建订单
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(o).Error
	}); err != nil {
		// 订单创建失败，回滚预扣
		for _, item := range o.Items {
			s.invClient.ReleaseStock(o.TenantID, item.SkuCode, o.WarehouseID, item.Quantity, o.OrderNo)
		}
		return err
	}

	s.publisher.Publish("order-events", o.OrderNo, event.OrderPaid{
		TenantID: o.TenantID, OrderNo: o.OrderNo, OrderID: o.ID, UserID: o.UserID,
	})
	return nil
}

// ListOrders 查询订单
func (s *OrderService) ListOrders(ctx context.Context, tenantID string, userID int64, status int32, p repo.Page) (*repo.PageResult[model.Order], error) {
	return s.orderRepo.List(ctx, tenantID, userID, status, p)
}

// GetOrder 获取订单详情
func (s *OrderService) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	return s.orderRepo.Get(ctx, id)
}

// PayOrder 支付：真实扣减库存
func (s *OrderService) PayOrder(ctx context.Context, id int64) error {
	o, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if o.Status != 1 {
		return fmt.Errorf("order status invalid, current=%d", o.Status)
	}

	// 真实扣减库存
	for _, item := range o.Items {
		if err := s.invClient.DeductStock(o.TenantID, item.SkuCode, o.WarehouseID, item.Quantity, o.OrderNo); err != nil {
			return fmt.Errorf("deduct stock failed for sku=%s: %w", item.SkuCode, err)
		}
	}

	// 更新状态
	if err := s.orderRepo.UpdateStatus(ctx, id, 2); err != nil {
		return err
	}
	s.writeLog(ctx, id, 1, 2, "支付成功")

	s.publisher.Publish("order-events", o.OrderNo, event.OrderPaid{
		TenantID: o.TenantID, OrderNo: o.OrderNo, OrderID: o.ID, UserID: o.UserID,
	})
	return nil
}

// CancelOrder 取消订单：释放预扣库存
func (s *OrderService) CancelOrder(ctx context.Context, id int64) error {
	o, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if o.Status != 1 {
		return fmt.Errorf("only unpaid order can be cancelled, current=%d", o.Status)
	}

	// 释放预扣库存
	for _, item := range o.Items {
		if err := s.invClient.ReleaseStock(o.TenantID, item.SkuCode, o.WarehouseID, item.Quantity, o.OrderNo); err != nil {
			fmt.Printf("WARN: release stock failed for sku=%s: %v\n", item.SkuCode, err)
		}
	}

	if err := s.orderRepo.UpdateStatus(ctx, id, 9); err != nil {
		return err
	}
	s.writeLog(ctx, id, o.Status, 9, "订单取消")

	s.publisher.Publish("order-events", o.OrderNo, event.OrderCancelled{
		TenantID: o.TenantID, OrderNo: o.OrderNo, OrderID: o.ID,
	})
	return nil
}

// ShipOrder 发货
func (s *OrderService) ShipOrder(ctx context.Context, id int64) error {
	o, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if o.Status != 2 {
		return fmt.Errorf("order not paid, current=%d", o.Status)
	}
	if err := s.orderRepo.UpdateStatus(ctx, id, 3); err != nil {
		return err
	}
	s.writeLog(ctx, id, 2, 3, "已发货")
	return nil
}

// CompleteOrder 确认收货
func (s *OrderService) CompleteOrder(ctx context.Context, id int64) error {
	o, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if o.Status != 3 {
		return fmt.Errorf("order not shipped, current=%d", o.Status)
	}
	if err := s.orderRepo.UpdateStatus(ctx, id, 4); err != nil {
		return err
	}
	s.writeLog(ctx, id, 3, 4, "已完成")
	return nil
}

func (s *OrderService) writeLog(ctx context.Context, orderID int64, from, to int32, remark string) {
	s.logRepo.Create(ctx, &model.OrderLog{
		OrderID: orderID, FromStatus: from, ToStatus: to, Remark: remark,
	})
}
