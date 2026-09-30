package repo

import (
	"context"

	"github.com/oa-portal/crm-analytics-service/internal/model"
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

// ---------- 销售快照 ----------
type SalesSnapshotRepo struct{ db *gorm.DB }

func NewSalesSnapshotRepo(db *gorm.DB) *SalesSnapshotRepo { return &SalesSnapshotRepo{db: db} }

func (r *SalesSnapshotRepo) Create(ctx context.Context, o *model.SalesSnapshot) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *SalesSnapshotRepo) List(ctx context.Context, tenantID, period string, p Page) (*PageResult[model.SalesSnapshot], error) {
	var result PageResult[model.SalesSnapshot]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if period != "" {
		q = q.Where("period = ?", period)
	}
	if err := q.Model(&model.SalesSnapshot{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

// 按期间汇总
type SalesSummary struct {
	Period      string
	TotalAmount float64
	OrderCount  int32
	AvgAmount   float64
}

func (r *SalesSnapshotRepo) SummaryByPeriod(ctx context.Context, tenantID string) ([]SalesSummary, error) {
	var rows []SalesSummary
	err := r.db.WithContext(ctx).Model(&model.SalesSnapshot{}).
		Select("period, COALESCE(SUM(total_amount),0) as total_amount, COALESCE(SUM(order_count),0) as order_count, COALESCE(AVG(avg_amount),0) as avg_amount").
		Where("tenant_id = ?", tenantID).
		Group("period").Order("period ASC").Scan(&rows).Error
	return rows, err
}

// 销售总额
func (r *SalesSnapshotRepo) TotalSales(ctx context.Context, tenantID string) (float64, int32, error) {
	type agg struct {
		Total  float64
		Orders int32
	}
	var a agg
	err := r.db.WithContext(ctx).Model(&model.SalesSnapshot{}).
		Select("COALESCE(SUM(total_amount),0) as total, COALESCE(SUM(order_count),0) as orders").
		Where("tenant_id = ?", tenantID).Scan(&a).Error
	return a.Total, a.Orders, err
}

// ---------- 漏斗 ----------
type FunnelStageRepo struct{ db *gorm.DB }

func NewFunnelStageRepo(db *gorm.DB) *FunnelStageRepo { return &FunnelStageRepo{db: db} }

func (r *FunnelStageRepo) Create(ctx context.Context, o *model.FunnelStage) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *FunnelStageRepo) List(ctx context.Context, tenantID, period string) ([]model.FunnelStage, error) {
	var items []model.FunnelStage
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if period != "" {
		q = q.Where("period = ?", period)
	}
	err := q.Order("id ASC").Find(&items).Error
	return items, err
}

// ---------- 客户统计 ----------
type CustomerStatRepo struct{ db *gorm.DB }

func NewCustomerStatRepo(db *gorm.DB) *CustomerStatRepo { return &CustomerStatRepo{db: db} }

func (r *CustomerStatRepo) Create(ctx context.Context, o *model.CustomerStat) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *CustomerStatRepo) List(ctx context.Context, tenantID string, p Page) (*PageResult[model.CustomerStat], error) {
	var result PageResult[model.CustomerStat]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if err := q.Model(&model.CustomerStat{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

// 最新一期客户统计
func (r *CustomerStatRepo) Latest(ctx context.Context, tenantID string) (*model.CustomerStat, error) {
	var o model.CustomerStat
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("id DESC").First(&o).Error
	return &o, err
}
