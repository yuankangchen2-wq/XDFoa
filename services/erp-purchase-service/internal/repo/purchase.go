package repo

import (
	"context"

	"github.com/oa-portal/purchase-service/internal/model"
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

// ---------- 采购订单 ----------
type PurchaseOrderRepo struct{ db *gorm.DB }

func NewPurchaseOrderRepo(db *gorm.DB) *PurchaseOrderRepo { return &PurchaseOrderRepo{db: db} }

func (r *PurchaseOrderRepo) DB() *gorm.DB { return r.db }

func (r *PurchaseOrderRepo) Create(ctx context.Context, o *model.PurchaseOrder) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *PurchaseOrderRepo) Get(ctx context.Context, id int64) (*model.PurchaseOrder, error) {
	var o model.PurchaseOrder
	err := r.db.WithContext(ctx).Preload("Items").First(&o, id).Error
	return &o, err
}

func (r *PurchaseOrderRepo) List(ctx context.Context, tenantID string, status int32, p Page) (*PageResult[model.PurchaseOrder], error) {
	var result PageResult[model.PurchaseOrder]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if err := q.Model(&model.PurchaseOrder{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Preload("Items").Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *PurchaseOrderRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.PurchaseOrder{}).Where("id = ?", id).Update("status", status).Error
}

func (r *PurchaseOrderRepo) UpdateItemReceived(ctx context.Context, orderID int64, skuCode string, qty int64) error {
	return r.db.WithContext(ctx).Model(&model.PurchaseOrderItem{}).
		Where("order_id = ? AND sku_code = ?", orderID, skuCode).
		UpdateColumn("received_qty", gorm.Expr("received_qty + ?", qty)).Error
}

// ---------- 到货单 ----------
type PurchaseReceiptRepo struct{ db *gorm.DB }

func NewPurchaseReceiptRepo(db *gorm.DB) *PurchaseReceiptRepo { return &PurchaseReceiptRepo{db: db} }

func (r *PurchaseReceiptRepo) Create(ctx context.Context, rcp *model.PurchaseReceipt) error {
	return r.db.WithContext(ctx).Create(rcp).Error
}

func (r *PurchaseReceiptRepo) Get(ctx context.Context, id int64) (*model.PurchaseReceipt, error) {
	var rcp model.PurchaseReceipt
	err := r.db.WithContext(ctx).Preload("Items").First(&rcp, id).Error
	return &rcp, err
}

func (r *PurchaseReceiptRepo) List(ctx context.Context, tenantID string, p Page) (*PageResult[model.PurchaseReceipt], error) {
	var result PageResult[model.PurchaseReceipt]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if err := q.Model(&model.PurchaseReceipt{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Preload("Items").Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *PurchaseReceiptRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.PurchaseReceipt{}).Where("id = ?", id).Update("status", status).Error
}
