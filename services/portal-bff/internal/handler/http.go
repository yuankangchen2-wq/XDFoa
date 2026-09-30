package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/portal-bff/internal/service"
)

type Handler struct {
	svc *service.BFFService
}

func NewHandler(svc *service.BFFService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up", "service": "portal-bff"})
	})
	api := r.Group("/api/bff")
	{
		api.GET("/health", h.Health)
		api.GET("/menus", h.Menus)
		api.GET("/me", h.Me)
	}
}

// Health 聚合各服务健康状态
func (h *Handler) Health(c *gin.Context) {
	results := h.svc.CheckAllHealth(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"items": results})
}

// Menus 菜单
func (h *Handler) Menus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"items": h.svc.GetMenus()})
}

// Me 当前用户信息（从 header 中提取，简化版）
func (h *Handler) Me(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	username := c.GetHeader("X-Username")
	if userID == "" {
		userID = "anonymous"
	}
	if username == "" {
		username = "匿名用户"
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
		"roles":    []string{"admin"},
	})
}
