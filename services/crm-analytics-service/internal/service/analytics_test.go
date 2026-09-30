package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/crm-analytics-service/internal/event"
	"github.com/oa-portal/crm-analytics-service/internal/model"
	"github.com/oa-portal/crm-analytics-service/internal/repo"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*AnalyticsService, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.AutoMigrate(&model.SalesSnapshot{}, &model.FunnelStage{}, &model.CustomerStat{})
	salesRepo := repo.NewSalesSnapshotRepo(db)
	funnelRepo := repo.NewFunnelStageRepo(db)
	custRepo := repo.NewCustomerStatRepo(db)
	svc := NewAnalyticsService(salesRepo, funnelRepo, custRepo, event.NoopPublisher{})
	return svc, db
}

// 1. 录入销售快照：客单价自动计算
func TestRecordSalesSnapshot(t *testing.T) {
	svc, _ := setupTest(t)
	o := &model.SalesSnapshot{
		TenantID: "t1", Period: "2026-09", ProductID: 1, CustomerID: 1,
		TotalAmount: 12000, OrderCount: 3,
	}
	if err := svc.RecordSalesSnapshot(context.Background(), o); err != nil {
		t.Fatalf("record: %v", err)
	}
	if o.AvgAmount != 4000 {
		t.Errorf("avg_amount = %v, want 4000", o.AvgAmount)
	}
}

// 2. 按期间汇总
func TestSalesSummary(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()
	svc.RecordSalesSnapshot(ctx, &model.SalesSnapshot{TenantID: "t1", Period: "2026-08", TotalAmount: 5000, OrderCount: 2})
	svc.RecordSalesSnapshot(ctx, &model.SalesSnapshot{TenantID: "t1", Period: "2026-09", TotalAmount: 3000, OrderCount: 1})
	svc.RecordSalesSnapshot(ctx, &model.SalesSnapshot{TenantID: "t1", Period: "2026-09", TotalAmount: 7000, OrderCount: 2})

	rows, err := svc.SalesSummaryByPeriod(ctx, "t1")
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("got %d periods, want 2", len(rows))
	}
	// 2026-09 汇总 = 10000, 3单
	for _, r := range rows {
		if r.Period == "2026-09" {
			if r.TotalAmount != 10000 {
				t.Errorf("2026-09 total = %v, want 10000", r.TotalAmount)
			}
			if r.OrderCount != 3 {
				t.Errorf("2026-09 orders = %d, want 3", r.OrderCount)
			}
		}
	}
}

// 3. 看板：销售额、订单数、转化率计算
func TestDashboard(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()
	svc.RecordSalesSnapshot(ctx, &model.SalesSnapshot{TenantID: "t1", Period: "2026-09", TotalAmount: 10000, OrderCount: 5})
	svc.RecordFunnelStage(ctx, &model.FunnelStage{TenantID: "t1", Period: "2026-09", Stage: "线索", OpportunityCount: 100})
	svc.RecordFunnelStage(ctx, &model.FunnelStage{TenantID: "t1", Period: "2026-09", Stage: "成交", OpportunityCount: 20})
	svc.RecordCustomerStat(ctx, &model.CustomerStat{TenantID: "t1", Period: "2026-09", NewCustomers: 15})

	d, err := svc.GetDashboard(ctx, "t1")
	if err != nil {
		t.Fatalf("dashboard: %v", err)
	}
	if d.TotalSales != 10000 {
		t.Errorf("total_sales = %v, want 10000", d.TotalSales)
	}
	if d.TotalOrders != 5 {
		t.Errorf("total_orders = %d, want 5", d.TotalOrders)
	}
	if d.AvgTicket != 2000 {
		t.Errorf("avg_ticket = %v, want 2000", d.AvgTicket)
	}
	if d.NewCustomers != 15 {
		t.Errorf("new_customers = %d, want 15", d.NewCustomers)
	}
	// 转化率 = 20/100 = 20%
	if d.ConversionRate < 19.9 || d.ConversionRate > 20.1 {
		t.Errorf("conversion_rate = %v, want ~20", d.ConversionRate)
	}
}

// 4. 漏斗列表
func TestFunnelList(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()
	svc.RecordFunnelStage(ctx, &model.FunnelStage{TenantID: "t1", Period: "2026-09", Stage: "线索", OpportunityCount: 100})
	svc.RecordFunnelStage(ctx, &model.FunnelStage{TenantID: "t1", Period: "2026-09", Stage: "商机", OpportunityCount: 60})
	svc.RecordFunnelStage(ctx, &model.FunnelStage{TenantID: "t1", Period: "2026-09", Stage: "成交", OpportunityCount: 20})

	items, err := svc.ListFunnel(ctx, "t1", "2026-09")
	if err != nil {
		t.Fatalf("list funnel: %v", err)
	}
	if len(items) != 3 {
		t.Errorf("got %d stages, want 3", len(items))
	}
}

// 5. 客户统计列表
func TestCustomerStatList(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()
	svc.RecordCustomerStat(ctx, &model.CustomerStat{TenantID: "t1", Period: "2026-08", NewCustomers: 10})
	svc.RecordCustomerStat(ctx, &model.CustomerStat{TenantID: "t1", Period: "2026-09", NewCustomers: 15})

	res, err := svc.ListCustomerStats(ctx, "t1", repo.Page{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if res.Total != 2 {
		t.Errorf("total = %d, want 2", res.Total)
	}
}
