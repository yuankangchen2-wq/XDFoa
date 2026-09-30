package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/crm-opportunity-service/internal/event"
	"github.com/oa-portal/crm-opportunity-service/internal/model"
	"github.com/oa-portal/crm-opportunity-service/internal/repo"
	"gorm.io/gorm"
)

type mockOrderClient struct {
	createCalled bool
}

func (m *mockOrderClient) CreateOrder(tenantID string, userID, warehouseID int64, items []map[string]interface{}) error {
	m.createCalled = true
	return nil
}

func setupTest(t *testing.T) (*OpportunityService, *mockOrderClient) {
	t.Helper()
	dbFile := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&model.Opportunity{})

	mockOrder := &mockOrderClient{}
	svc := NewOpportunityService(
		repo.NewOpportunityRepo(db),
		mockOrder,
		event.NoopPublisher{},
	)
	return svc, mockOrder
}

// 测试创建商机
func TestCreateOpportunity(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	o := &model.Opportunity{
		TenantID: "t1", Name: "测试商机", CustomerID: 1,
		Amount: 50000, OwnerID: "sales001",
	}
	if err := svc.CreateOpportunity(ctx, o); err != nil {
		t.Fatal(err)
	}
	if o.ID == 0 {
		t.Fatal("id should not be 0")
	}
	if o.Stage != StageLead {
		t.Fatalf("expected stage=%d, got %d", StageLead, o.Stage)
	}
	if o.WinRate != 10 {
		t.Fatalf("expected win_rate=10, got %d", o.WinRate)
	}
}

// 测试阶段流转
func TestTransitionStage(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	o := &model.Opportunity{
		TenantID: "t1", Name: "商机A", CustomerID: 1, Amount: 10000,
	}
	svc.CreateOpportunity(ctx, o)

	// 流转到报价
	if err := svc.TransitionStage(ctx, int32(o.ID), StageQuote); err != nil {
		t.Fatal(err)
	}
	updated, _ := svc.GetOpportunity(ctx, o.ID)
	if updated.Stage != StageQuote {
		t.Fatalf("expected stage=%d, got %d", StageQuote, updated.Stage)
	}
}

// 测试成交后不能再流转
func TestTransitionWonCannotChange(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	o := &model.Opportunity{
		TenantID: "t1", Name: "商机A", CustomerID: 1, Amount: 10000,
	}
	svc.CreateOpportunity(ctx, o)
	svc.TransitionStage(ctx, int32(o.ID), StageWon)

	// 成交后再流转应该失败
	err := svc.TransitionStage(ctx, int32(o.ID), StageLost)
	if err == nil {
		t.Fatal("expected error when transitioning won opportunity")
	}
}

// 测试商机转订单
func TestConvertToOrder(t *testing.T) {
	svc, mockOrder := setupTest(t)
	ctx := context.Background()

	o := &model.Opportunity{
		TenantID: "t1", Name: "商机A", CustomerID: 1, Amount: 80000,
	}
	svc.CreateOpportunity(ctx, o)

	// 未成交不能转订单
	err := svc.ConvertToOrder(ctx, o.ID, 1001, 1)
	if err == nil {
		t.Fatal("expected error when converting non-won opportunity")
	}

	// 成交后转订单
	svc.TransitionStage(ctx, int32(o.ID), StageWon)
	if err := svc.ConvertToOrder(ctx, o.ID, 1001, 1); err != nil {
		t.Fatal(err)
	}
	if !mockOrder.createCalled {
		t.Fatal("order client should be called")
	}
}

// 测试按阶段筛选
func TestListByStage(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()

	o1 := &model.Opportunity{TenantID: "t1", Name: "商机A", CustomerID: 1, Amount: 10000, Stage: StageLead}
	o2 := &model.Opportunity{TenantID: "t1", Name: "商机B", CustomerID: 1, Amount: 20000, Stage: StageQuote}
	svc.CreateOpportunity(ctx, o1)
	svc.CreateOpportunity(ctx, o2)

	res, err := svc.ListOpportunities(ctx, "t1", "", StageQuote, repo.Page{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 1 {
		t.Fatalf("expected 1 quote-stage opportunity, got %d", res.Total)
	}
}
