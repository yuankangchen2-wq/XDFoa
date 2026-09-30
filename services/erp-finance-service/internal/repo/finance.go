package repo

import (
	"context"

	"github.com/oa-portal/erp-finance-service/internal/model"
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

// ---------- 应收账款 ----------
type ReceivableRepo struct{ db *gorm.DB }

func NewReceivableRepo(db *gorm.DB) *ReceivableRepo { return &ReceivableRepo{db: db} }

func (r *ReceivableRepo) Create(ctx context.Context, o *model.Receivable) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *ReceivableRepo) Get(ctx context.Context, id int64) (*model.Receivable, error) {
	var o model.Receivable
	err := r.db.WithContext(ctx).First(&o, id).Error
	return &o, err
}

func (r *ReceivableRepo) List(ctx context.Context, tenantID string, status int32, p Page) (*PageResult[model.Receivable], error) {
	var result PageResult[model.Receivable]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if err := q.Model(&model.Receivable{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *ReceivableRepo) AddReceived(ctx context.Context, id int64, amount float64) error {
	return r.db.WithContext(ctx).Model(&model.Receivable{}).Where("id = ?", id).
		UpdateColumn("received_amount", gorm.Expr("received_amount + ?", amount)).Error
}

func (r *ReceivableRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.Receivable{}).Where("id = ?", id).Update("status", status).Error
}

// ---------- 应付账款 ----------
type PayableRepo struct{ db *gorm.DB }

func NewPayableRepo(db *gorm.DB) *PayableRepo { return &PayableRepo{db: db} }

func (r *PayableRepo) Create(ctx context.Context, o *model.Payable) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *PayableRepo) Get(ctx context.Context, id int64) (*model.Payable, error) {
	var o model.Payable
	err := r.db.WithContext(ctx).First(&o, id).Error
	return &o, err
}

func (r *PayableRepo) List(ctx context.Context, tenantID string, status int32, p Page) (*PageResult[model.Payable], error) {
	var result PageResult[model.Payable]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if err := q.Model(&model.Payable{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *PayableRepo) AddPaid(ctx context.Context, id int64, amount float64) error {
	return r.db.WithContext(ctx).Model(&model.Payable{}).Where("id = ?", id).
		UpdateColumn("paid_amount", gorm.Expr("paid_amount + ?", amount)).Error
}

func (r *PayableRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.Payable{}).Where("id = ?", id).Update("status", status).Error
}

// ---------- 收付款记录 ----------
type PaymentRepo struct{ db *gorm.DB }

func NewPaymentRepo(db *gorm.DB) *PaymentRepo { return &PaymentRepo{db: db} }

func (r *PaymentRepo) Create(ctx context.Context, p *model.FinancePayment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PaymentRepo) List(ctx context.Context, tenantID string, typ int32, p Page) (*PageResult[model.FinancePayment], error) {
	var result PageResult[model.FinancePayment]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if typ > 0 {
		q = q.Where("type = ?", typ)
	}
	if err := q.Model(&model.FinancePayment{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}
