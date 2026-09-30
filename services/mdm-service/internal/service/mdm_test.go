package service

import (
	"context"
	"testing"

	"github.com/oa-portal/mdm-service/internal/event"
	"github.com/oa-portal/mdm-service/internal/model"
	"github.com/oa-portal/mdm-service/internal/repo"
	_ "modernc.org/sqlite" // pure-go sqlite 驱动
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// 使用 modernc pure-go sqlite（无需 CGO）
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        ":memory:",
	}, &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.AutoMigrate(&model.Customer{}, &model.Product{}, &model.Supplier{}, &model.Organization{})
	return db
}

// ---------- 客户测试 ----------
func TestCustomerService_Create(t *testing.T) {
	db := setupTestDB(t)
	svc := NewCustomerService(repo.NewCustomerRepo(db), event.NoopPublisher{})
	ctx := context.Background()

	c := &model.Customer{
		TenantID:     "t1",
		CustomerCode: "C001",
		CustomerName: "测试客户",
		CustomerType: "enterprise",
	}
	if err := svc.Create(ctx, c); err != nil {
		t.Fatalf("create: %v", err)
	}
	if c.ID == 0 {
		t.Error("id should be set")
	}

	// 查询
	got, err := svc.Get(ctx, c.ID, "t1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.CustomerName != "测试客户" {
		t.Errorf("name = %s, want 测试客户", got.CustomerName)
	}
}

func TestCustomerService_List(t *testing.T) {
	db := setupTestDB(t)
	svc := NewCustomerService(repo.NewCustomerRepo(db), event.NoopPublisher{})
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		svc.Create(ctx, &model.Customer{
			TenantID:     "t1",
			CustomerCode: "C00" + string(rune('1'+i)),
			CustomerName: "客户" + string(rune('A'+i)),
		})
	}

	result, err := svc.List(ctx, "t1", "", "", 1, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if result.Total != 3 {
		t.Errorf("total = %d, want 3", result.Total)
	}
}

func TestCustomerService_Delete(t *testing.T) {
	db := setupTestDB(t)
	svc := NewCustomerService(repo.NewCustomerRepo(db), event.NoopPublisher{})
	ctx := context.Background()

	c := &model.Customer{TenantID: "t1", CustomerCode: "C001", CustomerName: "待删"}
	svc.Create(ctx, c)

	if err := svc.Delete(ctx, c.ID, "t1"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err := svc.Get(ctx, c.ID, "t1")
	if err == nil {
		t.Error("expected error after delete")
	}
}

// ---------- 组织测试 ----------
func TestOrganizationService_CreateWithParent(t *testing.T) {
	db := setupTestDB(t)
	svc := NewOrganizationService(repo.NewOrganizationRepo(db), event.NoopPublisher{})
	ctx := context.Background()

	root := &model.Organization{
		TenantID: "t1",
		OrgCode:  "root",
		OrgName:  "总部",
	}
	if err := svc.Create(ctx, root); err != nil {
		t.Fatalf("create root: %v", err)
	}
	if root.OrgPath != "/root" || root.OrgLevel != 1 {
		t.Errorf("root path=%s level=%d", root.OrgPath, root.OrgLevel)
	}

	child := &model.Organization{
		TenantID: "t1",
		ParentID: root.ID,
		OrgCode:  "dept",
		OrgName:  "部门",
	}
	if err := svc.Create(ctx, child); err != nil {
		t.Fatalf("create child: %v", err)
	}
	if child.OrgPath != "/root/dept" || child.OrgLevel != 2 {
		t.Errorf("child path=%s level=%d", child.OrgPath, child.OrgLevel)
	}
}
