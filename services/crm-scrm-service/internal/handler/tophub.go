package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ---------------- 全局搜索 ----------------

func (h *Handler) GlobalSearch(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(400, gin.H{"error": "请输入搜索关键词"})
		return
	}
	res, err := h.svc.SearchAll(c.Request.Context(), q)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

// ---------------- 线索详情 / 跟进日志 ----------------

func (h *Handler) GetLead(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	l, err := h.svc.GetLead(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "线索不存在"})
		return
	}
	c.JSON(200, l)
}

func (h *Handler) LeadFollowLogs(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	list, err := h.svc.ListFollowLogs(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (h *Handler) RecentCalls(c *gin.Context) {
	list, err := h.svc.RecentCalls(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

// ---------------- 云座席双呼 ----------------

func (h *Handler) CallDial(c *gin.Context) {
	var body struct {
		LeadID int64  `json:"lead_id"`
		Phone  string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	logEntry, err := h.svc.Dial(c.Request.Context(), body.LeadID, body.Phone)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, logEntry)
}

func (h *Handler) CallComplete(c *gin.Context) {
	var body struct {
		LogID       int64 `json:"log_id"`
		Connected   bool  `json:"connected"`
		DurationSec int32 `json:"duration_sec"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CompleteCall(c.Request.Context(), body.LogID, body.Connected, body.DurationSec); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

// ---------------- 站内通知 ----------------

func (h *Handler) ListNotifications(c *gin.Context) {
	res, err := h.svc.ListNotifications(c.Request.Context(), c.Query("unread"), parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) UnreadNotifications(c *gin.Context) {
	n, err := h.svc.UnreadNotifications(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"unread": n})
}

func (h *Handler) MarkNotificationsRead(c *gin.Context) {
	var body struct {
		IDs []int64 `json:"ids"`
		All bool    `json:"all"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.MarkNotificationsRead(c.Request.Context(), body.IDs, body.All); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}
