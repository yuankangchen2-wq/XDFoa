package repo

import (
	"context"

	"github.com/oa-portal/iam-service/internal/model"
	"gorm.io/gorm"
)

type PageQuery struct {
	Page     int
	PageSize int
}

type PageResult[T any] struct {
	Items    []T
	Total    int64
	Page     int
	PageSize int
}

func Paginate[T any](db *gorm.DB, query *gorm.DB, pq PageQuery) (*PageResult[T], error) {
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	var items []T
	offset := (pq.Page - 1) * pq.PageSize
	if err := query.Offset(offset).Limit(pq.PageSize).Find(&items).Error; err != nil {
		return nil, err
	}
	return &PageResult[T]{Items: items, Total: total, Page: pq.Page, PageSize: pq.PageSize}, nil
}

type ctxKey struct{}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKey{}, tx)
}

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(ctxKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}

// ---------- 用户 ----------
type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) DB() *gorm.DB { return r.db }

func (r *UserRepo) Create(ctx context.Context, u *model.User) error {
	return GetDB(ctx, r.db).Create(u).Error
}

func (r *UserRepo) Update(ctx context.Context, u *model.User) error {
	return GetDB(ctx, r.db).Save(u).Error
}

func (r *UserRepo) GetByID(ctx context.Context, id int64, tenantID string) (*model.User, error) {
	var u model.User
	err := GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username, tenantID string) (*model.User, error) {
	var u model.User
	err := GetDB(ctx, r.db).Where("username = ? AND tenant_id = ?", username, tenantID).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context, tenantID, keyword string, pq PageQuery) (*PageResult[model.User], error) {
	q := GetDB(ctx, r.db).Model(&model.User{}).Where("tenant_id = ?", tenantID)
	if keyword != "" {
		q = q.Where("username LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	return Paginate[model.User](r.db, q, pq)
}

func (r *UserRepo) Delete(ctx context.Context, id int64, tenantID string) error {
	return GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.User{}).Error
}

// ---------- 角色 ----------
type RoleRepo struct{ db *gorm.DB }

func NewRoleRepo(db *gorm.DB) *RoleRepo { return &RoleRepo{db: db} }

func (r *RoleRepo) Create(ctx context.Context, role *model.Role) error {
	return GetDB(ctx, r.db).Create(role).Error
}

func (r *RoleRepo) GetByID(ctx context.Context, id int64, tenantID string) (*model.Role, error) {
	var role model.Role
	err := GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) GetByCode(ctx context.Context, roleCode, tenantID string) (*model.Role, error) {
	var role model.Role
	err := GetDB(ctx, r.db).Where("role_code = ? AND tenant_id = ?", roleCode, tenantID).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) List(ctx context.Context, tenantID string, pq PageQuery) (*PageResult[model.Role], error) {
	q := GetDB(ctx, r.db).Model(&model.Role{}).Where("tenant_id = ?", tenantID)
	return Paginate[model.Role](r.db, q, pq)
}

func (r *RoleRepo) Delete(ctx context.Context, id int64, tenantID string) error {
	return GetDB(ctx, r.db).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Role{}).Error
}

// ---------- 用户角色 ----------
type UserRoleRepo struct{ db *gorm.DB }

func NewUserRoleRepo(db *gorm.DB) *UserRoleRepo { return &UserRoleRepo{db: db} }

func (r *UserRoleRepo) Add(ctx context.Context, ur *model.UserRole) error {
	return GetDB(ctx, r.db).Create(ur).Error
}

func (r *UserRoleRepo) RemoveByUser(ctx context.Context, userID int64) error {
	return GetDB(ctx, r.db).Where("user_id = ?", userID).Delete(&model.UserRole{}).Error
}

func (r *UserRoleRepo) GetRoleIDsByUser(ctx context.Context, userID int64) ([]int64, error) {
	var ids []int64
	err := GetDB(ctx, r.db).Model(&model.UserRole{}).Where("user_id = ?", userID).Pluck("role_id", &ids).Error
	return ids, err
}

// ---------- 权限 ----------
type PermissionRepo struct{ db *gorm.DB }

func NewPermissionRepo(db *gorm.DB) *PermissionRepo { return &PermissionRepo{db: db} }

func (r *PermissionRepo) List(ctx context.Context, tenantID string) ([]model.Permission, error) {
	var list []model.Permission
	err := GetDB(ctx, r.db).Where("tenant_id = ?", tenantID).Find(&list).Error
	return list, err
}

// ---------- 审计日志 ----------
type AuditRepo struct{ db *gorm.DB }

func NewAuditRepo(db *gorm.DB) *AuditRepo { return &AuditRepo{db: db} }

func (r *AuditRepo) Create(ctx context.Context, log *model.AuditLog) error {
	return GetDB(ctx, r.db).Create(log).Error
}
