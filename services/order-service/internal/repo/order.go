package repo

import (
	"context"

	"github.com/oa-portal/order-service/internal/model"
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

// ---------- 订单 ----------
type OrderRepo struct{ db *gorm.DB }

func NewOrderRepo(db *gorm.DB) *OrderRepo { return &OrderRepo{db: db} }

func (r *OrderRepo) DB() *gorm.DB { return r.db }

func (r *OrderRepo) Create(ctx context.Context, o *model.Order) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *OrderRepo) Get(ctx context.Context, id int64) (*model.Order, error) {
	var o model.Order
	err := r.db.WithContext(ctx).Preload("Items").First(&o, id).Error
	return &o, err
}

func (r *OrderRepo) List(ctx context.Context, tenantID string, userID int64, status int32, p Page) (*PageResult[model.Order], error) {
	var result PageResult[model.Order]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if userID > 0 {
		q = q.Where("user_id = ?", userID)
	}
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if err := q.Model(&model.Order{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Preload("Items").Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

// ---------- 订单日志 ----------
type OrderLogRepo struct{ db *gorm.DB }

func NewOrderLogRepo(db *gorm.DB) *OrderLogRepo { return &OrderLogRepo{db: db} }

func (r *OrderLogRepo) Create(ctx context.Context, l *model.OrderLog) error {
	return r.db.WithContext(ctx).Create(l).Error
}
