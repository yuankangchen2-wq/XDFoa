package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/oa-portal/iam-service/internal/auth"
	"github.com/oa-portal/iam-service/internal/model"
	"github.com/oa-portal/iam-service/internal/repo"
	"gorm.io/gorm"
)

// ---------- 认证服务 ----------
type AuthService struct {
	userRepo  *repo.UserRepo
	roleRepo  *repo.RoleRepo
	urRepo    *repo.UserRoleRepo
	auditRepo *repo.AuditRepo
	jwt       *auth.JWTManager
	casbin    *auth.CasbinManager
}

func NewAuthService(ur *repo.UserRepo, rr *repo.RoleRepo, urr *repo.UserRoleRepo, ar *repo.AuditRepo, j *auth.JWTManager, c *auth.CasbinManager) *AuthService {
	return &AuthService{userRepo: ur, roleRepo: rr, urRepo: urr, auditRepo: ar, jwt: j, casbin: c}
}

// Login 登录
func (s *AuthService) Login(ctx context.Context, username, password, tenantID string) (string, string, *model.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username, tenantID)
	if err != nil {
		return "", "", nil, errors.New("用户名或密码错误")
	}
	if user.Status != 1 {
		return "", "", nil, errors.New("账号已禁用")
	}
	if !auth.CheckPassword(password, user.PasswordHash) {
		s.recordAudit(ctx, user.ID, user.Username, tenantID, "login", "failed")
		return "", "", nil, errors.New("用户名或密码错误")
	}

	// 获取用户角色
	roleCodes, _ := s.casbin.GetRolesForUser(strconv.FormatInt(user.ID, 10))

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Username, user.TenantID, roleCodes, user.OrgPath)
	if err != nil {
		return "", "", nil, err
	}
	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID, user.Username, user.TenantID)
	if err != nil {
		return "", "", nil, err
	}

	s.recordAudit(ctx, user.ID, user.Username, tenantID, "login", "success")
	return accessToken, refreshToken, user, nil
}

// Refresh 刷新 token
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, *model.User, error) {
	claims, err := s.jwt.Parse(refreshToken)
	if err != nil || claims.Type != "refresh" {
		return "", "", nil, errors.New("invalid refresh token")
	}
	user, err := s.userRepo.GetByID(ctx, claims.UserID, claims.TenantID)
	if err != nil {
		return "", "", nil, errors.New("user not found")
	}
	roleCodes, _ := s.casbin.GetRolesForUser(strconv.FormatInt(user.ID, 10))
	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Username, user.TenantID, roleCodes, user.OrgPath)
	if err != nil {
		return "", "", nil, err
	}
	newRefresh, err := s.jwt.GenerateRefreshToken(user.ID, user.Username, user.TenantID)
	if err != nil {
		return "", "", nil, err
	}
	return accessToken, newRefresh, user, nil
}

// Verify 验证 token（供其他服务调用）
func (s *AuthService) Verify(ctx context.Context, tokenStr string) (*auth.Claims, []string, error) {
	claims, err := s.jwt.Parse(tokenStr)
	if err != nil {
		return nil, nil, err
	}
	// 获取权限码列表
	perms := []string{}
	for _, role := range claims.Roles {
		policies := s.casbin.GetPermissionsForRole(role)
		for _, p := range policies {
			if len(p) >= 2 {
				perms = append(perms, p[1])
			}
		}
	}
	return claims, perms, nil
}

func (s *AuthService) recordAudit(ctx context.Context, userID int64, username, tenantID, action, status string) {
	st := int32(0)
	if status == "success" {
		st = 1
	}
	s.auditRepo.Create(ctx, &model.AuditLog{
		TenantID: tenantID,
		UserID:   userID,
		Username: username,
		Module:   "auth",
		Action:   action,
		Status:   st,
	})
}

// ---------- 用户服务 ----------
type UserService struct {
	userRepo  *repo.UserRepo
	roleRepo  *repo.RoleRepo
	urRepo    *repo.UserRoleRepo
	auditRepo *repo.AuditRepo
	casbin    *auth.CasbinManager
}

func NewUserService(ur *repo.UserRepo, rr *repo.RoleRepo, urr *repo.UserRoleRepo, ar *repo.AuditRepo, c *auth.CasbinManager) *UserService {
	return &UserService{userRepo: ur, roleRepo: rr, urRepo: urr, auditRepo: ar, casbin: c}
}

