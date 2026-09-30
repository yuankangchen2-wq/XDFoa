package repo

import (
	"context"

	"github.com/oa-portal/audit-service/internal/model"
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

type AuditQuery struct {
	TenantID string
	Module   string
	Action   string
	UserID   string
	Status   int32
	StartAt  string
	EndAt    string
}

type AuditLogRepo struct{ db *gorm.DB }

func NewAuditLogRepo(db *gorm.DB) *AuditLogRepo { return &AuditLogRepo{db: db} }

func (r *AuditLogRepo) Create(ctx context.Context, o *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *AuditLogRepo) Get(ctx context.Context, id int64) (*model.AuditLog, error) {
	var o model.AuditLog
	err := r.db.WithContext(ctx).First(&o, id).Error
	return &o, err
}

func (r *AuditLogRepo) List(ctx context.Context, q AuditQuery, p Page) (*PageResult[model.AuditLog], error) {
	var result PageResult[model.AuditLog]
	db := r.db.WithContext(ctx).Where("tenant_id = ?", q.TenantID)
	if q.Module != "" {
		db = db.Where("module = ?", q.Module)
	}
	if q.Action != "" {
		db = db.Where("action = ?", q.Action)
	}
	if q.UserID != "" {
		db = db.Where("user_id = ?", q.UserID)
	}
	if q.Status > 0 {
		db = db.Where("status = ?", q.Status)
	}
	if q.StartAt != "" {
		db = db.Where("created_at >= ?", q.StartAt)
	}
	if q.EndAt != "" {
		db = db.Where("created_at <= ?", q.EndAt)
	}
	if err := db.Model(&model.AuditLog{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := db.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

// 按模块统计操作次数
type ModuleCount struct {
	Module string
	Count  int64
}

func (r *AuditLogRepo) CountByModule(ctx context.Context, tenantID string) ([]ModuleCount, error) {
	var rows []ModuleCount
	err := r.db.WithContext(ctx).Model(&model.AuditLog{}).
		Select("module, COUNT(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("module").Order("count DESC").Scan(&rows).Error
	return rows, err
}
