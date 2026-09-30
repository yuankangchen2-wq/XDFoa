package model

import "time"

// 用户
type User struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	Username     string    `gorm:"size:64;uniqueIndex:uk_tenant_username" json:"username"`
	PasswordHash string    `gorm:"size:256" json:"-"`
	RealName     string    `gorm:"size:64" json:"real_name"`
	Email        string    `gorm:"size:128" json:"email"`
	Phone        string    `gorm:"size:32" json:"phone"`
	OrgID        int64     `json:"org_id"`
	OrgPath      string    `gorm:"size:512" json:"org_path"`
	Status       int32     `gorm:"default:1" json:"status"` // 1=正常 0=禁用
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string { return "iam_user" }

// 角色
type Role struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	RoleCode    string    `gorm:"size:64;uniqueIndex:uk_tenant_role" json:"role_code"`
	RoleName    string    `gorm:"size:128" json:"role_name"`
	Description string    `gorm:"size:512" json:"description"`
	DataScope   string    `gorm:"size:32;default:self" json:"data_scope"` // self/dept/dept_and_sub/all
	Status      int32     `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Role) TableName() string { return "iam_role" }

// 用户-角色关联
type UserRole struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	TenantID string `gorm:"index;size:64" json:"tenant_id"`
	UserID   int64  `json:"user_id"`
	RoleID   int64  `json:"role_id"`
}

func (UserRole) TableName() string { return "iam_user_role" }

// 权限（菜单/按钮/接口）
type Permission struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	TenantID  string `gorm:"index;size:64" json:"tenant_id"`
	PermCode  string `gorm:"size:128;uniqueIndex:uk_tenant_perm" json:"perm_code"`
	PermName  string `gorm:"size:128" json:"perm_name"`
	PermType  int32  `json:"perm_type"` // 1=菜单 2=按钮 3=接口
	Path      string `gorm:"size:256" json:"path"`
	Method    string `gorm:"size:16" json:"method"`
	Status    int32  `gorm:"default:1" json:"status"`
}

func (Permission) TableName() string { return "iam_permission" }

// Casbin 策略（角色-权限关联）
type CasbinRule struct {
	ID    int64  `gorm:"primaryKey" json:"id"`
	PType string `gorm:"size:16" json:"ptype"` // p/g
	V0    string `gorm:"size:128" json:"v0"`    // role_code or user
	V1    string `gorm:"size:256" json:"v1"`    // perm_code or role_code
	V2    string `gorm:"size:128" json:"v2"`    // action
	V3    string `gorm:"size:128" json:"v3"`
	V4    string `gorm:"size:128" json:"v4"`
	V5    string `gorm:"size:128" json:"v5"`
}

func (CasbinRule) TableName() string { return "iam_casbin_rule" }

// 审计日志
type AuditLog struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	TenantID   string    `gorm:"index;size:64" json:"tenant_id"`
	UserID     int64     `json:"user_id"`
	Username   string    `gorm:"size:64" json:"username"`
	Module     string    `gorm:"size:64" json:"module"`
	Action     string    `gorm:"size:64" json:"action"`
	Resource   string    `gorm:"size:128" json:"resource"`
	ClientIP   string    `gorm:"size:64" json:"client_ip"`
	UserAgent  string    `gorm:"size:512" json:"user_agent"`
	Detail     string    `gorm:"type:text" json:"detail"`
	Status     int32     `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AuditLog) TableName() string { return "iam_audit_log" }
