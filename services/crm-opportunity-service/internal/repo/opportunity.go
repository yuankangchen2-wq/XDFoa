package repo

import (
	"context"

	"github.com/oa-portal/crm-opportunity-service/internal/model"
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

type OpportunityRepo struct{ db *gorm.DB }

func NewOpportunityRepo(db *gorm.DB) *OpportunityRepo { return &OpportunityRepo{db: db} }

func (r *OpportunityRepo) Create(ctx context.Context, o *model.Opportunity) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *OpportunityRepo) Get(ctx context.Context, id int64) (*model.Opportunity, error) {
	var o model.Opportunity
	err := r.db.WithContext(ctx).First(&o, id).Error
	return &o, err
}

func (r *OpportunityRepo) List(ctx context.Context, tenantID, ownerID string, stage int32, p Page) (*PageResult[model.Opportunity], error) {
	var result PageResult[model.Opportunity]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if ownerID != "" {
		q = q.Where("owner_id = ?", ownerID)
	}
	if stage > 0 {
		q = q.Where("stage = ?", stage)
	}
	if err := q.Model(&model.Opportunity{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *OpportunityRepo) Update(ctx context.Context, o *model.Opportunity) error {
	return r.db.WithContext(ctx).Save(o).Error
}

func (r *OpportunityRepo) UpdateStage(ctx context.Context, id int64, stage int32) error {
	return r.db.WithContext(ctx).Model(&model.Opportunity{}).Where("id = ?", id).Update("stage", stage).Error
}

func (r *OpportunityRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Opportunity{}, id).Error
}
