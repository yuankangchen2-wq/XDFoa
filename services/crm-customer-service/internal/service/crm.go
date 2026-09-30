package service

import (
	"context"

	"github.com/oa-portal/crm-customer-service/internal/event"
	"github.com/oa-portal/crm-customer-service/internal/model"
	"github.com/oa-portal/crm-customer-service/internal/repo"
)

type CrmService struct {
	customerRepo *repo.CustomerRepo
	contactRepo  *repo.ContactRepo
	followRepo   *repo.FollowUpRepo
	publisher    event.Publisher
}

func NewCrmService(
	customerRepo *repo.CustomerRepo,
	contactRepo *repo.ContactRepo,
	followRepo *repo.FollowUpRepo,
	pub event.Publisher,
) *CrmService {
	return &CrmService{
		customerRepo: customerRepo, contactRepo: contactRepo,
		followRepo: followRepo, publisher: pub,
	}
}

// ---------- 客户 ----------

func (s *CrmService) CreateCustomer(ctx context.Context, c *model.Customer) error {
	if c.Status == 0 {
		c.Status = 1
	}
	if err := s.customerRepo.Create(ctx, c); err != nil {
		return err
	}
	s.publisher.Publish("crm-events", c.Name, event.CustomerCreated{
		TenantID: c.TenantID, CustomerID: c.ID, Name: c.Name,
	})
	return nil
}

func (s *CrmService) ListCustomers(ctx context.Context, tenantID, level string, status int32, p repo.Page) (*repo.PageResult[model.Customer], error) {
	return s.customerRepo.List(ctx, tenantID, level, status, p)
}

func (s *CrmService) GetCustomer(ctx context.Context, id int64) (*model.Customer, error) {
	c, err := s.customerRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	// 加载跟进记录
	c.FollowUps, _ = s.followRepo.ListByCustomer(ctx, id)
	return c, nil
}

func (s *CrmService) UpdateCustomer(ctx context.Context, c *model.Customer) error {
	if err := s.customerRepo.Update(ctx, c); err != nil {
		return err
	}
	s.publisher.Publish("crm-events", c.Name, event.CustomerCreated{
		TenantID: c.TenantID, CustomerID: c.ID, Name: c.Name,
	})
	return nil
}

func (s *CrmService) DeleteCustomer(ctx context.Context, id int64) error {
	return s.customerRepo.Delete(ctx, id)
}

// ---------- 联系人 ----------

func (s *CrmService) AddContact(ctx context.Context, c *model.Contact) error {
	return s.contactRepo.Create(ctx, c)
}

func (s *CrmService) ListContacts(ctx context.Context, customerID int64) ([]model.Contact, error) {
	return s.contactRepo.ListByCustomer(ctx, customerID)
}

func (s *CrmService) DeleteContact(ctx context.Context, id int64) error {
	return s.contactRepo.Delete(ctx, id)
}

// ---------- 跟进记录 ----------

func (s *CrmService) AddFollowUp(ctx context.Context, f *model.FollowUp) error {
	if err := s.followRepo.Create(ctx, f); err != nil {
		return err
	}
	s.publisher.Publish("crm-events", "", event.FollowUpCreated{
		CustomerID: f.CustomerID, Type: f.Type,
	})
	return nil
}

func (s *CrmService) ListFollowUps(ctx context.Context, customerID int64) ([]model.FollowUp, error) {
	return s.followRepo.ListByCustomer(ctx, customerID)
}
