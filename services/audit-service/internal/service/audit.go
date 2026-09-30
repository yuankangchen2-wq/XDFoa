package service

import (
	"context"
	"time"

	"github.com/oa-portal/audit-service/internal/event"
	"github.com/oa-portal/audit-service/internal/model"
	"github.com/oa-portal/audit-service/internal/repo"
)

const (
	StatusSuccess = 1
	StatusFail    = 2
)

type AuditService struct {
	repo      *repo.AuditLogRepo
	publisher event.Publisher
}

func NewAuditService(r *repo.AuditLogRepo, pub event.Publisher) *AuditService {
	return &AuditService{repo: r, publisher: pub}
}

// Record 记录审计日志
func (s *AuditService) Record(ctx context.Context, log *model.AuditLog) error {
	if log.Status == 0 {
		log.Status = StatusSuccess
	}
	log.CreatedAt = time.Now()
	if err := s.repo.Create(ctx, log); err != nil {
		return err
	}
	s.publisher.Publish("audit-events", "", nil)
	return nil
}

// Query 按条件查询审计日志
func (s *AuditService) Query(ctx context.Context, q repo.AuditQuery, p repo.Page) (*repo.PageResult[model.AuditLog], error) {
	return s.repo.List(ctx, q, p)
}

// Get 获取单条
func (s *AuditService) Get(ctx context.Context, id int64) (*model.AuditLog, error) {
	return s.repo.Get(ctx, id)
}

// ModuleStats 按模块统计
func (s *AuditService) ModuleStats(ctx context.Context, tenantID string) ([]repo.ModuleCount, error) {
	return s.repo.CountByModule(ctx, tenantID)
}
