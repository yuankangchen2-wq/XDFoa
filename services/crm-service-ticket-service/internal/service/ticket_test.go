package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/crm-service-ticket-service/internal/event"
	"github.com/oa-portal/crm-service-ticket-service/internal/model"
	"github.com/oa-portal/crm-service-ticket-service/internal/repo"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) *TicketService {
	t.Helper()
	dbFile := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&model.ServiceTicket{}, &model.TicketReply{})

	return NewTicketService(
		repo.NewTicketRepo(db),
		repo.NewReplyRepo(db),
		event.NoopPublisher{},
	)
}

// 测试创建工单
func TestCreateTicket(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	tk := &model.ServiceTicket{
		TenantID: "t1", CustomerID: 1, Title: "系统无法登录",
		Description: "用户反馈无法登录系统", Priority: 3,
	}
	if err := svc.CreateTicket(ctx, tk); err != nil {
		t.Fatal(err)
	}
	if tk.TicketNo == "" {
		t.Fatal("ticket_no should not be empty")
	}
	if tk.Status != StatusPending {
		t.Fatalf("expected status=1, got %d", tk.Status)
	}
}

// 测试派单
func TestAssignTicket(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	tk := &model.ServiceTicket{TenantID: "t1", CustomerID: 1, Title: "测试工单", Priority: 2}
	svc.CreateTicket(ctx, tk)

	if err := svc.AssignTicket(ctx, tk.ID, "agent001"); err != nil {
		t.Fatal(err)
	}
	updated, _ := svc.GetTicket(ctx, tk.ID)
	if updated.AssigneeID != "agent001" {
		t.Fatalf("expected assignee=agent001, got %s", updated.AssigneeID)
	}
	if updated.Status != StatusProcessing {
		t.Fatalf("expected status=2 after assign, got %d", updated.Status)
	}
}

// 测试添加回复
func TestAddReply(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	tk := &model.ServiceTicket{TenantID: "t1", CustomerID: 1, Title: "测试工单", Priority: 2}
	svc.CreateTicket(ctx, tk)

	rp := &model.TicketReply{
		TicketID: tk.ID, UserID: "agent001", Content: "已收到，正在处理",
	}
	if err := svc.AddReply(ctx, rp); err != nil {
		t.Fatal(err)
	}

	replies, _ := svc.ListReplies(ctx, tk.ID)
	if len(replies) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(replies))
	}

	// GetTicket 应包含回复
	got, _ := svc.GetTicket(ctx, tk.ID)
	if len(got.Replies) != 1 {
		t.Fatalf("expected 1 reply in ticket, got %d", len(got.Replies))
	}
}

// 测试满意度评价（仅已解决工单可评价）
func TestRateTicket(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	tk := &model.ServiceTicket{TenantID: "t1", CustomerID: 1, Title: "测试工单", Priority: 2}
	svc.CreateTicket(ctx, tk)

	// 未解决不能评价
	if err := svc.RateTicket(ctx, tk.ID, 5); err == nil {
		t.Fatal("expected error when rating non-resolved ticket")
	}

	// 标记为已解决
	svc.UpdateStatus(ctx, tk.ID, StatusResolved)
	if err := svc.RateTicket(ctx, tk.ID, 5); err != nil {
		t.Fatal(err)
	}

	updated, _ := svc.GetTicket(ctx, tk.ID)
	if updated.Satisfaction != 5 {
		t.Fatalf("expected satisfaction=5, got %d", updated.Satisfaction)
	}

	// 评分超出范围
	if err := svc.RateTicket(ctx, tk.ID, 6); err == nil {
		t.Fatal("expected error for score=6")
	}
}

// 测试按状态筛选
func TestListByStatus(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	tk1 := &model.ServiceTicket{TenantID: "t1", CustomerID: 1, Title: "工单A", Priority: 2}
	tk2 := &model.ServiceTicket{TenantID: "t1", CustomerID: 1, Title: "工单B", Priority: 2}
	svc.CreateTicket(ctx, tk1)
	svc.CreateTicket(ctx, tk2)
	svc.UpdateStatus(ctx, tk2.ID, StatusResolved)

	// 查待处理
	res, err := svc.ListTickets(ctx, "t1", StatusPending, 0, repo.Page{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 1 {
		t.Fatalf("expected 1 pending ticket, got %d", res.Total)
	}
}
