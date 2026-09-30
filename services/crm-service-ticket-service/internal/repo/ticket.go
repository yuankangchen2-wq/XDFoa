package repo

import (
	"context"

	"github.com/oa-portal/crm-service-ticket-service/internal/model"
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

// ---------- 工单 ----------
type TicketRepo struct{ db *gorm.DB }

func NewTicketRepo(db *gorm.DB) *TicketRepo { return &TicketRepo{db: db} }

func (r *TicketRepo) Create(ctx context.Context, t *model.ServiceTicket) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *TicketRepo) Get(ctx context.Context, id int64) (*model.ServiceTicket, error) {
	var t model.ServiceTicket
	err := r.db.WithContext(ctx).Preload("Replies").First(&t, id).Error
	return &t, err
}

func (r *TicketRepo) List(ctx context.Context, tenantID string, status, priority int32, p Page) (*PageResult[model.ServiceTicket], error) {
	var result PageResult[model.ServiceTicket]
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if status > 0 {
		q = q.Where("status = ?", status)
	}
	if priority > 0 {
		q = q.Where("priority = ?", priority)
	}
	if err := q.Model(&model.ServiceTicket{}).Count(&result.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&result.Items).Error; err != nil {
		return nil, err
	}
	result.Page, result.PageSize = p.Page, p.PageSize
	return &result, nil
}

func (r *TicketRepo) UpdateStatus(ctx context.Context, id int64, status int32) error {
	return r.db.WithContext(ctx).Model(&model.ServiceTicket{}).Where("id = ?", id).Update("status", status).Error
}

func (r *TicketRepo) Assign(ctx context.Context, id int64, assigneeID string) error {
	return r.db.WithContext(ctx).Model(&model.ServiceTicket{}).Where("id = ?", id).
		Updates(map[string]interface{}{"assignee_id": assigneeID, "status": 2}).Error
}

func (r *TicketRepo) SetSatisfaction(ctx context.Context, id int64, score int32) error {
	return r.db.WithContext(ctx).Model(&model.ServiceTicket{}).Where("id = ?", id).Update("satisfaction", score).Error
}

// ---------- 工单回复 ----------
type ReplyRepo struct{ db *gorm.DB }

func NewReplyRepo(db *gorm.DB) *ReplyRepo { return &ReplyRepo{db: db} }

func (r *ReplyRepo) Create(ctx context.Context, rp *model.TicketReply) error {
	return r.db.WithContext(ctx).Create(rp).Error
}

func (r *ReplyRepo) ListByTicket(ctx context.Context, ticketID int64) ([]model.TicketReply, error) {
	var list []model.TicketReply
	err := r.db.WithContext(ctx).Where("ticket_id = ?", ticketID).Order("id ASC").Find(&list).Error
	return list, err
}
