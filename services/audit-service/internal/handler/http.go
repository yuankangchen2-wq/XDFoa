package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/audit-service/internal/model"
	"github.com/oa-portal/audit-service/internal/repo"
	"github.com/oa-portal/audit-service/internal/service"
)

type Handler struct {
	svc *service.AuditService
}

func NewHandler(svc *service.AuditService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/audit")
	{
		api.POST("/logs", h.Record)
		api.GET("/logs", h.List)
		api.GET("/logs/:id", h.Get)
		api.GET("/stats/modules", h.ModuleStats)
	}
}

func tenantID(c *gin.Context) string {
	return c.GetHeader("X-Tenant-Id")
}

// 记录审计日志
func (h *Handler) Record(c *gin.Context) {
	var log model.AuditLog
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.TenantID = tenantID(c)
	if err := h.svc.Record(c.Request.Context(), &log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, log)
}

// 查询审计日志
func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	q := repo.AuditQuery{
		TenantID: tenantID(c),
		Module:   c.Query("module"),
		Action:   c.Query("action"),
		UserID:   c.Query("user_id"),
		StartAt:  c.Query("start_at"),
		EndAt:    c.Query("end_at"),
	}
	if s := c.Query("status"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			q.Status = int32(v)
		}
	}
	res, err := h.svc.Query(c.Request.Context(), q, repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// 获取单条
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	log, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, log)
}

// 按模块统计
func (h *Handler) ModuleStats(c *gin.Context) {
	res, err := h.svc.ModuleStats(c.Request.Context(), tenantID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": res})
}
