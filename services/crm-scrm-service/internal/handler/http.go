package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/crm-scrm-service/internal/model"
	"github.com/oa-portal/crm-scrm-service/internal/repo"
	"github.com/oa-portal/crm-scrm-service/internal/service"
	"github.com/oa-portal/crm-scrm-service/internal/xlsxutil"
)

type Handler struct {
	svc *service.ScrmService
}

func NewHandler(svc *service.ScrmService) *Handler { return &Handler{svc: svc} }

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := r.Group("/api/scrm")
	{
		api.GET("/me", h.Me)
		api.GET("/dict", h.Dict)
		api.GET("/employees", h.Employees)

		api.GET("/home/todos", h.Todos)
		api.GET("/home/ranking", h.Ranking)
		api.GET("/home/grade-dist", h.GradeDist)

		api.GET("/schedules", h.ListSchedules)
		api.POST("/schedules", h.CreateSchedule)

		api.GET("/activities", h.ListActivities)
		api.GET("/activities/all", h.AllActivities)
		api.POST("/activities", h.CreateActivity)
		api.PUT("/activities/:id", h.UpdateActivity)

		api.GET("/leads", h.ListLeads)
		api.GET("/leads/export", h.ExportLeads)
		api.GET("/leads/stat", h.LeadStat)
		api.POST("/leads/handover", h.HandoverLeads)
		api.POST("/leads/batch-adjust", h.BatchAdjust)
		api.POST("/leads/batch-archive", h.BatchArchive)
		api.POST("/leads/activate", h.ActivateLeads)

		api.GET("/exceptions", h.ListExceptions)
		api.POST("/exceptions/ignore", h.IgnoreExceptions)
		api.POST("/exceptions/adjust", h.AdjustExceptions)

		api.GET("/task-templates", h.ListTaskTemplates)
		api.POST("/task-templates", h.CreateTaskTemplate)
		api.POST("/imports", h.UploadImport)
		api.GET("/imports/template", h.DownloadTemplate)
		api.GET("/import-logs", h.ListImportLogs)
		api.GET("/import-logs/:id/fail-file", h.DownloadFailFile)
		api.GET("/import-logs/:id/dedup", h.DedupDetails)

		api.GET("/data-records", h.ListDataRecords)
		api.POST("/data-records/handover", h.HandoverDataRecords)
		api.POST("/data-records/adjust", h.AdjustDataRecords)
	}
}

// ---------------- 基础 ----------------

func (h *Handler) Me(c *gin.Context) {
	c.JSON(200, h.svc.Me())
}

func (h *Handler) Dict(c *gin.Context) {
	c.JSON(200, h.svc.Dicts(c.Query("kind")))
}

