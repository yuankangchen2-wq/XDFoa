package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/oa-portal/mdm-service/internal/event"
	"github.com/oa-portal/mdm-service/internal/model"
	"github.com/oa-portal/mdm-service/internal/repo"
)

// ---------- 客户 ----------
type CustomerService struct {
	repo      *repo.CustomerRepo
	publisher event.Publisher
}

func NewCustomerService(r *repo.CustomerRepo, p event.Publisher) *CustomerService {
	return &CustomerService{repo: r, publisher: p}
}

func (s *CustomerService) Create(ctx context.Context, c *model.Customer) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	if err := s.repo.Create(ctx, c); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		c.TenantID, "customer", "created", strconv.FormatInt(c.ID, 10), c,
	))
}

func (s *CustomerService) Update(ctx context.Context, c *model.Customer) error {
	c.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, c); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		c.TenantID, "customer", "updated", strconv.FormatInt(c.ID, 10), c,
	))
}

func (s *CustomerService) Get(ctx context.Context, id int64, tenantID string) (*model.Customer, error) {
	return s.repo.GetByID(ctx, id, tenantID)
}

func (s *CustomerService) List(ctx context.Context, tenantID, keyword, customerType string, page, pageSize int) (*repo.PageResult[model.Customer], error) {
	return s.repo.List(ctx, tenantID, keyword, customerType, repo.PageQuery{Page: page, PageSize: pageSize})
}

func (s *CustomerService) Delete(ctx context.Context, id int64, tenantID string) error {
	if err := s.repo.Delete(ctx, id, tenantID); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		tenantID, "customer", "deleted", strconv.FormatInt(id, 10), nil,
	))
}

// ---------- 商品 ----------
type ProductService struct {
	repo      *repo.ProductRepo
	publisher event.Publisher
}

func NewProductService(r *repo.ProductRepo, p event.Publisher) *ProductService {
	return &ProductService{repo: r, publisher: p}
}

func (s *ProductService) Create(ctx context.Context, p *model.Product) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if err := s.repo.Create(ctx, p); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		p.TenantID, "product", "created", strconv.FormatInt(p.ID, 10), p,
	))
}

func (s *ProductService) Update(ctx context.Context, p *model.Product) error {
	p.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, p); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		p.TenantID, "product", "updated", strconv.FormatInt(p.ID, 10), p,
	))
}

func (s *ProductService) Get(ctx context.Context, id int64, tenantID string) (*model.Product, error) {
	return s.repo.GetByID(ctx, id, tenantID)
}

func (s *ProductService) List(ctx context.Context, tenantID, keyword, category string, page, pageSize int) (*repo.PageResult[model.Product], error) {
	return s.repo.List(ctx, tenantID, keyword, category, repo.PageQuery{Page: page, PageSize: pageSize})
}

func (s *ProductService) Delete(ctx context.Context, id int64, tenantID string) error {
	if err := s.repo.Delete(ctx, id, tenantID); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		tenantID, "product", "deleted", strconv.FormatInt(id, 10), nil,
	))
}

// ---------- 供应商 ----------
type SupplierService struct {
	repo      *repo.SupplierRepo
	publisher event.Publisher
}

func NewSupplierService(r *repo.SupplierRepo, p event.Publisher) *SupplierService {
	return &SupplierService{repo: r, publisher: p}
}

func (s *SupplierService) Create(ctx context.Context, sup *model.Supplier) error {
	sup.CreatedAt = time.Now()
	sup.UpdatedAt = time.Now()
	if err := s.repo.Create(ctx, sup); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		sup.TenantID, "supplier", "created", strconv.FormatInt(sup.ID, 10), sup,
	))
}

func (s *SupplierService) Update(ctx context.Context, sup *model.Supplier) error {
	sup.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, sup); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		sup.TenantID, "supplier", "updated", strconv.FormatInt(sup.ID, 10), sup,
	))
}

func (s *SupplierService) Get(ctx context.Context, id int64, tenantID string) (*model.Supplier, error) {
	return s.repo.GetByID(ctx, id, tenantID)
}

func (s *SupplierService) List(ctx context.Context, tenantID, keyword string, page, pageSize int) (*repo.PageResult[model.Supplier], error) {
	return s.repo.List(ctx, tenantID, keyword, repo.PageQuery{Page: page, PageSize: pageSize})
}

func (s *SupplierService) Delete(ctx context.Context, id int64, tenantID string) error {
	if err := s.repo.Delete(ctx, id, tenantID); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		tenantID, "supplier", "deleted", strconv.FormatInt(id, 10), nil,
	))
}

// ---------- 组织 ----------
type OrganizationService struct {
	repo      *repo.OrganizationRepo
	publisher event.Publisher
}

func NewOrganizationService(r *repo.OrganizationRepo, p event.Publisher) *OrganizationService {
	return &OrganizationService{repo: r, publisher: p}
}

func (s *OrganizationService) Create(ctx context.Context, o *model.Organization) error {
	// 计算 org_path 和 org_level
	if o.ParentID == 0 {
		o.OrgPath = fmt.Sprintf("/%s", o.OrgCode)
		o.OrgLevel = 1
	} else {
		parent, err := s.repo.GetByID(ctx, o.ParentID, o.TenantID)
		if err != nil {
			return fmt.Errorf("parent org not found: %w", err)
		}
		o.OrgPath = fmt.Sprintf("%s/%s", parent.OrgPath, o.OrgCode)
		o.OrgLevel = parent.OrgLevel + 1
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	if err := s.repo.Create(ctx, o); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		o.TenantID, "organization", "created", strconv.FormatInt(o.ID, 10), o,
	))
}

func (s *OrganizationService) Update(ctx context.Context, o *model.Organization) error {
	o.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, o); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		o.TenantID, "organization", "updated", strconv.FormatInt(o.ID, 10), o,
	))
}

func (s *OrganizationService) Get(ctx context.Context, id int64, tenantID string) (*model.Organization, error) {
	return s.repo.GetByID(ctx, id, tenantID)
}

func (s *OrganizationService) List(ctx context.Context, tenantID string, parentID int64) ([]model.Organization, error) {
	return s.repo.ListByParent(ctx, tenantID, parentID)
}

func (s *OrganizationService) Delete(ctx context.Context, id int64, tenantID string) error {
	if err := s.repo.Delete(ctx, id, tenantID); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event.NewEntitySyncedEvent(
		tenantID, "organization", "deleted", strconv.FormatInt(id, 10), nil,
	))
}
