package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/crm-opportunity-service/internal/model"
	"github.com/oa-portal/crm-opportunity-service/internal/repo"
	"github.com/oa-portal/crm-opportunity-service/internal/service"
)

type Handler struct {
	svc *service.OpportunityService
}

func NewHandler(svc *service.OpportunityService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/opp")
	{
		api.POST("/opportunities", h.CreateOpportunity)
		api.GET("/opportunities", h.ListOpportunities)
		api.GET("/opportunities/:id", h.GetOpportunity)
		api.PUT("/opportunities/:id", h.UpdateOpportunity)
		api.DELETE("/opportunities/:id", h.DeleteOpportunity)
		api.POST("/opportunities/:id/transition", h.TransitionStage)
		api.POST("/opportunities/:id/convert-order", h.ConvertToOrder)
	}
}

func tenantID(c *gin.Context) string {
	return c.GetHeader("X-Tenant-Id")
}

func (h *Handler) CreateOpportunity(c *gin.Context) {
	var o model.Opportunity
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o.TenantID = tenantID(c)
	if err := h.svc.CreateOpportunity(c.Request.Context(), &o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) ListOpportunities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	ownerID := c.Query("owner_id")
	stage, _ := strconv.Atoi(c.Query("stage"))
	res, err := h.svc.ListOpportunities(c.Request.Context(), tenantID(c), ownerID, int32(stage), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetOpportunity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	o, err := h.svc.GetOpportunity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *Handler) UpdateOpportunity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var o model.Opportunity
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o.ID = id
	if err := h.svc.UpdateOpportunity(c.Request.Context(), &o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeleteOpportunity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.svc.DeleteOpportunity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

type transitionReq struct {
	ToStage int32  `json:"to_stage" binding:"required"`
	Remark  string `json:"remark"`
}

func (h *Handler) TransitionStage(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req transitionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.TransitionStage(c.Request.Context(), int32(id), req.ToStage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

type convertReq struct {
	UserID      int64 `json:"user_id" binding:"required"`
	WarehouseID int64 `json:"warehouse_id" binding:"required"`
}

func (h *Handler) ConvertToOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var req convertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ConvertToOrder(c.Request.Context(), id, req.UserID, req.WarehouseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