func (h *Handler) Employees(c *gin.Context) {
	list, err := h.svc.Employees(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

// ---------------- 首页 ----------------

func (h *Handler) Todos(c *gin.Context) {
	role := c.DefaultQuery("role", "employee")
	list, err := h.svc.Todos(c.Request.Context(), role)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (h *Handler) Ranking(c *gin.Context) {
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))
	list, err := h.svc.Ranking(c.Request.Context(), month, c.Query("campus"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (h *Handler) GradeDist(c *gin.Context) {
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))
	list, err := h.svc.GradeDistribution(c.Request.Context(), month, c.Query("campus"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

// ---------------- 日程 ----------------

func (h *Handler) ListSchedules(c *gin.Context) {
	list, err := h.svc.ListSchedules(c.Request.Context(), c.Query("month"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (h *Handler) CreateSchedule(c *gin.Context) {
	var body struct {
		Title string `json:"title"`
		Date  string `json:"date"`
		Time  string `json:"time"`
		Type  string `json:"type"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if body.Type == "" {
		body.Type = "日程"
	}
	sc, err := h.svc.CreateSchedule(c.Request.Context(), body.Title, body.Date, body.Time, body.Type)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, sc)
}

// ---------------- 活动 ----------------

func parsePage(c *gin.Context) repo.Page {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return repo.Page{Page: page, PageSize: size}
}

func parseDay(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		return nil
	}
	return &t
}
func parseDayEnd(s string) *time.Time {
	t := parseDay(s)
	if t == nil {
		return nil
	}
	e := t.Add(24 * time.Hour).Add(-time.Second)
	return &e
}

func (h *Handler) ListActivities(c *gin.Context) {
	f := repo.ActivityFilter{
		Name: c.Query("name"), Type: c.Query("type"), SourceSystem: c.Query("source_system"),
		Creator: c.Query("creator"), CreatorDept: c.Query("creator_dept"),
		Code: c.Query("code"), ExternalID: c.Query("external_id"), LeadStatus: c.Query("lead_status"),
		CreatedStart: parseDay(c.Query("created_start")), CreatedEnd: parseDayEnd(c.Query("created_end")),
	}
	res, err := h.svc.ListActivities(c.Request.Context(), f, parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) AllActivities(c *gin.Context) {
	list, err := h.svc.AllActivities(c.Request.Context(), c.Query("keyword"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (h *Handler) CreateActivity(c *gin.Context) {
	var a model.Activity
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateActivity(c.Request.Context(), &a); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, a)
}

func (h *Handler) UpdateActivity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var a model.Activity
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	a.ID = id
	if err := h.svc.UpdateActivity(c.Request.Context(), &a); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

// ---------------- 线索 ----------------

func parseLeadFilter(c *gin.Context) repo.LeadFilter {
	q := c.Query
	activityID, _ := strconv.ParseInt(q("activity_id"), 10, 64)
	followMin, _ := strconv.Atoi(q("follow_min"))
	followMax, _ := strconv.Atoi(q("follow_max"))
	f := repo.LeadFilter{
		Keyword: q("keyword"), Name: q("name"), Phone: q("phone"), StudentNo: q("student_no"),
		StudentID: q("student_id"), Status: q("status"), EnterType: q("enter_type"),
		LeadType: q("lead_type"), StudentTag: q("student_tag"), NewOldType: q("new_old"),
		School: q("school"), Grade: q("grade"), District: q("district"), CampusCode: q("campus_code"),
		IntentionLevel: q("intention_level"), IntentionGrade: q("intention_grade"),
		IntentionDept: q("intention_dept"), IntentionProject: q("intention_project"),
		IntentionCampus: q("intention_campus"), IntentionSubject: q("intention_subject"),
		ChannelL1: q("channel_l1"), ChannelL2: q("channel_l2"), ChannelL3: q("channel_l3"),
		BusinessChannel: q("business_channel"), SystemSource: q("system_source"),
		CreateMethod: q("create_method"), ActivityID: activityID, ActivityName: q("activity_name"),
		OwnerName: q("owner_name"), Creator: q("creator"), CreatorDept: q("creator_dept"),
		Recommender: q("recommender"), ArchiveRule: q("archive_rule"), ArchiveWay: q("archive_way"),
		Remark: q("remark"), FollowMin: followMin, FollowMax: followMax, View: q("view"),
		CreatedStart: parseDay(q("created_start")), CreatedEnd: parseDayEnd(q("created_end")),
		AssignStart:  parseDay(q("assign_start")), AssignEnd:   parseDayEnd(q("assign_end")),
		HandoverStart: parseDay(q("handover_start")), HandoverEnd: parseDayEnd(q("handover_end")),
		FollowStart:  parseDay(q("follow_start")), FollowEnd:   parseDayEnd(q("follow_end")),
		ArchiveStart: parseDay(q("archive_start")), ArchiveEnd: parseDayEnd(q("archive_end")),
		TaskCreatedStart: parseDay(q("task_created_start")), TaskCreatedEnd: parseDayEnd(q("task_created_end")),
	}
	// 待办卡片下钻快捷参数
	switch q("revisit") {
	case "today":
		t := time.Now()
		s := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		e := s.Add(24 * time.Hour)
		f.Statuses = []string{model.StatusPendingFollow, model.StatusFollowing}
		f.NextFollowStart, f.NextFollowEnd = &s, &e
		f.Status = ""
	case "overdue":
		s := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
		e := time.Now()
		es := time.Date(e.Year(), e.Month(), e.Day(), 0, 0, 0, 0, e.Location())
		f.Statuses = []string{model.StatusPendingFollow, model.StatusFollowing}
		f.NextFollowStart, f.NextFollowEnd = &s, &es
		f.Status = ""
	}
	if dup := q("dup"); dup == "1" || dup == "0" {
		f.Dup = dup
	}
	if q("dup") == "follow" {
		f.Dup = "1"
		f.Status = model.StatusPendingFollow
	}
	return f
}

func (h *Handler) ListLeads(c *gin.Context) {
	f := parseLeadFilter(c)
	res, err := h.svc.ListLeads(c.Request.Context(), f, parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) ExportLeads(c *gin.Context) {
	f := parseLeadFilter(c)
	rows, err := h.svc.ExportLeads(c.Request.Context(), f)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	data, err := xlsxutil.Write(rows)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=leads.xlsx")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *Handler) LeadStat(c *gin.Context) {
	res, err := h.svc.LeadStat(c.Request.Context(), parseLeadFilter(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

type idsBody struct {
	IDs []int64 `json:"ids"`
}

func (h *Handler) HandoverLeads(c *gin.Context) {
	var body struct {
		IDs          []int64 `json:"ids"`
		ReceiverID   string  `json:"receiver_id"`
		ReceiverName string  `json:"receiver_name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Handover(c.Request.Context(), body.IDs, body.ReceiverID, body.ReceiverName, service.CurrentUserName); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) BatchAdjust(c *gin.Context) {
	var body struct {
		IDs    []int64         `json:"ids"`
		Fields map[string]string `json:"fields"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.BatchAdjust(c.Request.Context(), body.IDs, body.Fields, service.CurrentUserName); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) BatchArchive(c *gin.Context) {
	var body struct {
		IDs    []int64 `json:"ids"`
		Rule   string  `json:"rule"`
		Way    string  `json:"way"`
		Reason string  `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.BatchArchive(c.Request.Context(), body.IDs, body.Rule, body.Way, body.Reason); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) ActivateLeads(c *gin.Context) {
	var body idsBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Activate(c.Request.Context(), body.IDs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

// ---------------- 异常线索 ----------------

func (h *Handler) ListExceptions(c *gin.Context) {
	f := repo.ExceptionFilter{
		RelatedName: c.Query("related_name"), RelatedPhone: c.Query("related_phone"),
		ActivityName: c.Query("activity_name"), ActivityPhone: c.Query("activity_phone"),
		StudentNo: c.Query("student_no"),
	}
	res, err := h.svc.ListExceptions(c.Request.Context(), f, parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) IgnoreExceptions(c *gin.Context) {
	var body idsBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.IgnoreExceptions(c.Request.Context(), body.IDs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) AdjustExceptions(c *gin.Context) {
	var body struct {
		IDs    []int64          `json:"ids"`
		Fields map[string]string `json:"fields"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AdjustExceptions(c.Request.Context(), body.IDs, body.Fields); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

// ---------------- 任务模板 / 导入 ----------------

func (h *Handler) ListTaskTemplates(c *gin.Context) {
	list, err := h.svc.ListTaskTemplates(c.Request.Context(), c.Query("keyword"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, list)
}

func (h *Handler) CreateTaskTemplate(c *gin.Context) {
	var f service.ImportForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	t, err := h.svc.SaveTaskTemplate(c.Request.Context(), f)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, t)
}

func (h *Handler) UploadImport(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(20 << 20); err != nil {
		c.JSON(400, gin.H{"error": "解析上传表单失败: " + err.Error()})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "请选择要导入的Excel文件"})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pf := c.PostForm
	activityID, _ := strconv.ParseInt(pf("activity_id"), 10, 64)
	asOpp := pf("as_opportunity") == "1" || pf("as_opportunity") == "true"
	var dedup service.DedupConfig
	if raw := pf("dedup"); raw != "" {
		_ = json.Unmarshal([]byte(raw), &dedup)
	} else {
		dedup.Enable = true
		dedup.OldStudentChecked = true
	}

	form := service.ImportForm{
		TaskName: pf("task_name"), ActivityID: activityID, ActivityName: pf("activity_name"),
		ImportMethod: pf("import_method"), DataSource: pf("data_source"),
		ChannelL1: pf("channel_l1"), ChannelL2: pf("channel_l2"), ChannelL3: pf("channel_l3"),
		SubmitType: pf("submit_type"), AsOpportunity: asOpp, Dedup: dedup,
		FileName: file.Filename,
	}
	res, err := h.svc.ProcessImport(c.Request.Context(), form, data)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) DownloadTemplate(c *gin.Context) {
	method := c.DefaultQuery("method", "phone")
	data, name, err := h.svc.ImportTemplate(method)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// ---------------- 导入日志 ----------------

func (h *Handler) ListImportLogs(c *gin.Context) {
	q := c.Query
	f := repo.ImportTaskFilter{
		ActivityName: q("activity_name"), TaskName: q("task_name"), Operator: q("operator"),
		SerialNo: q("serial_no"), BatchNo: q("batch_no"),
		Start: parseDay(q("start")), End: parseDayEnd(q("end")),
	}
	res, err := h.svc.ListImportLogs(c.Request.Context(), f, parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) DownloadFailFile(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	f, err := h.svc.FailFile(c.Request.Context(), id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", f.Name))
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", f.Content)
}

func (h *Handler) DedupDetails(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	d, err := h.svc.DedupDetails(c.Request.Context(), id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"details": d})
}

// ---------------- 数据进入记录 ----------------

func (h *Handler) ListDataRecords(c *gin.Context) {
	q := c.Query
	f := repo.DataRecordFilter{
		Name: q("name"), Phone: q("phone"), Grade: q("grade"),
		School: q("school"), StudentNo: q("student_no"), Success: q("success"),
	}
	res, err := h.svc.ListDataRecords(c.Request.Context(), f, parsePage(c))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func (h *Handler) HandoverDataRecords(c *gin.Context) {
	var body struct {
		IDs      []int64 `json:"ids"`
		Receiver string  `json:"receiver"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.HandoverDataRecords(c.Request.Context(), body.IDs, body.Receiver); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

func (h *Handler) AdjustDataRecords(c *gin.Context) {
	var body struct {
		IDs    []int64          `json:"ids"`
		Fields map[string]string `json:"fields"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.AdjustDataRecords(c.Request.Context(), body.IDs, body.Fields); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}
