package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/crm-customer-service/internal/model"
	"github.com/oa-portal/crm-customer-service/internal/repo"
	"github.com/oa-portal/crm-customer-service/internal/service"
)

type Handler struct {
	svc *service.CrmService
}

func NewHandler(svc *service.CrmService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/crm")
	{
		// 客户
		api.POST("/customers", h.CreateCustomer)
		api.GET("/customers", h.ListCustomers)
		api.GET("/customers/:id", h.GetCustomer)
		api.PUT("/customers/:id", h.UpdateCustomer)
		api.DELETE("/customers/:id", h.DeleteCustomer)
		// 联系人
		api.POST("/customers/:id/contacts", h.AddContact)
		api.GET("/customers/:id/contacts", h.ListContacts)
		api.DELETE("/contacts/:id", h.DeleteContact)
		// 跟进
		api.POST("/customers/:id/follow-ups", h.AddFollowUp)
		api.GET("/customers/:id/follow-ups", h.ListFollowUps)
	}
}

func tenantID(c *gin.Context) string {
	return c.GetHeader("X-Tenant-Id")
}

// 客户
func (h *Handler) CreateCustomer(c *gin.Context) {
	var cu model.Customer
	if err := c.ShouldBindJSON(&cu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cu.TenantID = tenantID(c)
	if err := h.svc.CreateCustomer(c.Request.Context(), &cu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cu)
}

func (h *Handler) ListCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	level := c.Query("level")
	status, _ := strconv.Atoi(c.Query("status"))
	res, err := h.svc.ListCustomers(c.Request.Context(), tenantID(c), level, int32(status), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetCustomer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cu, err := h.svc.GetCustomer(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cu)
}

func (h *Handler) UpdateCustomer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var cu model.Customer
	if err := c.ShouldBindJSON(&cu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cu.ID = id
	if err := h.svc.UpdateCustomer(c.Request.Context(), &cu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeleteCustomer(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.DeleteCustomer(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 联系人
func (h *Handler) AddContact(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var ct model.Contact
	if err := c.ShouldBindJSON(&ct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ct.CustomerID = id
	if err := h.svc.AddContact(c.Request.Context(), &ct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ct)
}

func (h *Handler) ListContacts(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	list, err := h.svc.ListContacts(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) DeleteContact(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.DeleteContact(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// 跟进
func (h *Handler) AddFollowUp(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var f model.FollowUp
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f.CustomerID = id
	if err := h.svc.AddFollowUp(c.Request.Context(), &f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *Handler) ListFollowUps(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	list, err := h.svc.ListFollowUps(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
