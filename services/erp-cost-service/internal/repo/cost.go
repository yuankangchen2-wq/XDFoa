package repo

import (
	"context"

	"github.com/oa-portal/erp-cost-service/internal/model"
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

// ---------- 成本中心 ----------
type CostCenterRepo struct{ db *gorm.DB }

func NewCostCenterRepo(db *gorm.DB) *CostCenterRepo { return &CostCenterRepo{db: db} }

func (r *CostCenterRepo) Create(ctx context.Context, o *model.CostCenter) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *CostCenterRepo) Get(ctx context.Context, id int64) (*model.CostCenter, error) {
	var o model.CostCenter
	err := r.db.WithContext(ctx).First(&o, id).Error
	return &o, err
}

func (r *CostCenterRepo) List(ctx context.Context, tenantID string, typ int32, p Page) (*PageResult[model.CostCenter], error) {
	var result PageResult[model.CostCenter]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if typ > 0 {
		q = q.Where("type = ?", typ)
	}
	if err := q.Model(&model.CostCenter{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

// ---------- 产品成本 ----------
type ProductCostRepo struct{ db *gorm.DB }

func NewProductCostRepo(db *gorm.DB) *ProductCostRepo { return &ProductCostRepo{db: db} }

func (r *ProductCostRepo) Create(ctx context.Context, o *model.ProductCost) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *ProductCostRepo) Get(ctx context.Context, id int64) (*model.ProductCost, error) {
	var o model.ProductCost
	err := r.db.WithContext(ctx).First(&o, id).Error
	return &o, err
}

func (r *ProductCostRepo) List(ctx context.Context, tenantID string, productID int64, p Page) (*PageResult[model.ProductCost], error) {
	var result PageResult[model.ProductCost]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if productID > 0 {
		q = q.Where("product_id = ?", productID)
	}
	if err := q.Model(&model.ProductCost{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

// 按产品+期间汇总未归集的成本记录
func (r *ProductCostRepo) SumUncollected(ctx context.Context, tenantID string, productID int64, period string) (material, labor, overhead float64, err error) {
	type elemSum struct {
		Element int32
		Total   float64
	}
	var sums []elemSum
	err = r.db.WithContext(ctx).Model(&model.CostRecord{}).
		Select("element, COALESCE(SUM(amount),0) as total").
		Where("tenant_id = ? AND product_id = ? AND period = ? AND status = ?", tenantID, productID, period, 1).
		Group("element").Scan(&sums).Error
	for _, s := range sums {
		switch s.Element {
		case 1:
			material = s.Total
		case 2:
			labor = s.Total
		case 3:
			overhead = s.Total
		}
	}
	return
}

// ---------- 成本记录 ----------
type CostRecordRepo struct{ db *gorm.DB }

func NewCostRecordRepo(db *gorm.DB) *CostRecordRepo { return &CostRecordRepo{db: db} }

func (r *CostRecordRepo) Create(ctx context.Context, o *model.CostRecord) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *CostRecordRepo) List(ctx context.Context, tenantID string, status int32, p Page) (*PageResult[model.CostRecord], error) {
	var result PageResult[model.CostRecord]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if err := q.Model(&model.CostRecord{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *CostRecordRepo) MarkCollected(ctx context.Context, tenantID string, productID int64, period string) error {
	return r.db.WithContext(ctx).Model(&model.CostRecord{}).
		Where("tenant_id = ? AND product_id = ? AND period = ? AND status = ?", tenantID, productID, period, 1).
		Update("status", 2).Error
}
