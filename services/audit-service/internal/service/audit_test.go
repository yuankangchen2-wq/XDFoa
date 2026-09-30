package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/audit-service/internal/event"
	"github.com/oa-portal/audit-service/internal/model"
	"github.com/oa-portal/audit-service/internal/repo"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*AuditService, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.AutoMigrate(&model.AuditLog{})
	auditRepo := repo.NewAuditLogRepo(db)
	svc := NewAuditService(auditRepo, event.NoopPublisher{})
	return svc, db
}

// 1. 记录审计日志：状态默认为成功
func TestRecordAudit(t *testing.T) {
	svc, _ := setupTest(t)
	log := &model.AuditLog{
		TenantID: "t1", UserID: "u1", Username: "admin",
		Action: "CREATE", Module: "user", TargetID: "100", TargetName: "张三",
		RequestParams: `{"name":"张三"}`, IP: "127.0.0.1", DurationMs: 15,
	}
	if err := svc.Record(context.Background(), log); err != nil {
		t.Fatalf("record: %v", err)
	}
	if log.Status != StatusSuccess {
		t.Errorf("status = %d, want success(%d)", log.Status, StatusSuccess)
	}
	if log.ID == 0 {
		t.Error("expected id to be set")
	}
}

// 2. 按模块筛选
func TestQueryByModule(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()
	svc.Record(ctx, &model.AuditLog{TenantID: "t1", Module: "user", Action: "CREATE"})
	svc.Record(ctx, &model.AuditLog{TenantID: "t1", Module: "order", Action: "CREATE"})
	svc.Record(ctx, &model.AuditLog{TenantID: "t1", Module: "user", Action: "UPDATE"})

	res, err := svc.Query(ctx, repo.AuditQuery{TenantID: "t1", Module: "user"}, repo.Page{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if res.Total != 2 {
		t.Errorf("total = %d, want 2", res.Total)
	}
}

// 3. 按操作类型筛选
func TestQueryByAction(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()
	svc.Record(ctx, &model.AuditLog{TenantID: "t1", Module: "user", Action: "CREATE"})
	svc.Record(ctx, &model.AuditLog{TenantID: "t1", Module: "user", Action: "DELETE"})

	res, _ := svc.Query(ctx, repo.AuditQuery{TenantID: "t1", Action: "DELETE"}, repo.Page{Page: 1, PageSize: 10})
	if res.Total != 1 {
		t.Errorf("total = %d, want 1", res.Total)
	}
}

// 4. 按状态筛选（失败的日志）
func TestQueryByStatus(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()
	svc.Record(ctx, &model.AuditLog{TenantID: "t1", Action: "LOGIN", Status: StatusSuccess})
	svc.Record(ctx, &model.AuditLog{TenantID: "t1", Action: "LOGIN", Status: StatusFail})

	res, _ := svc.Query(ctx, repo.AuditQuery{TenantID: "t1", Status: StatusFail}, repo.Page{Page: 1, PageSize: 10})
	if res.Total != 1 {
		t.Errorf("total = %d, want 1", res.Total)
	}
}

// 5. 模块统计
func TestModuleStats(t *testing.T) {
	svc, _ := setupTest(t)
	ctx := context.Background()
	svc.Record(ctx, &model.AuditLog{TenantID: "t1", Module: "user", Action: "CREATE"})
	svc.Record(ctx, &model.AuditLog{TenantID: "t1", Module: "user", Action: "UPDATE"})
	svc.Record(ctx, &model.AuditLog{TenantID: "t1", Module: "order", Action: "CREATE"})

	rows, err := svc.ModuleStats(ctx, "t1")
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	var userCount int64
	for _, r := range rows {
		if r.Module == "user" {
			userCount = r.Count
		}
	}
	if userCount != 2 {
		t.Errorf("user module count = %d, want 2", userCount)
	}
}