func (s *UserService) Create(ctx context.Context, u *model.User, roleIDs []int64) error {
	hash, err := auth.HashPassword(u.PasswordHash)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	// 事务：创建用户 + 分配角色
	return s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		txCtx := repo.WithTx(ctx, tx)
		if err := s.userRepo.Create(txCtx, u); err != nil {
			return err
		}
		return s.assignRoles(txCtx, u.ID, u.TenantID, roleIDs)
	})
}

func (s *UserService) assignRoles(ctx context.Context, userID int64, tenantID string, roleIDs []int64) error {
	if err := s.urRepo.RemoveByUser(ctx, userID); err != nil {
		return err
	}
	for _, rid := range roleIDs {
		role, err := s.roleRepo.GetByID(ctx, rid, tenantID)
		if err != nil {
			continue
		}
		s.urRepo.Add(ctx, &model.UserRole{TenantID: tenantID, UserID: userID, RoleID: rid})
		s.casbin.AddRoleForUser(strconv.FormatInt(userID, 10), role.RoleCode)
	}
	return nil
}

func (s *UserService) Update(ctx context.Context, u *model.User, roleIDs []int64) error {
	u.UpdatedAt = time.Now()
	return s.userRepo.DB().Transaction(func(tx *gorm.DB) error {
		txCtx := repo.WithTx(ctx, tx)
		if err := s.userRepo.Update(txCtx, u); err != nil {
			return err
		}
		return s.assignRoles(txCtx, u.ID, u.TenantID, roleIDs)
	})
}

func (s *UserService) Get(ctx context.Context, id int64, tenantID string) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id, tenantID)
}

func (s *UserService) List(ctx context.Context, tenantID, keyword string, page, pageSize int) (*repo.PageResult[model.User], error) {
	return s.userRepo.List(ctx, tenantID, keyword, repo.PageQuery{Page: page, PageSize: pageSize})
}

func (s *UserService) Delete(ctx context.Context, id int64, tenantID string) error {
	return s.userRepo.Delete(ctx, id, tenantID)
}

// ---------- 角色服务 ----------
type RoleService struct {
	roleRepo  *repo.RoleRepo
	auditRepo *repo.AuditRepo
}

func NewRoleService(rr *repo.RoleRepo, ar *repo.AuditRepo) *RoleService {
	return &RoleService{roleRepo: rr, auditRepo: ar}
}

func (s *RoleService) Create(ctx context.Context, role *model.Role) error {
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()
	return s.roleRepo.Create(ctx, role)
}

func (s *RoleService) List(ctx context.Context, tenantID string, page, pageSize int) (*repo.PageResult[model.Role], error) {
	return s.roleRepo.List(ctx, tenantID, repo.PageQuery{Page: page, PageSize: pageSize})
}

func (s *RoleService) Delete(ctx context.Context, id int64, tenantID string) error {
	return s.roleRepo.Delete(ctx, id, tenantID)
}

// ---------- 权限服务 ----------
type PermissionService struct {
	permRepo *repo.PermissionRepo
	casbin   *auth.CasbinManager
}

func NewPermissionService(pr *repo.PermissionRepo, c *auth.CasbinManager) *PermissionService {
	return &PermissionService{permRepo: pr, casbin: c}
}

func (s *PermissionService) List(ctx context.Context, tenantID string) ([]model.Permission, error) {
	return s.permRepo.List(ctx, tenantID)
}

// ListPolicies 获取角色的权限策略
func (s *PermissionService) ListPolicies(ctx context.Context, roleCode string) [][]string {
	return s.casbin.GetPermissionsForRole(roleCode)
}

// AddPolicy 添加策略
func (s *PermissionService) AddPolicy(ctx context.Context, roleCode, permCode, action string) error {
	_, err := s.casbin.AddPolicy(roleCode, permCode, action)
	return err
}

// RemovePolicy 删除策略
func (s *PermissionService) RemovePolicy(ctx context.Context, roleCode, permCode, action string) error {
	_, err := s.casbin.RemovePolicy(roleCode, permCode, action)
	return err
}

// Enforce 权限校验
func (s *PermissionService) Enforce(ctx context.Context, userID int64, obj, act string) (bool, error) {
	return s.casbin.Enforce(strconv.FormatInt(userID, 10), obj, act)
}
