package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/inventory-service/internal/model"
	"github.com/oa-portal/inventory-service/internal/repo"
	"github.com/oa-portal/inventory-service/internal/service"
)

type Handler struct {
	svc *service.InventoryService
}

func NewHandler(svc *service.InventoryService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/inventory")
	{
		// 仓库
		api.POST("/warehouses", h.CreateWarehouse)
		api.GET("/warehouses", h.ListWarehouses)
		// 库存查询
		api.GET("/stocks", h.ListStocks)
		api.GET("/stocks/:sku/:warehouse", h.GetStock)
		// 库存操作
		api.POST("/stocks/allocate", h.AllocateStock)
		api.POST("/stocks/release", h.ReleaseStock)
		api.POST("/stocks/deduct", h.DeductStock)
		api.POST("/stocks/in", h.StockIn)
		// 流水
		api.GET("/journals", h.ListJournals)
		// 调拨
		api.POST("/transfers", h.CreateTransfer)
		api.POST("/transfers/:id/out", h.ConfirmTransferOut)
		api.POST("/transfers/:id/in", h.ConfirmTransferIn)
		api.GET("/transfers", h.ListTransfers)
		// 盘点
		api.POST("/stocktakes", h.CreateStocktake)
		api.GET("/stocktakes", h.ListStocktakes)
	}
}

func tenantID(c *gin.Context) string {
	return c.GetHeader("X-Tenant-Id")
}

type allocateReq struct {
	SkuCode     string `json:"sku_code" binding:"required"`
	WarehouseID int64  `json:"warehouse_id" binding:"required"`
	Quantity    int64  `json:"quantity" binding:"required"`
	ReferenceID string `json:"reference_id"`
}

func (h *Handler) AllocateStock(c *gin.Context) {
	var req allocateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AllocateStock(c.Request.Context(), tenantID(c), req.SkuCode, req.WarehouseID, req.Quantity, req.ReferenceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ReleaseStock(c *gin.Context) {
	var req allocateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ReleaseStock(c.Request.Context(), tenantID(c), req.SkuCode, req.WarehouseID, req.Quantity, req.ReferenceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) DeductStock(c *gin.Context) {
	var req allocateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.DeductStock(c.Request.Context(), tenantID(c), req.SkuCode, req.WarehouseID, req.Quantity, req.ReferenceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

type stockInReq struct {
	SkuCode     string `json:"sku_code" binding:"required"`
	WarehouseID int64  `json:"warehouse_id" binding:"required"`
	Quantity    int64  `json:"quantity" binding:"required"`
	ReferenceID string `json:"reference_id"`
	BatchNo     string `json:"batch_no"`
	ExpireAt    int64  `json:"expire_at"`
}

func (h *Handler) StockIn(c *gin.Context) {
	var req stockInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.StockIn(c.Request.Context(), tenantID(c), req.SkuCode, req.WarehouseID, req.Quantity, req.ReferenceID, req.BatchNo, req.ExpireAt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ListStocks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	res, err := h.svc.ListStocks(c.Request.Context(), tenantID(c), c.Query("sku_code"),
		parseInt64(c.Query("warehouse_id")), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetStock(c *gin.Context) {
	sku := c.Param("sku")
	wh := parseInt64(c.Param("warehouse"))
	// 简化：返回该 sku+warehouse 的库存（从 list 取第一条）
	res, err := h.svc.ListStocks(c.Request.Context(), tenantID(c), sku, wh, repo.Page{Page: 1, PageSize: 1})
	if err != nil || len(res.Items) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, res.Items[0])
}

func (h *Handler) ListJournals(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	res, err := h.svc.ListJournals(c.Request.Context(), tenantID(c), c.Query("sku_code"),
		parseInt64(c.Query("warehouse_id")), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// 仓库
func (h *Handler) CreateWarehouse(c *gin.Context) {
	var w model.Warehouse
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	w.TenantID = tenantID(c)
	if err := h.svc.CreateWarehouse(c.Request.Context(), &w); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *Handler) ListWarehouses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	res, err := h.svc.ListWarehouses(c.Request.Context(), tenantID(c), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// 调拨
type transferReq struct {
	FromWarehouseID int64  `json:"from_warehouse_id" binding:"required"`
	ToWarehouseID   int64  `json:"to_warehouse_id" binding:"required"`
	SkuCode         string `json:"sku_code" binding:"required"`
	Quantity        int64  `json:"quantity" binding:"required"`
}

func (h *Handler) CreateTransfer(c *gin.Context) {
	var req transferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := h.svc.CreateTransfer(c.Request.Context(), tenantID(c), req.FromWarehouseID, req.ToWarehouseID, req.SkuCode, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *Handler) ConfirmTransferOut(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	if err := h.svc.ConfirmTransferOut(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ConfirmTransferIn(c *gin.Context) {
	id := parseInt64(c.Param("id"))
	if err := h.svc.ConfirmTransferIn(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) ListTransfers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	res, err := h.svc.ListTransfers(c.Request.Context(), tenantID(c), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// 盘点
type stocktakeReq struct {
	WarehouseID int64                 `json:"warehouse_id" binding:"required"`
	Items       []model.StocktakeItem `json:"items"`
}

func (h *Handler) CreateStocktake(c *gin.Context) {
	var req stocktakeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	st, err := h.svc.CreateStocktake(c.Request.Context(), tenantID(c), req.WarehouseID, req.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, st)
}

func (h *Handler) ListStocktakes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	res, err := h.svc.ListStocktakes(c.Request.Context(), tenantID(c), repo.Page{Page: page, PageSize: pageSize})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func parseInt64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}
