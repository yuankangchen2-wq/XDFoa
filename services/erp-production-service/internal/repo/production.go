package repo

import (
	"context"

	"github.com/oa-portal/production-service/internal/model"
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

// ---------- BOM ----------
type BomRepo struct{ db *gorm.DB }

func NewBomRepo(db *gorm.DB) *BomRepo { return &BomRepo{db: db} }

func (r *BomRepo) Create(ctx context.Context, b *model.Bom) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *BomRepo) Get(ctx context.Context, id int64) (*model.Bom, error) {
	var b model.Bom
	err := r.db.WithContext(ctx).Preload("Items").First(&b, id).Error
	return &b, err
}

func (r *BomRepo) GetByProductSku(ctx context.Context, tenantID, productSku string) (*model.Bom, error) {
	var b model.Bom
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND product_sku = ? AND status = 2", tenantID, productSku).
		Preload("Items").First(&b).Error
	return &b, err
}

func (r *BomRepo) List(ctx context.Context, tenantID string, p Page) (*PageResult[model.Bom], error) {
	var result PageResult[model.Bom]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if err := q.Model(&model.Bom{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Preload("Items").Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *BomRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.Bom{}).Where("id = ?", id).Update("status", status).Error
}

// ---------- 生产工单 ----------
type WorkOrderRepo struct{ db *gorm.DB }

func NewWorkOrderRepo(db *gorm.DB) *WorkOrderRepo { return &WorkOrderRepo{db: db} }

func (r *WorkOrderRepo) Create(ctx context.Context, w *model.WorkOrder) error {
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *WorkOrderRepo) Get(ctx context.Context, id int64) (*model.WorkOrder, error) {
	var w model.WorkOrder
	err := r.db.WithContext(ctx).First(&w, id).Error
	return &w, err
}

func (r *WorkOrderRepo) List(ctx context.Context, tenantID string, status int32, p Page) (*PageResult[model.WorkOrder], error) {
	var result PageResult[model.WorkOrder]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if err := q.Model(&model.WorkOrder{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *WorkOrderRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.WorkOrder{}).Where("id = ?", id).Update("status", status).Error
}

func (r *WorkOrderRepo) AddProducedQty(ctx context.Context, id int64, qty int64) error {
	return r.db.WithContext(ctx).Model(&model.WorkOrder{}).Where("id = ?", id).
		UpdateColumn("produced_qty", gorm.Expr("produced_qty + ?", qty)).Error
}

// ---------- 领料单 ----------
type MaterialIssueRepo struct{ db *gorm.DB }

func NewMaterialIssueRepo(db *gorm.DB) *MaterialIssueRepo { return &MaterialIssueRepo{db: db} }

func (r *MaterialIssueRepo) Create(ctx context.Context, m *model.MaterialIssue) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MaterialIssueRepo) Get(ctx context.Context, id int64) (*model.MaterialIssue, error) {
	var m model.MaterialIssue
	err := r.db.WithContext(ctx).Preload("Items").First(&m, id).Error
	return &m, err
}

func (r *MaterialIssueRepo) List(ctx context.Context, tenantID string, p Page) (*PageResult[model.MaterialIssue], error) {
	var result PageResult[model.MaterialIssue]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if err := q.Model(&model.MaterialIssue{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Preload("Items").Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *MaterialIssueRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.MaterialIssue{}).Where("id = ?", id).Update("status", status).Error
}

// ---------- 报工 ----------
type ProductionReportRepo struct{ db *gorm.DB }

func NewProductionReportRepo(db *gorm.DB) *ProductionReportRepo { return &ProductionReportRepo{db: db} }

func (r *ProductionReportRepo) Create(ctx context.Context, p *model.ProductionReport) error {
	return r.db.WithContext(ctx).Create(p).Error
}
