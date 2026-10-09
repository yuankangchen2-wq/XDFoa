package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/portal-bff/internal/model"
	"github.com/oa-portal/portal-bff/internal/service"
)

type Handler struct {
	svc *service.BFFService
	p   *service.PortalService
}

func NewHandler(svc *service.BFFService, p *service.PortalService) *Handler {
	return &Handler{svc: svc, p: p}
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

		// 员工/邮箱搜索
		api.GET("/employees", h.SearchEmployees)

		// 审批
		api.GET("/approvals", h.ListApprovals)
		api.GET("/approvals/counts", h.ApprovalCounts)
		api.POST("/approvals", h.CreateApproval)

		// 日程
		api.GET("/schedules", h.ListSchedules)

		// 我的常用
		api.GET("/favorites", h.ListFavorites)
		api.POST("/favorites", h.AddFavorite)
		api.DELETE("/favorites/:id", h.DeleteFavorite)
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

// Me 当前用户信息
func (h *Handler) Me(c *gin.Context) {
	// 入职日期：2024-07-25，相伴天数按当天动态计算
	joinDate := time.Date(2024, 7, 25, 0, 0, 0, 0, time.Local)
	days := int(time.Since(joinDate).Hours()/24) + 1
	if days < 0 {
		days = 0
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":    service.CurrentUserID,
		"username":   "陈元康",
		"title":      "老师",
		"join_days":  days,
		"join_date":  joinDate.Format("2006-01-02"),
		"avatar":     "",
		"roles":      []string{"admin"},
	})
}

// SearchEmployees 员工/邮箱搜索
func (h *Handler) SearchEmployees(c *gin.Context) {
	keyword := c.Query("keyword")
	list, err := h.p.SearchEmployees(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list})
}

// ListApprovals 审批列表
func (h *Handler) ListApprovals(c *gin.Context) {
	typ := c.Query("type")
	status := c.DefaultQuery("status", "todo") // todo / done / mine
	keyword := c.Query("keyword")
	list, err := h.p.ListApprovals(typ, status, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list})
}

// ApprovalCounts 审批数量
func (h *Handler) ApprovalCounts(c *gin.Context) {
	counts, err := h.p.ApprovalCounts(c.Query("type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, counts)
}

// CreateApproval 发起审批
func (h *Handler) CreateApproval(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Type      string `json:"type" binding:"required"`
		SerialNo  string `json:"serial_no"`
		Applicant string `json:"applicant"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a := &model.PortalApproval{
		Title:       req.Title,
		Type:        req.Type,
		SerialNo:    req.SerialNo,
		Applicant:   req.Applicant,
		ApplicantID: service.CurrentUserID,
		Status:      1,
	}
	if err := h.p.CreateApproval(a); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

// ListSchedules 日程
func (h *Handler) ListSchedules(c *gin.Context) {
	list, err := h.p.ListSchedules(c.Query("month"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list})
}

// ListFavorites 我的常用
func (h *Handler) ListFavorites(c *gin.Context) {
	list, err := h.p.ListFavorites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": list})
}

// AddFavorite 添加常用
func (h *Handler) AddFavorite(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Icon string `json:"icon"`
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f, err := h.p.AddFavorite(req.Name, req.Icon, req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, f)
}

// DeleteFavorite 删除常用
func (h *Handler) DeleteFavorite(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.p.DeleteFavorite(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
