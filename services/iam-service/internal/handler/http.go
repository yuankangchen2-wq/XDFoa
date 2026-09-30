package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/iam-service/internal/auth"
	"github.com/oa-portal/iam-service/internal/model"
	"github.com/oa-portal/iam-service/internal/service"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

func fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: msg})
}

// ---------- 认证 ----------
type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler { return &AuthHandler{svc: s} }

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/auth")
	{
		g.POST("/login", h.Login)
		g.POST("/refresh", h.Refresh)
		g.POST("/verify", h.Verify)
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		TenantID string `json:"tenant_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	access, refresh, user, err := h.svc.Login(c.Request.Context(), req.Username, req.Password, req.TenantID)
	if err != nil {
		fail(c, 401, err.Error())
		return
	}
	ok(c, gin.H{"access_token": access, "refresh_token": refresh, "user": user})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	access, refresh, user, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		fail(c, 401, err.Error())
		return
	}
	ok(c, gin.H{"access_token": access, "refresh_token": refresh, "user": user})
}

func (h *AuthHandler) Verify(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	claims, perms, err := h.svc.Verify(c.Request.Context(), req.Token)
	if err != nil {
		fail(c, 401, "invalid token")
		return
	}
	ok(c, gin.H{
		"valid":       true,
		"user_id":     claims.UserID,
		"username":    claims.Username,
		"tenant_id":   claims.TenantID,
		"roles":       claims.Roles,
		"permissions": perms,
		"org_path":    claims.OrgPath,
	})
}

// ---------- 用户 ----------
type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler { return &UserHandler{svc: s} }

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/users")
	{
		g.GET("", h.List)
		g.GET("/:id", h.Get)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	}
}

func (h *UserHandler) List(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	result, err := h.svc.List(c.Request.Context(), tenantID, keyword, page, pageSize)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, result)
}

func (h *UserHandler) Get(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	u, err := h.svc.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		fail(c, 404, "not found")
		return
	}
	ok(c, u)
}

func (h *UserHandler) Create(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	var req struct {
		Username string  `json:"username"`
		Password string  `json:"password"`
		RealName string  `json:"real_name"`
		Email    string  `json:"email"`
		Phone    string  `json:"phone"`
		OrgID    int64   `json:"org_id"`
		RoleIDs  []int64 `json:"role_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	u := &model.User{
		TenantID:     tenantID,
		Username:     req.Username,
		PasswordHash: req.Password,
		RealName:     req.RealName,
		Email:        req.Email,
		Phone:        req.Phone,
		OrgID:        req.OrgID,
		Status:       1,
	}
	if err := h.svc.Create(c.Request.Context(), u, req.RoleIDs); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, u)
}

func (h *UserHandler) Update(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req struct {
		RealName string  `json:"real_name"`
		Email    string  `json:"email"`
		Phone    string  `json:"phone"`
		OrgID    int64   `json:"org_id"`
		Status   int32   `json:"status"`
		RoleIDs  []int64 `json:"role_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	u := &model.User{
		ID:       id,
		TenantID: tenantID,
		RealName: req.RealName,
		Email:    req.Email,
		Phone:    req.Phone,
		OrgID:    req.OrgID,
		Status:   req.Status,
	}
	if err := h.svc.Update(c.Request.Context(), u, req.RoleIDs); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, u)
}

func (h *UserHandler) Delete(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), id, tenantID); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

// ---------- 角色 ----------
type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(s *service.RoleService) *RoleHandler { return &RoleHandler{svc: s} }

func (h *RoleHandler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/roles")
	{
		g.GET("", h.List)
		g.POST("", h.Create)
		g.DELETE("/:id", h.Delete)
	}
}

func (h *RoleHandler) List(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	result, err := h.svc.List(c.Request.Context(), tenantID, page, pageSize)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, result)
}

func (h *RoleHandler) Create(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		fail(c, 400, err.Error())
		return
	}
	role.TenantID = tenantID
	if err := h.svc.Create(c.Request.Context(), &role); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, role)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), id, tenantID); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

// ---------- 权限 + 策略 ----------
type PermissionHandler struct {
	svc *service.PermissionService
}

func NewPermissionHandler(s *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{svc: s}
}

func (h *PermissionHandler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/permissions")
	{
		g.GET("", h.List)
	}
	pg := r.Group("/policies")
	{
		pg.GET("", h.ListPolicies)
		pg.POST("", h.AddPolicy)
		pg.DELETE("", h.RemovePolicy)
	}
}

func (h *PermissionHandler) List(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	list, err := h.svc.List(c.Request.Context(), tenantID)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, list)
}

func (h *PermissionHandler) ListPolicies(c *gin.Context) {
	roleCode := c.Query("role_code")
	policies := h.svc.ListPolicies(c.Request.Context(), roleCode)
	ok(c, policies)
}

func (h *PermissionHandler) AddPolicy(c *gin.Context) {
	var req struct {
		RoleCode string `json:"role_code"`
		PermCode string `json:"perm_code"`
		Action   string `json:"action"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	if err := h.svc.AddPolicy(c.Request.Context(), req.RoleCode, req.PermCode, req.Action); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

func (h *PermissionHandler) RemovePolicy(c *gin.Context) {
	var req struct {
		RoleCode string `json:"role_code"`
		PermCode string `json:"perm_code"`
		Action   string `json:"action"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	if err := h.svc.RemovePolicy(c.Request.Context(), req.RoleCode, req.PermCode, req.Action); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

// AuthMiddleware JWT 鉴权中间件
func AuthMiddleware(jwtMgr *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		claims, err := jwtMgr.Parse(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{Code: 401, Message: "unauthorized"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("tenant_id", claims.TenantID)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}
