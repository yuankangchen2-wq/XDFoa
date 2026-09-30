package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/production-service/internal/model"
	"github.com/oa-portal/production-service/internal/repo"
	"github.com/oa-portal/production-service/internal/service"
)

type Handler struct {
	svc *service.ProductionService
}

func NewHandler(svc *service.ProductionService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/production")
	{
		// BOM
		api.POST("/boms", h.CreateBom)
		api.GET("/boms", h.ListBoms)
		api.POST("/boms/:id/enable", h.EnableBom)
		// 工单
		api.POST("/work-orders", h.CreateWorkOrder)
		api.GET("/work-orders", h.ListWorkOrders)
		api.GET("/work-orders/:id", h.GetWorkOrder)
		// 领料
		api.POST("/work-orders/:id/issue", h.IssueMaterial)
		api.GET("/material-issues", h.ListMaterialIssues)
		// 报工
		api.POST("/work-orders/:id/report", h.ReportProduction)
	}
}

func tenantID(c *gin.Context) string {
	return c.GetHeader("X-Tenant-Id")
}

// BOM
func (h *Handler) CreateBom(c *gin.Context) {
	var b model.Bom
	if err := c.ShouldBindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b.TenantID = tenantID(c)
	if err := h.svc.CreateBom(c.Request.Context(), &b); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *Handler) ListBoms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	res, err := h.svc.ListBoms(c.Request.Context(), tenantID(c), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) EnableBom(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.EnableBom(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 工单
func (h *Handler) CreateWorkOrder(c *gin.Context) {
	var w model.WorkOrder
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	w.TenantID = tenantID(c)
	if err := h.svc.CreateWorkOrder(c.Request.Context(), &w); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *Handler) ListWorkOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.Query("status"))
	res, err := h.svc.ListWorkOrders(c.Request.Context(), tenantID(c), int32(status), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetWorkOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	w, err := h.svc.GetWorkOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, w)
}

// 领料
func (h *Handler) IssueMaterial(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	issue, err := h.svc.IssueMaterial(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, issue)
}

func (h *Handler) ListMaterialIssues(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	res, err := h.svc.ListMaterialIssues(c.Request.Context(), tenantID(c), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// 报工
type reportReq struct {
	ProducedQty int64 `json:"produced_qty" binding:"required"`
}

func (h *Handler) ReportProduction(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req reportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ReportProduction(c.Request.Context(), id, req.ProducedQty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
