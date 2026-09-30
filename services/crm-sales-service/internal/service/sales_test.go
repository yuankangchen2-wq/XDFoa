package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/crm-sales-service/internal/event"
	"github.com/oa-portal/crm-sales-service/internal/model"
	"github.com/oa-portal/crm-sales-service/internal/repo"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) *SalesService {
	t.Helper()
	dbFile := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&model.SalesOrder{}, &model.Contract{}, &model.Payment{})

	return NewSalesService(
		db,
		repo.NewSalesOrderRepo(db),
		repo.NewContractRepo(db),
		repo.NewPaymentRepo(db),
		event.NoopPublisher{},
	)
}

// 测试销售订单完整流程
func TestSalesOrderFlow(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	o := &model.SalesOrder{
		TenantID: "t1", CustomerID: 1, TotalAmount: 50000,
	}
	if err := svc.CreateSalesOrder(ctx, o); err != nil {
		t.Fatal(err)
	}
	if o.Status != 1 {
		t.Fatalf("expected status=1, got %d", o.Status)
	}

	// 审批
	if err := svc.ApproveSalesOrder(ctx, o.ID); err != nil {
		t.Fatal(err)
	}
	// 发货
	if err := svc.ShipSalesOrder(ctx, o.ID); err != nil {
		t.Fatal(err)
	}
	// 完成
	if err := svc.CompleteSalesOrder(ctx, o.ID); err != nil {
		t.Fatal(err)
	}

	final, _ := svc.GetSalesOrder(ctx, o.ID)
	if final.Status != 4 {
		t.Fatalf("expected status=4, got %d", final.Status)
	}
}

// 测试销售订单状态校验
func TestSalesOrderStatusValidation(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	o := &model.SalesOrder{TenantID: "t1", CustomerID: 1, TotalAmount: 10000}
	svc.CreateSalesOrder(ctx, o)

	// 待审批不能直接发货
	if err := svc.ShipSalesOrder(ctx, o.ID); err == nil {
		t.Fatal("expected error when shipping unapproved order")
	}
}

// 测试合同 + 回款
func TestContractAndPayment(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	c := &model.Contract{
		TenantID: "t1", CustomerID: 1, Name: "测试合同",
		Amount: 100000, SignDate: "2026-09-01",
	}
	if err := svc.CreateContract(ctx, c); err != nil {
		t.Fatal(err)
	}
	if c.Status != 1 {
		t.Fatalf("expected status=1, got %d", c.Status)
	}

	// 草稿状态不能回款
	p := &model.Payment{TenantID: "t1", ContractID: c.ID, Amount: 30000, PayDate: "2026-09-10", PayMethod: "银行转账"}
	if err := svc.AddPayment(ctx, p); err == nil {
		t.Fatal("expected error when paying draft contract")
	}

	// 审批合同
	if err := svc.ApproveContract(ctx, c.ID); err != nil {
		t.Fatal(err)
	}

	// 回款
	if err := svc.AddPayment(ctx, p); err != nil {
		t.Fatal(err)
	}

	// 验证合同已回款金额
	updated, _ := svc.GetContract(ctx, c.ID)
	if updated.PaidAmount != 30000 {
		t.Fatalf("expected paid_amount=30000, got %f", updated.PaidAmount)
	}

	// 验证回款记录
	payments, _ := svc.ListPayments(ctx, c.ID)
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment, got %d", len(payments))
	}
}

// 测试合同归档
func TestArchiveContract(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	c := &model.Contract{TenantID: "t1", CustomerID: 1, Name: "合同A", Amount: 50000}
	svc.CreateContract(ctx, c)

	// 未生效不能归档
	if err := svc.ArchiveContract(ctx, c.ID); err == nil {
		t.Fatal("expected error when archiving non-effective contract")
	}

	// 审批后归档
	svc.ApproveContract(ctx, c.ID)
	if err := svc.ArchiveContract(ctx, c.ID); err != nil {
		t.Fatal(err)
	}
	updated, _ := svc.GetContract(ctx, c.ID)
	if updated.Status != 4 {
		t.Fatalf("expected status=4, got %d", updated.Status)
	}
}
