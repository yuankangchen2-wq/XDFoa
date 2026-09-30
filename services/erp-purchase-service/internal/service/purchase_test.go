package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/purchase-service/internal/event"
	"github.com/oa-portal/purchase-service/internal/model"
	"github.com/oa-portal/purchase-service/internal/repo"
	"gorm.io/gorm"
)

// mockInventoryClient 模拟 inventory-service 入库
type mockInventoryClient struct {
	stockInFunc func(tenantID, skuCode string, warehouseID, qty int64, referenceID, batchNo string) error
}

func (m *mockInventoryClient) StockIn(tenantID, skuCode string, warehouseID, qty int64, referenceID, batchNo string) error {
	if m.stockInFunc != nil {
		return m.stockInFunc(tenantID, skuCode, warehouseID, qty, referenceID, batchNo)
	}
	return nil
}

func setupTest(t *testing.T) (*PurchaseService, *mockInventoryClient) {
	t.Helper()
	dbFile := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&model.PurchaseOrder{}, &model.PurchaseOrderItem{},
		&model.PurchaseReceipt{}, &model.PurchaseReceiptItem{})

	mockInv := &mockInventoryClient{}
	svc := NewPurchaseService(
		db,
		repo.NewPurchaseOrderRepo(db),
		repo.NewPurchaseReceiptRepo(db),
		mockInv,
		event.NoopPublisher{},
	)
	return svc, mockInv
}

// 测试创建采购订单
func TestCreateOrder(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	o := &model.PurchaseOrder{
		TenantID:     "t1",
		SupplierCode: "SUP001",
		WarehouseID:  1,
		Remark:       "测试采购",
		Items: []model.PurchaseOrderItem{
			{SkuCode: "SKU001", Quantity: 100, UnitPrice: 10.5},
			{SkuCode: "SKU002", Quantity: 50, UnitPrice: 20.0},
		},
	}

	if err := svc.CreateOrder(ctx, o); err != nil {
		t.Fatal(err)
	}
	if o.OrderNo == "" {
		t.Fatal("order_no should not be empty")
	}
	if o.TotalAmount != 100*10.5+50*20.0 {
		t.Fatalf("total amount mismatch: got %f", o.TotalAmount)
	}
	if o.Status != 1 {
		t.Fatalf("expected status=1, got %d", o.Status)
	}
}

// 测试审批订单
func TestApproveOrder(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	o := &model.PurchaseOrder{
		TenantID:     "t1",
		SupplierCode: "SUP001",
		WarehouseID:  1,
		Items: []model.PurchaseOrderItem{
			{SkuCode: "SKU001", Quantity: 100, UnitPrice: 10},
		},
	}
	svc.CreateOrder(ctx, o)

	if err := svc.ApproveOrder(ctx, o.ID); err != nil {
		t.Fatal(err)
	}
	updated, _ := svc.GetOrder(ctx, o.ID)
	if updated.Status != 2 {
		t.Fatalf("expected status=2, got %d", updated.Status)
	}
}

// 测试到货入库
func TestReceiveGoods(t *testing.T) {
	svc, mockInv := setupTest(t)
	ctx := context.Background()

	stockInCalled := false
	mockInv.stockInFunc = func(tenantID, skuCode string, warehouseID, qty int64, referenceID, batchNo string) error {
		stockInCalled = true
		if qty != 100 {
			t.Errorf("expected stockin qty=100, got %d", qty)
		}
		return nil
	}

	o := &model.PurchaseOrder{
		TenantID:     "t1",
		SupplierCode: "SUP001",
		WarehouseID:  1,
		Items: []model.PurchaseOrderItem{
			{SkuCode: "SKU001", Quantity: 100, UnitPrice: 10},
		},
	}
	svc.CreateOrder(ctx, o)
	svc.ApproveOrder(ctx, o.ID)

	// 到货
	rcp := &model.PurchaseReceipt{
		OrderID: o.ID,
		Items: []model.PurchaseReceiptItem{
			{SkuCode: "SKU001", Quantity: 100},
		},
	}
	if err := svc.ReceiveGoods(ctx, rcp); err != nil {
		t.Fatal(err)
	}

	if !stockInCalled {
		t.Fatal("inventory stockin should be called")
	}

	// 订单状态应为"全部到货"
	updated, _ := svc.GetOrder(ctx, o.ID)
	if updated.Status != 4 {
		t.Fatalf("expected status=4 (all received), got %d", updated.Status)
	}

	// 明细 received_qty 应为 100
	if updated.Items[0].ReceivedQty != 100 {
		t.Fatalf("expected received_qty=100, got %d", updated.Items[0].ReceivedQty)
	}
}

// 测试部分到货
func TestPartialReceive(t *testing.T) {
	svc, mockInv := setupTest(t)
	ctx := context.Background()
	mockInv.stockInFunc = func(tenantID, skuCode string, warehouseID, qty int64, referenceID, batchNo string) error {
		return nil
	}

	o := &model.PurchaseOrder{
		TenantID:     "t1",
		SupplierCode: "SUP001",
		WarehouseID:  1,
		Items: []model.PurchaseOrderItem{
			{SkuCode: "SKU001", Quantity: 100, UnitPrice: 10},
		},
	}
	svc.CreateOrder(ctx, o)
	svc.ApproveOrder(ctx, o.ID)

	// 到货 60（部分）
	rcp := &model.PurchaseReceipt{
		OrderID: o.ID,
		Items: []model.PurchaseReceiptItem{
			{SkuCode: "SKU001", Quantity: 60},
		},
	}
	svc.ReceiveGoods(ctx, rcp)

	updated, _ := svc.GetOrder(ctx, o.ID)
	if updated.Status != 3 {
		t.Fatalf("expected status=3 (partial), got %d", updated.Status)
	}
}

// 测试结算
func TestSettleOrder(t *testing.T) {
	svc, mockInv := setupTest(t)
	ctx := context.Background()
	mockInv.stockInFunc = func(tenantID, skuCode string, warehouseID, qty int64, referenceID, batchNo string) error {
		return nil
	}

	o := &model.PurchaseOrder{
		TenantID:     "t1",
		SupplierCode: "SUP001",
		WarehouseID:  1,
		Items: []model.PurchaseOrderItem{
			{SkuCode: "SKU001", Quantity: 100, UnitPrice: 10},
		},
	}
	svc.CreateOrder(ctx, o)
	svc.ApproveOrder(ctx, o.ID)
	svc.ReceiveGoods(ctx, &model.PurchaseReceipt{
		OrderID: o.ID,
		Items:   []model.PurchaseReceiptItem{{SkuCode: "SKU001", Quantity: 100}},
	})

	if err := svc.SettleOrder(ctx, o.ID); err != nil {
		t.Fatal(err)
	}
	updated, _ := svc.GetOrder(ctx, o.ID)
	if updated.Status != 5 {
		t.Fatalf("expected status=5 (settled), got %d", updated.Status)
	}
}
