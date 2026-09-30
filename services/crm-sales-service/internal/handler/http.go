package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/crm-sales-service/internal/model"
	"github.com/oa-portal/crm-sales-service/internal/repo"
	"github.com/oa-portal/crm-sales-service/internal/service"
)

type Handler struct {
	svc *service.SalesService
}

func NewHandler(svc *service.SalesService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/sales")
	{
		// 销售订单
		api.POST("/orders", h.CreateSalesOrder)
		api.GET("/orders", h.ListSalesOrders)
		api.GET("/orders/:id", h.GetSalesOrder)
		api.POST("/orders/:id/approve", h.ApproveSalesOrder)
		api.POST("/orders/:id/ship", h.ShipSalesOrder)
		api.POST("/orders/:id/complete", h.CompleteSalesOrder)
		// 合同
		api.POST("/contracts", h.CreateContract)
		api.GET("/contracts", h.ListContracts)
		api.GET("/contracts/:id", h.GetContract)
		api.POST("/contracts/:id/approve", h.ApproveContract)
		api.POST("/contracts/:id/archive", h.ArchiveContract)
		// 回款
		api.POST("/contracts/:id/payments", h.AddPayment)
		api.GET("/contracts/:id/payments", h.ListPayments)
	}
}

func tenantID(c *gin.Context) string {
	return c.GetHeader("X-Tenant-Id")
}

// 销售订单
func (h *Handler) CreateSalesOrder(c *gin.Context) {
	var o model.SalesOrder
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o.TenantID = tenantID(c)
	if err := h.svc.CreateSalesOrder(c.Request.Context(), &o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) ListSalesOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.Query("status"))
	res, err := h.svc.ListSalesOrders(c.Request.Context(), tenantID(c), int32(status), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetSalesOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	o, err := h.svc.GetSalesOrder(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) ApproveSalesOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.ApproveSalesOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ShipSalesOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.ShipSalesOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) CompleteSalesOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.CompleteSalesOrder(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 合同
func (h *Handler) CreateContract(c *gin.Context) {
	var ct model.Contract
	if err := c.ShouldBindJSON(&ct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ct.TenantID = tenantID(c)
	if err := h.svc.CreateContract(c.Request.Context(), &ct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ct)
}

func (h *Handler) ListContracts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.Query("status"))
	res, err := h.svc.ListContracts(c.Request.Context(), tenantID(c), int32(status), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetContract(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ct, err := h.svc.GetContract(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ct)
}

func (h *Handler) ApproveContract(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.ApproveContract(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ArchiveContract(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.ArchiveContract(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 回款
func (h *Handler) AddPayment(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var p model.Payment
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ContractID = id
	p.TenantID = tenantID(c)
	if err := h.svc.AddPayment(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) ListPayments(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	list, err := h.svc.ListPayments(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
