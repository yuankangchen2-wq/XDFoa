package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/crm-analytics-service/internal/model"
	"github.com/oa-portal/crm-analytics-service/internal/repo"
	"github.com/oa-portal/crm-analytics-service/internal/service"
)

type Handler struct {
	svc *service.AnalyticsService
}

func NewHandler(svc *service.AnalyticsService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/analytics")
	{
		api.GET("/dashboard", h.GetDashboard)
		// 销售快照
		api.POST("/sales", h.RecordSales)
		api.GET("/sales", h.ListSales)
		api.GET("/sales/summary", h.SalesSummary)
		// 漏斗
		api.POST("/funnel", h.RecordFunnel)
		api.GET("/funnel", h.ListFunnel)
		// 客户统计
		api.POST("/customers", h.RecordCustomerStat)
		api.GET("/customers", h.ListCustomerStats)
	}
}

func tenantID(c *gin.Context) string {
	return c.GetHeader("X-Tenant-Id")
}

// 看板
func (h *Handler) GetDashboard(c *gin.Context) {
	d, err := h.svc.GetDashboard(c.Request.Context(), tenantID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}

// 销售
func (h *Handler) RecordSales(c *gin.Context) {
	var o model.SalesSnapshot
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o.TenantID = tenantID(c)
	if err := h.svc.RecordSalesSnapshot(c.Request.Context(), &o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) ListSales(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	period := c.Query("period")
	res, err := h.svc.ListSalesSnapshots(c.Request.Context(), tenantID(c), period, repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) SalesSummary(c *gin.Context) {
	res, err := h.svc.SalesSummaryByPeriod(c.Request.Context(), tenantID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": res})
}

// 漏斗
func (h *Handler) RecordFunnel(c *gin.Context) {
	var o model.FunnelStage
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o.TenantID = tenantID(c)
	if err := h.svc.RecordFunnelStage(c.Request.Context(), &o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) ListFunnel(c *gin.Context) {
	period := c.Query("period")
	res, err := h.svc.ListFunnel(c.Request.Context(), tenantID(c), period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": res})
}

// 客户统计
func (h *Handler) RecordCustomerStat(c *gin.Context) {
	var o model.CustomerStat
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o.TenantID = tenantID(c)
	if err := h.svc.RecordCustomerStat(c.Request.Context(), &o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) ListCustomerStats(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	res, err := h.svc.ListCustomerStats(c.Request.Context(), tenantID(c), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
