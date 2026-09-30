package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/oa-portal/iam-service/internal/auth"
	"github.com/oa-portal/iam-service/internal/model"
	"github.com/oa-portal/iam-service/internal/repo"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	db.AutoMigrate(
		&model.User{}, &model.Role{}, &model.UserRole{},
		&model.Permission{}, &model.CasbinRule{}, &model.AuditLog{},
	)
	return db
}

func TestAuthService_Login(t *testing.T) {
	db := setupTestDB(t)

	jwtMgr := auth.NewJWTManager("test-secret", 900, 604800)
	casbinMgr, err := auth.NewCasbinManager(db, "../../configs/casbin_model.conf")
	if err != nil {
		t.Fatalf("casbin: %v", err)
	}

	userRepo := repo.NewUserRepo(db)
	roleRepo := repo.NewRoleRepo(db)
	urRepo := repo.NewUserRoleRepo(db)
	auditRepo := repo.NewAuditRepo(db)

	authSvc := NewAuthService(userRepo, roleRepo, urRepo, auditRepo, jwtMgr, casbinMgr)
	userSvc := NewUserService(userRepo, roleRepo, urRepo, auditRepo, casbinMgr)

	ctx := context.Background()

	// 创建角色
	role := &model.Role{TenantID: "t1", RoleCode: "admin", RoleName: "管理员", DataScope: "all"}
	roleRepo.Create(ctx, role)

	// 创建用户
	user := &model.User{
		TenantID:     "t1",
		Username:     "testuser",
		PasswordHash: "password123",
		RealName:     "测试用户",
	}
	if err := userSvc.Create(ctx, user, []int64{role.ID}); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// 登录
	access, refresh, u, err := authSvc.Login(ctx, "testuser", "password123", "t1")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if access == "" || refresh == "" {
		t.Error("tokens should not be empty")
	}
	if u.Username != "testuser" {
		t.Errorf("username = %s", u.Username)
	}

	// 验证 token
	claims, perms, err := authSvc.Verify(ctx, access)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if claims.Username != "testuser" {
		t.Errorf("claims username = %s", claims.Username)
	}
	_ = perms

	// 错误密码
	_, _, _, err = authSvc.Login(ctx, "testuser", "wrong", "t1")
	if err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestPermissionService_Enforce(t *testing.T) {
	db := setupTestDB(t)

	casbinMgr, err := auth.NewCasbinManager(db, "../../configs/casbin_model.conf")
	if err != nil {
		t.Fatalf("casbin: %v", err)
	}

	permSvc := NewPermissionService(repo.NewPermissionRepo(db), casbinMgr)
	ctx := context.Background()

	// 添加策略: admin 角色可以 read customer
	if err := permSvc.AddPolicy(ctx, "admin", "customer", "read"); err != nil {
		t.Fatalf("add policy: %v", err)
	}

	// 给用户分配角色
	userID := int64(1)
	casbinMgr.AddRoleForUser(strconv.FormatInt(userID, 10), "admin")

	// 校验权限
	ok, err := permSvc.Enforce(ctx, userID, "customer", "read")
	if err != nil {
		t.Fatalf("enforce: %v", err)
	}
	if !ok {
		t.Error("should have permission")
	}

	// 无权限的操作
	ok, _ = permSvc.Enforce(ctx, userID, "customer", "delete")
	if ok {
		t.Error("should not have delete permission")
	}
}
