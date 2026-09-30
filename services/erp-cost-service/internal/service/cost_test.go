package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/erp-cost-service/internal/event"
	"github.com/oa-portal/erp-cost-service/internal/model"
	"github.com/oa-portal/erp-cost-service/internal/repo"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*CostService, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.AutoMigrate(&model.CostCenter{}, &model.ProductCost{}, &model.CostRecord{})
	ccRepo := repo.NewCostCenterRepo(db)
	pcRepo := repo.NewProductCostRepo(db)
	recordRepo := repo.NewCostRecordRepo(db)
	svc := NewCostService(db, ccRepo, pcRepo, recordRepo, event.NoopPublisher{})
	return svc, db
}

// 1. 手动计算产品成本：总成本=物料+人工+制造费用，单位成本=总成本/数量
func TestCalculateProductCost(t *testing.T) {
	svc, _ := setupTest(t)
	pc := &model.ProductCost{
		TenantID: "t1", ProductID: 1,
		MaterialCost: 100, LaborCost: 50, OverheadCost: 30,
		Quantity: 10, Period: "2026-09",
	}
	if err := svc.CalculateProductCost(context.Background(), pc); err != nil {
		t.Fatalf("calculate: %v", err)
	}
	if pc.TotalCost != 180 {
		t.Errorf("total_cost = %v, want 180", pc.TotalCost)
	}
	if pc.UnitCost != 18 {
		t.Errorf("unit_cost = %v, want 18", pc.UnitCost)
	}
}

// 2. 数量为零报错
func TestCalculateZeroQuantity(t *testing.T) {
	svc, _ := setupTest(t)
	pc := &model.ProductCost{TenantID: "t1", ProductID: 1, Quantity: 0}
	if err := svc.CalculateProductCost(context.Background(), pc); err == nil {
		t.Fatal("expected error for zero quantity")
	}
}

// 3. 归集：添加3条成本记录（物料/人工/制造费用），归集后生成产品成本，记录标记已归集
func TestCollectCost(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	svc.AddCostRecord(ctx, &model.CostRecord{TenantID: "t1", ProductID: 1, Element: 1, Amount: 200, Period: "2026-09"})
	svc.AddCostRecord(ctx, &model.CostRecord{TenantID: "t1", ProductID: 1, Element: 2, Amount: 80, Period: "2026-09"})
	svc.AddCostRecord(ctx, &model.CostRecord{TenantID: "t1", ProductID: 1, Element: 3, Amount: 40, Period: "2026-09"})

	pc, err := svc.CollectCost(ctx, "t1", 1, "2026-09", 20)
	if err != nil {
		t.Fatalf("collect: %v", err)
	}
	if pc.TotalCost != 320 {
		t.Errorf("total_cost = %v, want 320", pc.TotalCost)
	}
	if pc.UnitCost != 16 {
		t.Errorf("unit_cost = %v, want 16", pc.UnitCost)
	}
	if pc.MaterialCost != 200 || pc.LaborCost != 80 || pc.OverheadCost != 40 {
		t.Errorf("cost breakdown wrong: %+v", pc)
	}

	// 验证记录已标记为已归集
	res, _ := svc.ListCostRecords(ctx, "t1", StatusCollected, repo.Page{Page: 1, PageSize: 10})
	if res.Total != 3 {
		t.Errorf("collected records = %d, want 3", res.Total)
	}
}

// 4. 归集时无待归集记录报错
func TestCollectNoRecords(t *testing.T) {
	svc, _ := setupTest(t)
	_, err := svc.CollectCost(context.Background(), "t1", 99, "2026-09", 10)
	if err == nil {
		t.Fatal("expected error when no pending records")
	}
}

// 5. 成本记录金额必须为正
func TestCostRecordInvalidAmount(t *testing.T) {
	svc, _ := setupTest(t)
	err := svc.AddCostRecord(context.Background(), &model.CostRecord{
		TenantID: "t1", ProductID: 1, Element: 1, Amount: -10, Period: "2026-09",
	})
	if err == nil {
		t.Fatal("expected error for negative amount")
	}
}
