package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oa-portal/crm-service-ticket-service/internal/event"
	"github.com/oa-portal/crm-service-ticket-service/internal/model"
	"github.com/oa-portal/crm-service-ticket-service/internal/repo"
)

// 工单状态
const (
	StatusPending  = 1 // 待处理
	StatusProcessing = 2 // 处理中
	StatusResolved = 3 // 已解决
	StatusClosed   = 4 // 已关闭
)

type TicketService struct {
	ticketRepo *repo.TicketRepo
	replyRepo  *repo.ReplyRepo
	publisher  event.Publisher
}

func NewTicketService(
	ticketRepo *repo.TicketRepo,
	replyRepo *repo.ReplyRepo,
	pub event.Publisher,
) *TicketService {
	return &TicketService{
		ticketRepo: ticketRepo, replyRepo: replyRepo, publisher: pub,
	}
}

// CreateTicket 创建工单
func (s *TicketService) CreateTicket(ctx context.Context, t *model.ServiceTicket) error {
	t.TicketNo = fmt.Sprintf("TK%d", time.Now().UnixNano())
	t.Status = StatusPending
	if t.Priority == 0 {
		t.Priority = 2
	}
	if err := s.ticketRepo.Create(ctx, t); err != nil {
		return err
	}
	s.publisher.Publish("crm-events", t.TicketNo, nil)
	return nil
}

// ListTickets 查询工单
func (s *TicketService) ListTickets(ctx context.Context, tenantID string, status, priority int32, p repo.Page) (*repo.PageResult[model.ServiceTicket], error) {
	return s.ticketRepo.List(ctx, tenantID, status, priority, p)
}

// GetTicket 获取工单详情（含回复）
func (s *TicketService) GetTicket(ctx context.Context, id int64) (*model.ServiceTicket, error) {
	return s.ticketRepo.Get(ctx, id)
}

// UpdateStatus 更新工单状态
func (s *TicketService) UpdateStatus(ctx context.Context, id int64, status int32) error {
	if status < StatusPending || status > StatusClosed {
		return fmt.Errorf("invalid status: %d", status)
	}
	return s.ticketRepo.UpdateStatus(ctx, id, status)
}

// AssignTicket 派单
func (s *TicketService) AssignTicket(ctx context.Context, id int64, assigneeID string) error {
	if assigneeID == "" {
		return fmt.Errorf("assignee_id is required")
	}
	return s.ticketRepo.Assign(ctx, id, assigneeID)
}

// AddReply 添加工单回复
func (s *TicketService) AddReply(ctx context.Context, rp *model.TicketReply) error {
	// 校验工单存在
	if _, err := s.ticketRepo.Get(ctx, rp.TicketID); err != nil {
		return err
	}
	return s.replyRepo.Create(ctx, rp)
}

// ListReplies 查询工单回复
func (s *TicketService) ListReplies(ctx context.Context, ticketID int64) ([]model.TicketReply, error) {
	return s.replyRepo.ListByTicket(ctx, ticketID)
}

// RateTicket 满意度评价
func (s *TicketService) RateTicket(ctx context.Context, id int64, score int32) error {
	if score < 1 || score > 5 {
		return fmt.Errorf("satisfaction score must be 1-5, got %d", score)
	}
	t, err := s.ticketRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if t.Status != StatusResolved {
		return fmt.Errorf("only resolved ticket can be rated, current=%d", t.Status)
	}
	return s.ticketRepo.SetSatisfaction(ctx, id, score)
}
