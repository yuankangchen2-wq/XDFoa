package service

import (
	"context"
	"fmt"

	"github.com/oa-portal/erp-finance-service/internal/event"
	"github.com/oa-portal/erp-finance-service/internal/model"
	"github.com/oa-portal/erp-finance-service/internal/repo"
	"gorm.io/gorm"
)

const (
	// 应收状态
	RecvUnpaid     = 1 // 未收
	RecvPartial    = 2 // 部分已收
	RecvSettled    = 3 // 已收清
	// 应付状态
	PayUnpaid      = 1 // 未付
	PayPartial     = 2 // 部分已付
	PaySettled     = 3 // 已付清
	// 收付款类型
	TypeReceive = 1
	TypePay     = 2
)

type FinanceService struct {
	db           *gorm.DB
	recvRepo     *repo.ReceivableRepo
	payRepo      *repo.PayableRepo
	paymentRepo  *repo.PaymentRepo
	publisher    event.Publisher
}

func NewFinanceService(
	db *gorm.DB,
	recvRepo *repo.ReceivableRepo,
	payRepo *repo.PayableRepo,
	paymentRepo *repo.PaymentRepo,
	pub event.Publisher,
) *FinanceService {
	return &FinanceService{
		db: db, recvRepo: recvRepo, payRepo: payRepo,
		paymentRepo: paymentRepo, publisher: pub,
	}
}

// ---------- 应收账款 ----------

func (s *FinanceService) CreateReceivable(ctx context.Context, r *model.Receivable) error {
	r.Status = RecvUnpaid
	return s.recvRepo.Create(ctx, r)
}

func (s *FinanceService) ListReceivables(ctx context.Context, tenantID string, status int32, p repo.Page) (*repo.PageResult[model.Receivable], error) {
	return s.recvRepo.List(ctx, tenantID, status, p)
}

// ---------- 应付账款 ----------

func (s *FinanceService) CreatePayable(ctx context.Context, p *model.Payable) error {
	p.Status = PayUnpaid
	return s.payRepo.Create(ctx, p)
}

func (s *FinanceService) ListPayables(ctx context.Context, tenantID string, status int32, p repo.Page) (*repo.PageResult[model.Payable], error) {
	return s.payRepo.List(ctx, tenantID, status, p)
}

// ---------- 收款 ----------

// ReceivePayment 收款：更新应收已收金额 + 状态
func (s *FinanceService) ReceivePayment(ctx context.Context, p *model.FinancePayment) error {
	p.Type = TypeReceive
	r, err := s.recvRepo.Get(ctx, p.RefID)
	if err != nil {
		return fmt.Errorf("receivable not found: %w", err)
	}
	if p.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if r.ReceivedAmount+p.Amount > r.Amount+0.001 {
		return fmt.Errorf("received amount exceeds receivable amount")
	}

	// 事务：创建收款记录 + 更新应收
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		return tx.Model(&model.Receivable{}).Where("id = ?", p.RefID).
			UpdateColumn("received_amount", gorm.Expr("received_amount + ?", p.Amount)).Error
	}); err != nil {
		return err
	}

	// 更新状态
	updated, _ := s.recvRepo.Get(ctx, p.RefID)
	newStatus := RecvPartial
	if updated.ReceivedAmount >= updated.Amount-0.001 {
		newStatus = RecvSettled
	}
	s.recvRepo.UpdateStatus(ctx, p.RefID, int32(newStatus))

	s.publisher.Publish("finance-events", "", nil)
	return nil
}

// ---------- 付款 ----------

// MakePayment 付款：更新应付已付金额 + 状态
func (s *FinanceService) MakePayment(ctx context.Context, p *model.FinancePayment) error {
	p.Type = TypePay
	pa, err := s.payRepo.Get(ctx, p.RefID)
	if err != nil {
		return fmt.Errorf("payable not found: %w", err)
	}
	if p.Amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if pa.PaidAmount+p.Amount > pa.Amount+0.001 {
		return fmt.Errorf("paid amount exceeds payable amount")
	}

	// 事务：创建付款记录 + 更新应付
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		return tx.Model(&model.Payable{}).Where("id = ?", p.RefID).
			UpdateColumn("paid_amount", gorm.Expr("paid_amount + ?", p.Amount)).Error
	}); err != nil {
		return err
	}

	// 更新状态
	updated, _ := s.payRepo.Get(ctx, p.RefID)
	newStatus := PayPartial
	if updated.PaidAmount >= updated.Amount-0.001 {
		newStatus = PaySettled
	}
	s.payRepo.UpdateStatus(ctx, p.RefID, int32(newStatus))

	s.publisher.Publish("finance-events", "", nil)
	return nil
}

// ---------- 收付款记录 ----------

func (s *FinanceService) ListPayments(ctx context.Context, tenantID string, typ int32, p repo.Page) (*repo.PageResult[model.FinancePayment], error) {
	return s.paymentRepo.List(ctx, tenantID, typ, p)
}
