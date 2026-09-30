package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/mdm-service/internal/model"
	"github.com/oa-portal/mdm-service/internal/service"
)

// 统一响应
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

// ---------- 客户 Handler ----------
type CustomerHandler struct {
	svc *service.CustomerService
}

func NewCustomerHandler(svc *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{svc: svc}
}

func (h *CustomerHandler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/customers")
	{
		g.GET("", h.List)
		g.GET("/:id", h.Get)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	}
}

func (h *CustomerHandler) List(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	keyword := c.Query("keyword")
	customerType := c.Query("customer_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.svc.List(c.Request.Context(), tenantID, keyword, customerType, page, pageSize)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, result)
}

func (h *CustomerHandler) Get(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	customer, err := h.svc.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		fail(c, 404, "not found")
		return
	}
	ok(c, customer)
}

func (h *CustomerHandler) Create(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	var req model.Customer
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	req.TenantID = tenantID
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, req)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req model.Customer
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	req.ID = id
	req.TenantID = tenantID
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, req)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), id, tenantID); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

// ---------- 商品 Handler ----------
type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/products")
	{
		g.GET("", h.List)
		g.GET("/:id", h.Get)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	}
}

func (h *ProductHandler) List(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	keyword := c.Query("keyword")
	category := c.Query("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	result, err := h.svc.List(c.Request.Context(), tenantID, keyword, category, page, pageSize)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, result)
}

func (h *ProductHandler) Get(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	p, err := h.svc.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		fail(c, 404, "not found")
		return
	}
	ok(c, p)
}

func (h *ProductHandler) Create(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	var req model.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	req.TenantID = tenantID
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, req)
}

func (h *ProductHandler) Update(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req model.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	req.ID = id
	req.TenantID = tenantID
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, req)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), id, tenantID); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

// ---------- 供应商 Handler ----------
type SupplierHandler struct {
	svc *service.SupplierService
}

func NewSupplierHandler(svc *service.SupplierService) *SupplierHandler {
	return &SupplierHandler{svc: svc}
}

func (h *SupplierHandler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/suppliers")
	{
		g.GET("", h.List)
		g.GET("/:id", h.Get)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	}
}

func (h *SupplierHandler) List(c *gin.Context) {
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

func (h *SupplierHandler) Get(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	s, err := h.svc.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		fail(c, 404, "not found")
		return
	}
	ok(c, s)
}

func (h *SupplierHandler) Create(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	var req model.Supplier
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	req.TenantID = tenantID
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, req)
}

func (h *SupplierHandler) Update(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req model.Supplier
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	req.ID = id
	req.TenantID = tenantID
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, req)
}

func (h *SupplierHandler) Delete(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), id, tenantID); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

// ---------- 组织 Handler ----------
type OrganizationHandler struct {
	svc *service.OrganizationService
}

func NewOrganizationHandler(svc *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{svc: svc}
}

func (h *OrganizationHandler) RegisterRoutes(r *gin.RouterGroup) {
	g := r.Group("/organizations")
	{
		g.GET("", h.List)
		g.GET("/:id", h.Get)
		g.POST("", h.Create)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Delete)
	}
}

func (h *OrganizationHandler) List(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	parentID, _ := strconv.ParseInt(c.DefaultQuery("parent_id", "0"), 10, 64)
	list, err := h.svc.List(c.Request.Context(), tenantID, parentID)
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, list)
}

func (h *OrganizationHandler) Get(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	o, err := h.svc.Get(c.Request.Context(), id, tenantID)
	if err != nil {
		fail(c, 404, "not found")
		return
	}
	ok(c, o)
}

func (h *OrganizationHandler) Create(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	var req model.Organization
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	req.TenantID = tenantID
	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, req)
}

func (h *OrganizationHandler) Update(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req model.Organization
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, err.Error())
		return
	}
	req.ID = id
	req.TenantID = tenantID
	if err := h.svc.Update(c.Request.Context(), &req); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, req)
}

func (h *OrganizationHandler) Delete(c *gin.Context) {
	tenantID := c.GetHeader("X-Tenant-Id")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), id, tenantID); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}
