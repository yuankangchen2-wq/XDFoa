package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/erp-finance-service/internal/event"
	"github.com/oa-portal/erp-finance-service/internal/model"
	"github.com/oa-portal/erp-finance-service/internal/repo"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*FinanceService, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.AutoMigrate(&model.Receivable{}, &model.Payable{}, &model.FinancePayment{})
	recvRepo := repo.NewReceivableRepo(db)
	payRepo := repo.NewPayableRepo(db)
	paymentRepo := repo.NewPaymentRepo(db)
	svc := NewFinanceService(db, recvRepo, payRepo, paymentRepo, event.NoopPublisher{})
	return svc, db
}

// 1. 创建应收后状态为未收
func TestCreateReceivable(t *testing.T) {
	svc, _ := setupTest(t)
	r := &model.Receivable{TenantID: "t1", CustomerID: 1, Amount: 1000, DueDate: "2026-12-31"}
	if err := svc.CreateReceivable(context.Background(), r); err != nil {
		t.Fatalf("create: %v", err)
	}
	if r.Status != RecvUnpaid {
		t.Errorf("status = %d, want %d", r.Status, RecvUnpaid)
	}
}

// 2. 部分收款：状态变为部分已收
func TestReceivePartial(t *testing.T) {
	svc, _ := setupTest(t)
	r := &model.Receivable{TenantID: "t1", CustomerID: 1, Amount: 1000}
	svc.CreateReceivable(context.Background(), r)

	err := svc.ReceivePayment(context.Background(), &model.FinancePayment{
		TenantID: "t1", RefID: r.ID, Amount: 400, PayMethod: "转账",
	})
	if err != nil {
		t.Fatalf("receive: %v", err)
	}

	updated, _ := svc.recvRepo.Get(context.Background(), r.ID)
	if updated.ReceivedAmount != 400 {
		t.Errorf("received_amount = %v, want 400", updated.ReceivedAmount)
	}
	if updated.Status != RecvPartial {
		t.Errorf("status = %d, want partial(%d)", updated.Status, RecvPartial)
	}
}

// 3. 全额收款：状态变为已收清
func TestReceiveFull(t *testing.T) {
	svc, _ := setupTest(t)
	r := &model.Receivable{TenantID: "t1", CustomerID: 1, Amount: 1000}
	svc.CreateReceivable(context.Background(), r)

	svc.ReceivePayment(context.Background(), &model.FinancePayment{
		TenantID: "t1", RefID: r.ID, Amount: 1000,
	})
	updated, _ := svc.recvRepo.Get(context.Background(), r.ID)
	if updated.Status != RecvSettled {
		t.Errorf("status = %d, want settled(%d)", updated.Status, RecvSettled)
	}
}

// 4. 超额收款：报错
func TestReceiveExceeds(t *testing.T) {
	svc, _ := setupTest(t)
	r := &model.Receivable{TenantID: "t1", CustomerID: 1, Amount: 1000}
	svc.CreateReceivable(context.Background(), r)

	err := svc.ReceivePayment(context.Background(), &model.FinancePayment{
		TenantID: "t1", RefID: r.ID, Amount: 1500,
	})
	if err == nil {
		t.Fatal("expected error for over-receive, got nil")
	}
}

// 5. 付款同理
func TestMakePaymentSettled(t *testing.T) {
	svc, _ := setupTest(t)
	p := &model.Payable{TenantID: "t1", SupplierID: 1, Amount: 5000}
	svc.CreatePayable(context.Background(), p)

	svc.MakePayment(context.Background(), &model.FinancePayment{
		TenantID: "t1", RefID: p.ID, Amount: 5000,
	})
	updated, _ := svc.payRepo.Get(context.Background(), p.ID)
	if updated.Status != PaySettled {
		t.Errorf("status = %d, want settled(%d)", updated.Status, PaySettled)
	}
	if updated.PaidAmount != 5000 {
		t.Errorf("paid_amount = %v, want 5000", updated.PaidAmount)
	}
}
