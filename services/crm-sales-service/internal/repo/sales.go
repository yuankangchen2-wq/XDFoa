package repo

import (
	"context"

	"github.com/oa-portal/crm-sales-service/internal/model"
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

// ---------- 销售订单 ----------
type SalesOrderRepo struct{ db *gorm.DB }

func NewSalesOrderRepo(db *gorm.DB) *SalesOrderRepo { return &SalesOrderRepo{db: db} }

func (r *SalesOrderRepo) Create(ctx context.Context, o *model.SalesOrder) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *SalesOrderRepo) Get(ctx context.Context, id int64) (*model.SalesOrder, error) {
	var o model.SalesOrder
	err := r.db.WithContext(ctx).First(&o, id).Error
	return &o, err
}

func (r *SalesOrderRepo) List(ctx context.Context, tenantID string, status int32, p Page) (*PageResult[model.SalesOrder], error) {
	var result PageResult[model.SalesOrder]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if err := q.Model(&model.SalesOrder{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *SalesOrderRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.SalesOrder{}).Where("id = ?", id).Update("status", status).Error
}

// ---------- 合同 ----------
type ContractRepo struct{ db *gorm.DB }

func NewContractRepo(db *gorm.DB) *ContractRepo { return &ContractRepo{db: db} }

func (r *ContractRepo) Create(ctx context.Context, c *model.Contract) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *ContractRepo) Get(ctx context.Context, id int64) (*model.Contract, error) {
	var c model.Contract
	err := r.db.WithContext(ctx).First(&c, id).Error
	return &c, err
}

func (r *ContractRepo) List(ctx context.Context, tenantID string, status int32, p Page) (*PageResult[model.Contract], error) {
	var result PageResult[model.Contract]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if err := q.Model(&model.Contract{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *ContractRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.Contract{}).Where("id = ?", id).Update("status", status).Error
}

func (r *ContractRepo) AddPaidAmount(ctx context.Context, id int64, amount float64) error {
	return r.db.WithContext(ctx).Model(&model.Contract{}).Where("id = ?", id).
		UpdateColumn("paid_amount", gorm.Expr("paid_amount + ?", amount)).Error
}

// ---------- 回款 ----------
type PaymentRepo struct{ db *gorm.DB }

func NewPaymentRepo(db *gorm.DB) *PaymentRepo { return &PaymentRepo{db: db} }

func (r *PaymentRepo) Create(ctx context.Context, p *model.Payment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PaymentRepo) ListByContract(ctx context.Context, contractID int64) ([]model.Payment, error) {
	var list []model.Payment
	err := r.db.WithContext(ctx).Where("contract_id = ?", contractID).Order("id DESC").Find(&list).Error
	return list, err
}
