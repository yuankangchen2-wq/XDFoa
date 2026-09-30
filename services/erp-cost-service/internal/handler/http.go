package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/erp-cost-service/internal/model"
	"github.com/oa-portal/erp-cost-service/internal/repo"
	"github.com/oa-portal/erp-cost-service/internal/service"
)

type Handler struct {
	svc *service.CostService
}

func NewHandler(svc *service.CostService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/cost")
	{
		// 成本中心
		api.POST("/centers", h.CreateCostCenter)
		api.GET("/centers", h.ListCostCenters)
		// 产品成本
		api.POST("/products", h.CalculateProductCost)
		api.GET("/products", h.ListProductCosts)
		api.POST("/collect", h.CollectCost)
		// 成本记录
		api.POST("/records", h.AddCostRecord)
		api.GET("/records", h.ListCostRecords)
	}
}

func tenantID(c *gin.Context) string {
	return c.GetHeader("X-Tenant-Id")
}

// 成本中心
func (h *Handler) CreateCostCenter(c *gin.Context) {
	var cc model.CostCenter
	if err := c.ShouldBindJSON(&cc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cc.TenantID = tenantID(c)
	if err := h.svc.CreateCostCenter(c.Request.Context(), &cc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cc)
}

func (h *Handler) ListCostCenters(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	typ, _ := strconv.Atoi(c.Query("type"))
	res, err := h.svc.ListCostCenters(c.Request.Context(), tenantID(c), int32(typ), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// 产品成本
func (h *Handler) CalculateProductCost(c *gin.Context) {
	var pc model.ProductCost
	if err := c.ShouldBindJSON(&pc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pc.TenantID = tenantID(c)
	if err := h.svc.CalculateProductCost(c.Request.Context(), &pc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pc)
}

func (h *Handler) ListProductCosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	productID, _ := strconv.Atoi(c.Query("product_id"))
	res, err := h.svc.ListProductCosts(c.Request.Context(), tenantID(c), int64(productID), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

type collectReq struct {
	ProductID int64   `json:"product_id" binding:"required"`
	Period    string  `json:"period" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
}

func (h *Handler) CollectCost(c *gin.Context) {
	var req collectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pc, err := h.svc.CollectCost(c.Request.Context(), tenantID(c), req.ProductID, req.Period, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pc)
}

// 成本记录
func (h *Handler) AddCostRecord(c *gin.Context) {
	var r model.CostRecord
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r.TenantID = tenantID(c)
	if err := h.svc.AddCostRecord(c.Request.Context(), &r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (h *Handler) ListCostRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.Query("status"))
	res, err := h.svc.ListCostRecords(c.Request.Context(), tenantID(c), int32(status), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
