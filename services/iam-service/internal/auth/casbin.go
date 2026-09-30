package auth

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Casbin 权限管理器
type CasbinManager struct {
	enforcer *casbin.Enforcer
}

// NewCasbinManager 初始化 Casbin，使用 GORM 适配器存储策略
func NewCasbinManager(db *gorm.DB, modelPath string) (*CasbinManager, error) {
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, nil, "iam_casbin_rule")
	if err != nil {
		return nil, err
	}
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, err
	}
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	return &CasbinManager{enforcer: enforcer}, nil
}

// Enforce 权限校验
func (c *CasbinManager) Enforce(sub, obj, act string) (bool, error) {
	return c.enforcer.Enforce(sub, obj, act)
}

// AddPolicy 添加策略 p = sub, obj, act
func (c *CasbinManager) AddPolicy(roleCode, permCode, action string) (bool, error) {
	return c.enforcer.AddPolicy(roleCode, permCode, action)
}

// RemovePolicy 删除策略
func (c *CasbinManager) RemovePolicy(roleCode, permCode, action string) (bool, error) {
	return c.enforcer.RemovePolicy(roleCode, permCode, action)
}

// AddRoleForUser 给用户分配角色 g = user, role
func (c *CasbinManager) AddRoleForUser(userID, roleCode string) (bool, error) {
	return c.enforcer.AddGroupingPolicy(userID, roleCode)
}

// RemoveRoleForUser 移除用户角色
func (c *CasbinManager) RemoveRoleForUser(userID, roleCode string) (bool, error) {
	return c.enforcer.RemoveGroupingPolicy(userID, roleCode)
}

// GetRolesForUser 获取用户的角色
func (c *CasbinManager) GetRolesForUser(userID string) ([]string, error) {
	return c.enforcer.GetRolesForUser(userID)
}

// GetPermissionsForRole 获取角色的权限
func (c *CasbinManager) GetPermissionsForRole(roleCode string) [][]string {
	return c.enforcer.GetPermissionsForUser(roleCode)
}

// ReloadPolicy 重新加载策略
func (c *CasbinManager) ReloadPolicy() error {
	return c.enforcer.LoadPolicy()
}
