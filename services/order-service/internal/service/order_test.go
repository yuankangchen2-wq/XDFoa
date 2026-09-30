package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/order-service/internal/event"
	"github.com/oa-portal/order-service/internal/model"
	"github.com/oa-portal/order-service/internal/repo"
	"gorm.io/gorm"
)

type mockInventoryClient struct {
	allocateCalls map[string]int64
	releaseCalls  map[string]int64
	deductCalls   map[string]int64
}

func (m *mockInventoryClient) AllocateStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	m.allocateCalls[skuCode] += qty
	return nil
}

func (m *mockInventoryClient) ReleaseStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	m.releaseCalls[skuCode] += qty
	return nil
}

func (m *mockInventoryClient) DeductStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	m.deductCalls[skuCode] += qty
	return nil
}

func setupTest(t *testing.T) (*OrderService, *mockInventoryClient) {
	t.Helper()
	dbFile := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&model.Order{}, &model.OrderItem{}, &model.OrderLog{})

	mockInv := &mockInventoryClient{
		allocateCalls: make(map[string]int64),
		releaseCalls:  make(map[string]int64),
		deductCalls:   make(map[string]int64),
	}
	svc := NewOrderService(
		db,
		repo.NewOrderRepo(db),
		repo.NewOrderLogRepo(db),
		mockInv,
		event.NoopPublisher{},
	)
	return svc, mockInv
}

// 测试下单（预扣库存）
func TestCreateOrder(t *testing.T) {
	svc, mockInv := setupTest(t)
	ctx := context.Background()

	o := &model.Order{
		TenantID:    "t1",
		UserID:      1001,
		WarehouseID: 1,
		Items: []model.OrderItem{
			{SkuCode: "SKU001", Quantity: 2, UnitPrice: 100},
			{SkuCode: "SKU002", Quantity: 3, UnitPrice: 50},
		},
	}

	if err := svc.CreateOrder(ctx, o); err != nil {
		t.Fatal(err)
	}
	if o.OrderNo == "" {
		t.Fatal("order_no should not be empty")
	}
	if o.TotalAmount != 2*100+3*50 {
		t.Fatalf("total amount mismatch: got %f", o.TotalAmount)
	}
	if o.Status != 1 {
		t.Fatalf("expected status=1, got %d", o.Status)
	}
	// 验证预扣被调用
	if mockInv.allocateCalls["SKU001"] != 2 {
		t.Errorf("expected allocate SKU001=2, got %d", mockInv.allocateCalls["SKU001"])
	}
	if mockInv.allocateCalls["SKU002"] != 3 {
		t.Errorf("expected allocate SKU002=3, got %d", mockInv.allocateCalls["SKU002"])
	}
}

// 测试支付（真实扣减）
func TestPayOrder(t *testing.T) {
	svc, mockInv := setupTest(t)
	ctx := context.Background()

	o := &model.Order{
		TenantID: "t1", UserID: 1001, WarehouseID: 1,
		Items: []model.OrderItem{{SkuCode: "SKU001", Quantity: 5, UnitPrice: 10}},
	}
	svc.CreateOrder(ctx, o)

	if err := svc.PayOrder(ctx, o.ID); err != nil {
		t.Fatal(err)
	}
	updated, _ := svc.GetOrder(ctx, o.ID)
	if updated.Status != 2 {
		t.Fatalf("expected status=2, got %d", updated.Status)
	}
	// 验证扣减被调用
	if mockInv.deductCalls["SKU001"] != 5 {
		t.Errorf("expected deduct SKU001=5, got %d", mockInv.deductCalls["SKU001"])
	}
}

// 测试取消订单（释放预扣）
func TestCancelOrder(t *testing.T) {
	svc, mockInv := setupTest(t)
	ctx := context.Background()

	o := &model.Order{
		TenantID: "t1", UserID: 1001, WarehouseID: 1,
		Items: []model.OrderItem{{SkuCode: "SKU001", Quantity: 5, UnitPrice: 10}},
	}
	svc.CreateOrder(ctx, o)

	if err := svc.CancelOrder(ctx, o.ID); err != nil {
		t.Fatal(err)
	}
	updated, _ := svc.GetOrder(ctx, o.ID)
	if updated.Status != 9 {
		t.Fatalf("expected status=9, got %d", updated.Status)
	}
	// 验证释放被调用
	if mockInv.releaseCalls["SKU001"] != 5 {
		t.Errorf("expected release SKU001=5, got %d", mockInv.releaseCalls["SKU001"])
	}
}

// 测试完整订单流程：下单 → 支付 → 发货 → 完成
func TestFullOrderFlow(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	o := &model.Order{
		TenantID: "t1", UserID: 1001, WarehouseID: 1,
		Items: []model.OrderItem{{SkuCode: "SKU001", Quantity: 1, UnitPrice: 99}},
	}
	svc.CreateOrder(ctx, o)

	// 支付
	if err := svc.PayOrder(ctx, o.ID); err != nil {
		t.Fatalf("pay failed: %v", err)
	}
	// 发货
	if err := svc.ShipOrder(ctx, o.ID); err != nil {
		t.Fatalf("ship failed: %v", err)
	}
	// 完成
	if err := svc.CompleteOrder(ctx, o.ID); err != nil {
		t.Fatalf("complete failed: %v", err)
	}

	final, _ := svc.GetOrder(ctx, o.ID)
	if final.Status != 4 {
		t.Fatalf("expected status=4, got %d", final.Status)
	}
}

// 测试状态流转校验：已支付订单不能取消
func TestCancelPaidOrder(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	o := &model.Order{
		TenantID: "t1", UserID: 1001, WarehouseID: 1,
		Items: []model.OrderItem{{SkuCode: "SKU001", Quantity: 1, UnitPrice: 10}},
	}
	svc.CreateOrder(ctx, o)
	svc.PayOrder(ctx, o.ID)

	err := svc.CancelOrder(ctx, o.ID)
	if err == nil {
		t.Fatal("expected error when cancelling paid order")
	}
}
