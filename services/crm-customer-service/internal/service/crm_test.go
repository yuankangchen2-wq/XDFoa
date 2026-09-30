package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/crm-customer-service/internal/event"
	"github.com/oa-portal/crm-customer-service/internal/model"
	"github.com/oa-portal/crm-customer-service/internal/repo"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) *CrmService {
	t.Helper()
	dbFile := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&model.Customer{}, &model.Contact{}, &model.FollowUp{})

	return NewCrmService(
		repo.NewCustomerRepo(db),
		repo.NewContactRepo(db),
		repo.NewFollowUpRepo(db),
		event.NoopPublisher{},
	)
}

// 测试创建客户
func TestCreateCustomer(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	c := &model.Customer{
		TenantID: "t1", Name: "测试客户", ContactPerson: "张三",
		Phone: "13800138000", Level: "A", Status: 1,
	}
	if err := svc.CreateCustomer(ctx, c); err != nil {
		t.Fatal(err)
	}
	if c.ID == 0 {
		t.Fatal("customer id should not be 0")
	}
	if c.Status != 1 {
		t.Fatalf("expected status=1, got %d", c.Status)
	}
}

// 测试查询客户列表（按分级筛选）
func TestListCustomers(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	svc.CreateCustomer(ctx, &model.Customer{TenantID: "t1", Name: "客户A", Level: "A", Status: 1})
	svc.CreateCustomer(ctx, &model.Customer{TenantID: "t1", Name: "客户B", Level: "B", Status: 1})
	svc.CreateCustomer(ctx, &model.Customer{TenantID: "t1", Name: "客户A2", Level: "A", Status: 1})

	// 查 A 级客户
	res, err := svc.ListCustomers(ctx, "t1", "A", 0, repo.Page{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 2 {
		t.Fatalf("expected 2 A-level customers, got %d", res.Total)
	}
}

// 测试添加联系人
func TestAddContact(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	c := &model.Customer{TenantID: "t1", Name: "客户A", Status: 1}
	svc.CreateCustomer(ctx, c)

	ct := &model.Contact{
		CustomerID: c.ID, Name: "李四", Position: "经理",
		Phone: "13900139000", IsPrimary: 1,
	}
	if err := svc.AddContact(ctx, ct); err != nil {
		t.Fatal(err)
	}

	contacts, _ := svc.ListContacts(ctx, c.ID)
	if len(contacts) != 1 {
		t.Fatalf("expected 1 contact, got %d", len(contacts))
	}
}

// 测试添加跟进记录
func TestAddFollowUp(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	c := &model.Customer{TenantID: "t1", Name: "客户A", Status: 1}
	svc.CreateCustomer(ctx, c)

	f := &model.FollowUp{
		CustomerID: c.ID, Type: "电话", Content: "首次沟通，客户有意向",
		UserID: "sales001", FollowUpAt: time.Now(),
	}
	if err := svc.AddFollowUp(ctx, f); err != nil {
		t.Fatal(err)
	}

	followUps, _ := svc.ListFollowUps(ctx, c.ID)
	if len(followUps) != 1 {
		t.Fatalf("expected 1 follow-up, got %d", len(followUps))
	}

	// GetCustomer 应该包含跟进记录
	got, _ := svc.GetCustomer(ctx, c.ID)
	if len(got.FollowUps) != 1 {
		t.Fatalf("expected 1 follow-up in customer, got %d", len(got.FollowUps))
	}
}

// 测试删除客户
func TestDeleteCustomer(t *testing.T) {
	svc := setupTest(t)
	ctx := context.Background()

	c := &model.Customer{TenantID: "t1", Name: "客户A", Status: 1}
	svc.CreateCustomer(ctx, c)

	if err := svc.DeleteCustomer(ctx, c.ID); err != nil {
		t.Fatal(err)
	}
	_, err := svc.GetCustomer(ctx, c.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}
