package repo

import (
	"context"

	"github.com/oa-portal/crm-customer-service/internal/model"
	"gorm.io/gorm"
)

type Page struct {
	Page     int
	PageSize int
}

type PageResult[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

func Paginate(p Page) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.Page <= 0 {
			p.Page = 1
		}
		if p.PageSize <= 0 || p.PageSize > 100 {
			p.PageSize = 20
		}
		return db.Offset((p.Page - 1) * p.PageSize).Limit(p.PageSize)
	}
}

// ---------- 客户 ----------
type CustomerRepo struct{ db *gorm.DB }

func NewCustomerRepo(db *gorm.DB) *CustomerRepo { return &CustomerRepo{db: db} }

func (r *CustomerRepo) Create(ctx context.Context, c *model.Customer) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CustomerRepo) Get(ctx context.Context, id int64) (*model.Customer, error) {
	var c model.Customer
	err := r.db.WithContext(ctx).Preload("Contacts").First(&c, id).Error
	return &c, err
}

func (r *CustomerRepo) List(ctx context.Context, tenantID, level string, status int32, p Page) (*PageResult[model.Customer], error) {
	var result PageResult[model.Customer]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if level != "" {
		q = q.Where("level = ?", level)
	}
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if err := q.Model(&model.Customer{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Preload("Contacts").Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *CustomerRepo) Update(ctx context.Context, c *model.Customer) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *CustomerRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Customer{}, id).Error
}

// ---------- 联系人 ----------
type ContactRepo struct{ db *gorm.DB }

func NewContactRepo(db *gorm.DB) *ContactRepo { return &ContactRepo{db: db} }

func (r *ContactRepo) Create(ctx context.Context, c *model.Contact) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *ContactRepo) ListByCustomer(ctx context.Context, customerID int64) ([]model.Contact, error) {
	var list []model.Contact
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&list).Error
	return list, err
}

func (r *ContactRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Contact{}, id).Error
}

// ---------- 跟进记录 ----------
type FollowUpRepo struct{ db *gorm.DB }

func NewFollowUpRepo(db *gorm.DB) *FollowUpRepo { return &FollowUpRepo{db: db} }

func (r *FollowUpRepo) Create(ctx context.Context, f *model.FollowUp) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *FollowUpRepo) ListByCustomer(ctx context.Context, customerID int64) ([]model.FollowUp, error) {
	var list []model.FollowUp
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Order("follow_up_at DESC").Find(&list).Error
	return list, err
}
