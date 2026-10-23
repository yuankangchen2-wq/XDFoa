package handler

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/crm-scrm-service/internal/model"
	"github.com/oa-portal/crm-scrm-service/internal/repo"
	"github.com/oa-portal/crm-scrm-service/internal/service"
)

// RegisterP2 二三期路由（在 Register 中调用）
func (h *Handler) RegisterP2(api *gin.RouterGroup) {
	api.POST("/leads/recycle", h.RecycleLeads)
	api.GET("/lead-task-snapshot", h.LeadTaskSnapshot)
	api.GET("/employee-tasks", h.ListEmployeeTasks)

	api.GET("/opportunities", h.ListOpportunities)
	api.POST("/opportunities", h.CreateOpportunity)
	api.PUT("/opportunities/:id", h.UpdateOpportunity)
	api.GET("/opportunities/score-rules", h.ScoreRules)
	api.POST("/opportunities/assign", h.AssignOpportunities)
	api.POST("/opportunities/recycle", h.RecycleOpportunities)
	api.POST("/opportunities/batch-archive", h.ArchiveOpportunities)
	api.POST("/opportunities/activate", h.ActivateOpportunities)
	api.POST("/opportunities/push-91", h.Push91)
	api.POST("/opportunities/push-hedgehog", h.PushHedgehog)
	api.POST("/opportunities/apply-deal", h.ApplyDeal)
	api.GET("/assign-records", h.ListAssignRecords)
	api.GET("/inventory", h.Inventory)

	api.GET("/appeals", h.ListAppeals)
	api.POST("/appeals", h.CreateAppeal)
	api.POST("/appeals/handle", h.HandleAppeals)
	api.GET("/approvals", h.ListApprovals)
	api.POST("/approvals/handle", h.HandleApprovals)
	api.GET("/deal-logs", h.ListDealLogs)
	api.GET("/tickets", h.ListTickets)
	api.GET("/mall-orders", h.ListOrders)

	api.GET("/grid-config", h.GetGridConfig)
	api.POST("/grid-config", h.SaveGridConfig)
}

// ---------------- 线索工作台 ----------------

func (h *Handler) RecycleLeads(c *gin.Context) {
	var body struct {
		IDs    []int64 `json:"ids"`
		Reason string  `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RecycleLeads(c.Request.Context(), body.IDs, body.Reason); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) LeadTaskSnapshot(c *gin.Context) {
	res, err := h.svc.LeadTaskSnapshot(c.Request.Context(), c.DefaultQuery("scope", "mine"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"items": res})
}

func (h *Handler) ListEmployeeTasks(c *gin.Context) {
	q := c.Query
	f := repo.EmployeeTaskFilter{
		BizKind: q("biz_kind"), OwnerName: q("owner_name"), TaskName: q("task_name"),
		TaskType: q("task_type"), Status: q("status"),
		Start: parseDay(q("start")), End: parseDayEnd(q("end")),
	}
	res, err := h.svc.ListEmployeeTasks(c.Request.Context(), f, parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

// ---------------- 商机 ----------------

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func parseOppFilter(c *gin.Context) repo.OppFilter {
	q := c.Query
	f := repo.OppFilter{
		Keyword: q("keyword"), Name: q("name"), Phone: q("phone"), Code: q("code"),
		StudentNo: q("student_no"), Status: q("status"), Level: q("level"), Grade2: q("grade2"),
		EnterType: q("enter_type"), OppType: q("opp_type"), Category: q("category"),
		NewOld: q("new_old"), Grade: q("grade"), District: q("district"), School: q("school"),
		IntentionDept: q("intention_dept"), IntentionProject: q("intention_project"),
		IntentionCampus: q("intention_campus"), ContactPoint: q("contact_point"),
		ChannelL1: q("channel_l1"), ChannelL2: q("channel_l2"), ChannelL3: q("channel_l3"),
		ActivityName: q("activity_name"), OwnerName: q("owner_name"), Creator: q("creator"),
		Pushed91: q("pushed_91"), ResourceNo91: q("resource_no_91"), IsIntl: q("is_intl"),
		AddWechat: q("add_wechat"), Scope: q("scope"), Quick: q("quick"), OwnerID: q("owner_id"),
		ScoreMin: atoiDefault(q("score_min"), 0), ScoreMax: atoiDefault(q("score_max"), 0),
		InitMin: atoiDefault(q("init_min"), 0), InitMax: atoiDefault(q("init_max"), 0),
		MaxMin: atoiDefault(q("max_min"), 0), MaxMax: atoiDefault(q("max_max"), 0),
		CreatedStart: parseDay(q("created_start")), CreatedEnd: parseDayEnd(q("created_end")),
		FollowStart: parseDay(q("follow_start")), FollowEnd: parseDayEnd(q("follow_end")),
		AssignStart: parseDay(q("assign_start")), AssignEnd: parseDayEnd(q("assign_end")),
		VisitStart: parseDay(q("visit_start")), VisitEnd: parseDayEnd(q("visit_end")),
		NextStart: parseDay(q("next_start")), NextEnd: parseDayEnd(q("next_end")),
		ArchiveStart: parseDay(q("archive_start")), ArchiveEnd: parseDayEnd(q("archive_end")),
	}
	if s := q("statuses"); s != "" {
		f.Statuses = strings.Split(s, ",")
	}
	return f
}

func (h *Handler) ListOpportunities(c *gin.Context) {
	scope := c.DefaultQuery("scope", "mine")
	res, err := h.svc.ListOpportunities(c.Request.Context(), scope, parseOppFilter(c), parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) CreateOpportunity(c *gin.Context) {
	var o model.Opportunity
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateOpportunity(c.Request.Context(), &o); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, o)
}

func (h *Handler) UpdateOpportunity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var fields map[string]interface{}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateOpportunityFields(c.Request.Context(), id, fields); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) ScoreRules(c *gin.Context) {
	c.JSON(200, service.ScoreRuleText())
}

type oppIDsBody struct {
	IDs []int64 `json:"ids"`
}

func (h *Handler) AssignOpportunities(c *gin.Context) {
	var body struct {
		IDs          []int64 `json:"ids"`
		ReceiverID   string  `json:"receiver_id"`
		ReceiverName string  `json:"receiver_name"`
		ReceiverDept string  `json:"receiver_dept"`
		Auto         int32   `json:"auto"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AssignOpportunities(c.Request.Context(), body.IDs, body.ReceiverID, body.ReceiverName, body.ReceiverDept, body.Auto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) RecycleOpportunities(c *gin.Context) {
	var body struct {
		IDs    []int64 `json:"ids"`
		Reason string  `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RecycleOpportunities(c.Request.Context(), body.IDs, body.Reason); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) ArchiveOpportunities(c *gin.Context) {
	var body struct {
		IDs  []int64 `json:"ids"`
		Rule string  `json:"rule"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ArchiveOpportunities(c.Request.Context(), body.IDs, body.Rule); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) ActivateOpportunities(c *gin.Context) {
	var body oppIDsBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ActivateOpportunities(c.Request.Context(), body.IDs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) Push91(c *gin.Context) {
	var body oppIDsBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	n, err := h.svc.Push91(c.Request.Context(), body.IDs)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "pushed": n})
}

func (h *Handler) PushHedgehog(c *gin.Context) {
	var body oppIDsBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	n, err := h.svc.PushHedgehog(c.Request.Context(), body.IDs)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "pushed": n})
}

