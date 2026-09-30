package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oa-portal/crm-sales-service/internal/event"
	"github.com/oa-portal/crm-sales-service/internal/model"
	"github.com/oa-portal/crm-sales-service/internal/repo"
	"gorm.io/gorm"
)

type SalesService struct {
	db           *gorm.DB
	orderRepo    *repo.SalesOrderRepo
	contractRepo *repo.ContractRepo
	paymentRepo  *repo.PaymentRepo
	publisher    event.Publisher
}

func NewSalesService(
	db *gorm.DB,
	orderRepo *repo.SalesOrderRepo,
	contractRepo *repo.ContractRepo,
	paymentRepo *repo.PaymentRepo,
	pub event.Publisher,
) *SalesService {
	return &SalesService{
		db: db, orderRepo: orderRepo, contractRepo: contractRepo,
		paymentRepo: paymentRepo, publisher: pub,
	}
}

// ---------- 销售订单 ----------

func (s *SalesService) CreateSalesOrder(ctx context.Context, o *model.SalesOrder) error {
	o.OrderNo = fmt.Sprintf("SO%d", time.Now().UnixNano())
	o.Status = 1 // 待审批
	if err := s.orderRepo.Create(ctx, o); err != nil {
		return err
	}
	s.publisher.Publish("crm-events", o.OrderNo, nil)
	return nil
}

func (s *SalesService) ListSalesOrders(ctx context.Context, tenantID string, status int32, p repo.Page) (*repo.PageResult[model.SalesOrder], error) {
	return s.orderRepo.List(ctx, tenantID, status, p)
}

func (s *SalesService) GetSalesOrder(ctx context.Context, id int64) (*model.SalesOrder, error) {
	return s.orderRepo.Get(ctx, id)
}

func (s *SalesService) ApproveSalesOrder(ctx context.Context, id int64) error {
	o, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if o.Status != 1 {
		return fmt.Errorf("sales order status invalid, current=%d", o.Status)
	}
	return s.orderRepo.UpdateStatus(ctx, id, 2)
}

func (s *SalesService) ShipSalesOrder(ctx context.Context, id int64) error {
	o, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if o.Status != 2 {
		return fmt.Errorf("sales order not approved, current=%d", o.Status)
	}
	return s.orderRepo.UpdateStatus(ctx, id, 3)
}

func (s *SalesService) CompleteSalesOrder(ctx context.Context, id int64) error {
	o, err := s.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if o.Status != 3 {
		return fmt.Errorf("sales order not shipped, current=%d", o.Status)
	}
	return s.orderRepo.UpdateStatus(ctx, id, 4)
}

// ---------- 合同 ----------

func (s *SalesService) CreateContract(ctx context.Context, c *model.Contract) error {
	c.ContractNo = fmt.Sprintf("CT%d", time.Now().UnixNano())
	c.Status = 1 // 草稿
	if err := s.contractRepo.Create(ctx, c); err != nil {
		return err
	}
	s.publisher.Publish("crm-events", c.ContractNo, nil)
	return nil
}

func (s *SalesService) ListContracts(ctx context.Context, tenantID string, status int32, p repo.Page) (*repo.PageResult[model.Contract], error) {
	return s.contractRepo.List(ctx, tenantID, status, p)
}

func (s *SalesService) GetContract(ctx context.Context, id int64) (*model.Contract, error) {
	return s.contractRepo.Get(ctx, id)
}

func (s *SalesService) ApproveContract(ctx context.Context, id int64) error {
	c, err := s.contractRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if c.Status != 1 {
		return fmt.Errorf("contract status invalid, current=%d", c.Status)
	}
	return s.contractRepo.UpdateStatus(ctx, id, 3) // 审批通过直接生效
}

func (s *SalesService) ArchiveContract(ctx context.Context, id int64) error {
	c, err := s.contractRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if c.Status != 3 {
		return fmt.Errorf("only effective contract can be archived, current=%d", c.Status)
	}
	return s.contractRepo.UpdateStatus(ctx, id, 4)
}

// ---------- 回款 ----------

func (s *SalesService) AddPayment(ctx context.Context, p *model.Payment) error {
	// 校验合同存在且生效
	c, err := s.contractRepo.Get(ctx, p.ContractID)
	if err != nil {
		return err
	}
	if c.Status != 3 {
		return fmt.Errorf("contract not effective, current=%d", c.Status)
	}
	if p.Amount <= 0 {
		return fmt.Errorf("payment amount must be positive")
	}

	// 事务：创建回款 + 更新合同已回款金额
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		return tx.Model(&model.Contract{}).Where("id = ?", p.ContractID).
			UpdateColumn("paid_amount", gorm.Expr("paid_amount + ?", p.Amount)).Error
	}); err != nil {
		return err
	}

	s.publisher.Publish("crm-events", "", event.PaymentReceived{
		TenantID: p.TenantID, ContractID: p.ContractID, Amount: p.Amount,
	})
	return nil
}

func (s *SalesService) ListPayments(ctx context.Context, contractID int64) ([]model.Payment, error) {
	return s.paymentRepo.ListByContract(ctx, contractID)
}
