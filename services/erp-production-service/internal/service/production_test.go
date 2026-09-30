package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/production-service/internal/event"
	"github.com/oa-portal/production-service/internal/model"
	"github.com/oa-portal/production-service/internal/repo"
	"gorm.io/gorm"
)

type mockInventoryClient struct {
	deductCalls map[string]int64 // sku -> total deducted
	stockInQty  int64
}

func (m *mockInventoryClient) DeductStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	m.deductCalls[skuCode] += qty
	return nil
}

func (m *mockInventoryClient) StockIn(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	m.stockInQty += qty
	return nil
}

func setupTest(t *testing.T) (*ProductionService, *mockInventoryClient) {
	t.Helper()
	dbFile := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&model.Bom{}, &model.BomItem{}, &model.WorkOrder{},
		&model.MaterialIssue{}, &model.MaterialIssueItem{}, &model.ProductionReport{})

	mockInv := &mockInventoryClient{deductCalls: make(map[string]int64)}
	svc := NewProductionService(
		db,
		repo.NewBomRepo(db),
		repo.NewWorkOrderRepo(db),
		repo.NewMaterialIssueRepo(db),
		repo.NewProductionReportRepo(db),
		mockInv,
		event.NoopPublisher{},
	)
	return svc, mockInv
}

// 测试创建 BOM
func TestCreateBom(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	bom := &model.Bom{
		TenantID:   "t1",
		ProductSku: "PROD001",
		BomName:    "成品A",
		Items: []model.BomItem{
			{ComponentSku: "MAT001", Quantity: 2, Level: 1},
			{ComponentSku: "MAT002", Quantity: 3, Level: 1},
		},
	}
	if err := svc.CreateBom(ctx, bom); err != nil {
		t.Fatal(err)
	}
	if bom.ID == 0 {
		t.Fatal("bom id should not be 0")
	}
}

// 测试 MRP 计算
func TestMRPCalculate(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	bom := &model.Bom{
		TenantID: "t1", ProductSku: "PROD001", BomName: "成品A",
		Items: []model.BomItem{
			{ComponentSku: "MAT001", Quantity: 2, Level: 1},
			{ComponentSku: "MAT002", Quantity: 3, Level: 1},
		},
	}
	svc.CreateBom(ctx, bom)
	svc.EnableBom(ctx, bom.ID)

	// 生产 10 个成品
	req, err := svc.MRPCalculate(ctx, "t1", "PROD001", 10)
	if err != nil {
		t.Fatal(err)
	}
	if req["MAT001"] != 20 {
		t.Errorf("expected MAT001=20, got %d", req["MAT001"])
	}
	if req["MAT002"] != 30 {
		t.Errorf("expected MAT002=30, got %d", req["MAT002"])
	}
}

// 测试完整生产流程：创建工单 → 领料 → 报工
func TestFullProductionFlow(t *testing.T) {
	svc, mockInv := setupTest(t)
	ctx := context.Background()

	// 1. 创建并启用 BOM
	bom := &model.Bom{
		TenantID: "t1", ProductSku: "PROD001", BomName: "成品A",
		Items: []model.BomItem{
			{ComponentSku: "MAT001", Quantity: 2, Level: 1},
		},
	}
	svc.CreateBom(ctx, bom)
	svc.EnableBom(ctx, bom.ID)

	// 2. 创建生产工单（生产 5 个成品）
	wo := &model.WorkOrder{
		TenantID: "t1", ProductSku: "PROD001", Quantity: 5, WarehouseID: 1,
	}
	if err := svc.CreateWorkOrder(ctx, wo); err != nil {
		t.Fatal(err)
	}
	if wo.Status != 1 {
		t.Fatalf("expected status=1, got %d", wo.Status)
	}

	// 3. 领料（应扣减 MAT001 × 10）
	if _, err := svc.IssueMaterial(ctx, wo.ID); err != nil {
		t.Fatal(err)
	}
	if mockInv.deductCalls["MAT001"] != 10 {
		t.Errorf("expected deduct MAT001=10, got %d", mockInv.deductCalls["MAT001"])
	}

	// 工单状态应为"生产中"
	updated, _ := svc.GetWorkOrder(ctx, wo.ID)
	if updated.Status != 2 {
		t.Fatalf("expected status=2 (production), got %d", updated.Status)
	}

	// 4. 报工（完工 5 个）
	if err := svc.ReportProduction(ctx, wo.ID, 5); err != nil {
		t.Fatal(err)
	}
	if mockInv.stockInQty != 5 {
		t.Errorf("expected stockin qty=5, got %d", mockInv.stockInQty)
	}

	// 工单状态应为"已完工"
	final, _ := svc.GetWorkOrder(ctx, wo.ID)
	if final.Status != 3 {
		t.Fatalf("expected status=3 (completed), got %d", final.Status)
	}
	if final.ProducedQty != 5 {
		t.Fatalf("expected produced_qty=5, got %d", final.ProducedQty)
	}
}

// 测试领料前必须有已启用的 BOM
func TestIssueMaterialWithoutBom(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	wo := &model.WorkOrder{
		TenantID: "t1", ProductSku: "NO_BOM", Quantity: 5, WarehouseID: 1,
	}
	svc.CreateWorkOrder(ctx, wo)

	_, err := svc.IssueMaterial(ctx, wo.ID)
	if err == nil {
		t.Fatal("expected error when BOM not found")
	}
}
