package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/erp-finance-service/internal/model"
	"github.com/oa-portal/erp-finance-service/internal/repo"
	"github.com/oa-portal/erp-finance-service/internal/service"
)

type Handler struct {
	svc *service.FinanceService
}

func NewHandler(svc *service.FinanceService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/finance")
	{
		// 应收
		api.POST("/receivables", h.CreateReceivable)
		api.GET("/receivables", h.ListReceivables)
		// 应付
		api.POST("/payables", h.CreatePayable)
		api.GET("/payables", h.ListPayables)
		// 收付款
		api.POST("/payments/receive", h.ReceivePayment)
		api.POST("/payments/make", h.MakePayment)
		api.GET("/payments", h.ListPayments)
	}
}

func tenantID(c *gin.Context) string {
	return c.GetHeader("X-Tenant-Id")
}

// 应收
func (h *Handler) CreateReceivable(c *gin.Context) {
	var r model.Receivable
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r.TenantID = tenantID(c)
	if err := h.svc.CreateReceivable(c.Request.Context(), &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, r)
}

func (h *Handler) ListReceivables(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.Query("status"))
	res, err := h.svc.ListReceivables(c.Request.Context(), tenantID(c), int32(status), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// 应付
func (h *Handler) CreatePayable(c *gin.Context) {
	var p model.Payable
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.TenantID = tenantID(c)
	if err := h.svc.CreatePayable(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) ListPayables(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.Query("status"))
	res, err := h.svc.ListPayables(c.Request.Context(), tenantID(c), int32(status), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// 收款
func (h *Handler) ReceivePayment(c *gin.Context) {
	var p model.FinancePayment
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.TenantID = tenantID(c)
	if err := h.svc.ReceivePayment(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

// 付款
func (h *Handler) MakePayment(c *gin.Context) {
	var p model.FinancePayment
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.TenantID = tenantID(c)
	if err := h.svc.MakePayment(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

// 收付款列表
func (h *Handler) ListPayments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	typ, _ := strconv.Atoi(c.Query("type"))
	res, err := h.svc.ListPayments(c.Request.Context(), tenantID(c), int32(typ), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