func (h *Handler) ApplyDeal(c *gin.Context) {
	var a model.DealApproval
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ApplyDeal(c.Request.Context(), &a); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, a)
}

func (h *Handler) ListAssignRecords(c *gin.Context) {
	q := c.Query
	f := repo.AssignFilter{
		Keyword: q("keyword"), Receiver: q("receiver"), Operator: q("operator"),
		Auto: q("auto"), Status: q("status"),
		Start: parseDay(q("start")), End: parseDayEnd(q("end")),
	}
	res, err := h.svc.ListAssignRecords(c.Request.Context(), c.DefaultQuery("scope", "mine"), f, parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) Inventory(c *gin.Context) {
	res, err := h.svc.Inventory(c.Request.Context(), c.DefaultQuery("dim", "school"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

// ---------------- 申诉 / 审批 / 成单 / 工单 / 订单 ----------------

func (h *Handler) ListAppeals(c *gin.Context) {
	res, err := h.svc.ListAppeals(c.Request.Context(), c.Query("biz_line"), c.Query("keyword"), c.Query("status"), parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) CreateAppeal(c *gin.Context) {
	var a model.Appeal
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateAppeal(c.Request.Context(), &a); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, a)
}

func (h *Handler) HandleAppeals(c *gin.Context) {
	var body struct {
		IDs  []int64 `json:"ids"`
		Pass bool    `json:"pass"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.HandleAppeals(c.Request.Context(), body.IDs, body.Pass); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) ListApprovals(c *gin.Context) {
	res, err := h.svc.ListApprovals(c.Request.Context(), c.Query("keyword"), c.Query("status"),
		c.Query("applicant"), c.Query("biz_type"), parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) HandleApprovals(c *gin.Context) {
	var body struct {
		IDs  []int64 `json:"ids"`
		Pass bool    `json:"pass"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.HandleApprovals(c.Request.Context(), body.IDs, body.Pass); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) ListDealLogs(c *gin.Context) {
	q := c.Query
	f := repo.DealFilter{BizLine: q("biz_line"), Keyword: q("keyword"), OwnerName: q("owner_name"),
		Start: parseDay(q("start")), End: parseDayEnd(q("end"))}
	res, err := h.svc.ListDeals(c.Request.Context(), f, parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) ListTickets(c *gin.Context) {
	res, err := h.svc.ListTickets(c.Request.Context(), c.DefaultQuery("scope", "all"),
		c.Query("status"), c.Query("biz_type"), c.Query("keyword"), parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) ListOrders(c *gin.Context) {
	q := c.Query
	f := repo.OrderFilter{
		Keyword: q("keyword"), ProductType: q("product_type"), Status: q("status"),
		BizType: q("biz_type"), SaleType: q("sale_type"),
		Start: parseDay(q("start")), End: parseDayEnd(q("end")),
	}
	res, err := h.svc.ListOrders(c.Request.Context(), f, parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

// ---------------- 列/筛选配置 ----------------

func (h *Handler) GetGridConfig(c *gin.Context) {
	cfg, err := h.svc.GetGridConfig(c.Request.Context(), c.Query("page_key"), c.DefaultQuery("cfg_type", "column"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"config": cfg})
}

func (h *Handler) SaveGridConfig(c *gin.Context) {
	var body struct {
		PageKey string `json:"page_key"`
		CfgType string `json:"cfg_type"`
		Config  string `json:"config"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if body.CfgType == "" {
		body.CfgType = "column"
	}
	if err := h.svc.SaveGridConfig(c.Request.Context(), body.PageKey, body.CfgType, body.Config); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}
