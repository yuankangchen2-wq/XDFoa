package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/oa-portal/crm-scrm-service/internal/dict"
	"github.com/oa-portal/crm-scrm-service/internal/model"
	"github.com/oa-portal/crm-scrm-service/internal/repo"
	"github.com/oa-portal/crm-scrm-service/internal/xlsxutil"
	"gorm.io/gorm"
)

const (
	CurrentTenant   = "default"
	CurrentUserID   = "1"
	CurrentUserName = "陈元素"
)

type ScrmService struct {
	db      *gorm.DB
	actRepo *repo.ActivityRepo
	leadRepo *repo.LeadRepo
	taskRepo *repo.ImportTaskRepo
	recRepo  *repo.DataRecordRepo
	excRepo  *repo.ExceptionRepo
	schedRepo *repo.ScheduleRepo
	perfRepo *repo.PerfRepo
	todoRepo *repo.TodoSeedRepo
	handRepo *repo.HandoverRepo
	logRepo  *repo.FollowLogRepo
	notifyRepo *repo.NotificationRepo
	oppRepo *repo.OppRepo
	assignRepo *repo.AssignRecordRepo
	task2Repo *repo.EmployeeTaskRepo
	appealRepo *repo.AppealRepo
	approvalRepo *repo.DealApprovalRepo
	dealRepo *repo.DealRepo
	ticketRepo *repo.TicketRepo
	orderRepo *repo.OrderRepo
	gridRepo *repo.GridConfigRepo
}

func NewScrmService(db *gorm.DB,
	act *repo.ActivityRepo, lead *repo.LeadRepo, task *repo.ImportTaskRepo,
	rec *repo.DataRecordRepo, exc *repo.ExceptionRepo, sched *repo.ScheduleRepo,
	perf *repo.PerfRepo, todo *repo.TodoSeedRepo, hand *repo.HandoverRepo,
	fl *repo.FollowLogRepo, notify *repo.NotificationRepo,
	opp *repo.OppRepo, assign *repo.AssignRecordRepo, task2 *repo.EmployeeTaskRepo,
	appeal *repo.AppealRepo, approval *repo.DealApprovalRepo, deal *repo.DealRepo,
	ticket *repo.TicketRepo, order *repo.OrderRepo, grid *repo.GridConfigRepo,
) *ScrmService {
	return &ScrmService{db: db, actRepo: act, leadRepo: lead, taskRepo: task,
		recRepo: rec, excRepo: exc, schedRepo: sched, perfRepo: perf, todoRepo: todo, handRepo: hand,
		logRepo: fl, notifyRepo: notify,
		oppRepo: opp, assignRepo: assign, task2Repo: task2, appealRepo: appeal,
		approvalRepo: approval, dealRepo: deal, ticketRepo: ticket, orderRepo: order, gridRepo: grid}
}

// ---------------- 首页：当前用户 ----------------

type MeInfo struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
	Store      string `json:"store"`
	Unread     int    `json:"unread"`
	Roles      []string `json:"roles"`
}

func (s *ScrmService) Me(ctx context.Context) (MeInfo, error) {
	unread, _ := s.notifyRepo.UnreadCount(ctx, CurrentTenant, CurrentUserID)
	return MeInfo{
		Name: CurrentUserName, Title: "课程顾问",
		Phone: "13807110668",
		City: "武汉", Store: "江汉区北诉市场校区", Unread: int(unread),
		Roles: []string{"employee", "supervisor"},
	}, nil
}

func (s *ScrmService) Dicts(kind string) interface{} {
	if kind == "channels" {
		return dict.Channels()
	}
	return dict.Options(kind)
}

// ---------------- 待办聚合 ----------------

type TodoItem struct {
	Key    string `json:"key"`
	Label  string `json:"label"`
	Group  string `json:"group"`
	Count  int64  `json:"count"`
	Route  string `json:"route"`
	Query  string `json:"query,omitempty"`
}

var employeeTodos = []TodoItem{
	// 商机类
	{Key: "opp_follow", Label: "待跟进商机", Group: "商机类", Route: "/opportunities"},
	{Key: "opp_active_revisit", Label: "活跃待回访商机", Group: "商机类", Route: "/opportunities"},
	{Key: "opp_today_revisit", Label: "今日待回访商机", Group: "商机类", Route: "/opportunities"},
	{Key: "opp_overdue", Label: "逾期未回访商机", Group: "商机类", Route: "/opportunities"},
	{Key: "opp_appeal_pending", Label: "申诉待审商机", Group: "商机类", Route: "/opportunities"},
	{Key: "opp_3day_recycle", Label: "3天即回收商机", Group: "商机类", Route: "/opportunities"},
	{Key: "opp_3day_archive", Label: "3天即归档商机", Group: "商机类", Route: "/opportunities"},
	// 体验课类
	{Key: "trial_confirm", Label: "体验课待确访", Group: "体验课类", Route: "/home"},
	{Key: "trial_visit_notconv", Label: "体验课到访未转化", Group: "体验课类", Route: "/home"},
	// 日程类
	{Key: "schedule_confirm", Label: "日程待确访", Group: "日程类", Route: "/home"},
	{Key: "schedule_visit_notconv", Label: "日程到访未转化", Group: "日程类", Route: "/home"},
	// 线索类
	{Key: "lead_follow", Label: "待跟进线索", Group: "线索类", Route: "/all-leads", Query: "status=待跟进"},
	{Key: "lead_today_revisit", Label: "今日待回访线索", Group: "线索类", Route: "/all-leads", Query: "revisit=today"},
	{Key: "lead_overdue", Label: "逾期线索", Group: "线索类", Route: "/all-leads", Query: "revisit=overdue"},
	{Key: "lead_dup_follow", Label: "待跟进重复线索", Group: "线索类", Route: "/all-leads", Query: "dup=follow"},
	{Key: "lead_dup_handover", Label: "待交接重复线索", Group: "线索类", Route: "/handover-leads", Query: "dup=1"},
	{Key: "lead_handover", Label: "待交接线索", Group: "线索类", Route: "/handover-leads"},
	{Key: "lead_exception", Label: "异常线索", Group: "线索类", Route: "/exception-leads"},
	// 任务/其他
	{Key: "task_expiring_3d", Label: "3天即将失效任务", Group: "任务/其他", Route: "/data-record"},
	{Key: "ticket_refund_follow", Label: "待跟进退挽工单", Group: "任务/其他", Route: "/coming-soon", Query: "name=待跟进退挽工单"},
	{Key: "new_sign_schedule", Label: "新签待排课", Group: "任务/其他", Route: "/coming-soon", Query: "name=新签待排课"},
	{Key: "first_lesson", Label: "首课任务", Group: "任务/其他", Route: "/coming-soon", Query: "name=首课任务"},
	{Key: "student_handover", Label: "待交接学员", Group: "任务/其他", Route: "/coming-soon", Query: "name=待交接学员"},
}

var supervisorTodoKeys = []string{"schedule_confirm", "schedule_visit_notconv", "lead_handover"}

func (s *ScrmService) Todos(ctx context.Context, role string) ([]TodoItem, error) {
	var defs []TodoItem
	if role == "supervisor" {
		allow := map[string]bool{}
		for _, k := range supervisorTodoKeys {
			allow[k] = true
		}
		for _, t := range employeeTodos {
			if allow[t.Key] {
				defs = append(defs, t)
			}
		}
	} else {
		defs = employeeTodos
	}

	// 种子数量（非线索/日程类）
	seeds, err := s.todoRepo.ListByRole(ctx, CurrentTenant, role)
	if err != nil {
		return nil, err
	}
	counts := map[string]int64{}
	for _, sd := range seeds {
		counts[sd.TodoKey] = int64(sd.Count)
	}

	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	countLead := func(f repo.LeadFilter) int64 {
		n, _ := s.leadRepo.Count(ctx, CurrentTenant, f)
		return n
	}
	// 实时计算的线索/日程类数量
	counts["lead_follow"] = countLead(repo.LeadFilter{Status: model.StatusPendingFollow})
	counts["lead_today_revisit"] = countLead(repo.LeadFilter{
		Statuses: []string{model.StatusPendingFollow, model.StatusFollowing},
	}) // 用专门方法更精确，下面覆盖
	// 今日回访 / 逾期：按下次回访时间过滤
	n, _ := s.leadRepo.Count(ctx, CurrentTenant, repo.LeadFilter{
		Statuses:        []string{model.StatusPendingFollow, model.StatusFollowing},
		NextFollowStart: &startOfDay, NextFollowEnd: &endOfDay,
	})
	counts["lead_today_revisit"] = n
	overdueStart := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	n, _ = s.leadRepo.Count(ctx, CurrentTenant, repo.LeadFilter{
		Statuses:        []string{model.StatusPendingFollow, model.StatusFollowing},
		NextFollowStart: &overdueStart, NextFollowEnd: &startOfDay,
	})
	counts["lead_overdue"] = n
	counts["lead_dup_follow"] = countLead(repo.LeadFilter{Status: model.StatusPendingFollow})
	// 重复标记需要单独过滤，借助 List 统计
	if dupN, err := s.countDuplicate(ctx, model.StatusPendingFollow); err == nil {
		counts["lead_dup_follow"] = dupN
	}
	if dupN, err := s.countDuplicate(ctx, model.StatusPendingHandover); err == nil {
		counts["lead_dup_handover"] = dupN
	}
	counts["lead_handover"] = countLead(repo.LeadFilter{Status: model.StatusPendingHandover})
	if excN, err := s.excRepo.CountPending(ctx, CurrentTenant); err == nil {
		counts["lead_exception"] = excN
	}
	if v, err := s.schedRepo.CountByStatus(ctx, CurrentTenant, CurrentUserID, "待确访"); err == nil {
		counts["schedule_confirm"] = v
		counts["trial_confirm"] = v
	}
	if v, err := s.schedRepo.CountByStatus(ctx, CurrentTenant, CurrentUserID, "到访未转化"); err == nil {
		counts["schedule_visit_notconv"] = v
		counts["trial_visit_notconv"] = v
	}

	out := make([]TodoItem, 0, len(defs))
	for _, d := range defs {
		d.Count = counts[d.Key]
		out = append(out, d)
	}
	return out, nil
}

func (s *ScrmService) countDuplicate(ctx context.Context, status string) (int64, error) {
	var n int64
	err := s.db.WithContext(ctx).Model(&model.Lead{}).
		Where("tenant_id = ? AND status = ? AND is_duplicate = 1", CurrentTenant, status).Count(&n).Error
	return n, err
}

// ---------------- 顾问排名 / 年级分布 ----------------

func (s *ScrmService) Ranking(ctx context.Context, month, campus string) ([]model.ConsultantPerf, error) {
	return s.perfRepo.Ranking(ctx, CurrentTenant, month, campus)
}

type GradeDist struct {
	Grade string `json:"grade"`
	Count int64  `json:"count"`
}

func (s *ScrmService) GradeDistribution(ctx context.Context, month, campus string) ([]GradeDist, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	start, err := time.ParseInLocation("2006-01", month, time.Local)
	if err != nil {
		return nil, err
	}
	end := start.AddDate(0, 1, 0)
	rows, err := s.leadRepo.GradeDistribution(ctx, CurrentTenant, campus, start, end)
	if err != nil {
		return nil, err
	}
	out := make([]GradeDist, 0, len(rows))
	for _, r := range rows {
		g := r.Grade
		if g == "" {
			g = "未填写"
		}
		out = append(out, GradeDist{Grade: g, Count: r.Cnt})
	}
	return out, nil
}

// ---------------- 日程 ----------------

func (s *ScrmService) ListSchedules(ctx context.Context, month string) ([]model.Schedule, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	return s.schedRepo.ListByMonth(ctx, CurrentTenant, CurrentUserID, month)
}

func (s *ScrmService) CreateSchedule(ctx context.Context, title, date, tm, typ string) (*model.Schedule, error) {
	evDate, err := time.ParseInLocation("2006-01-02", date, time.Local)
	if err != nil {
		return nil, fmt.Errorf("日程日期格式错误")
	}
	sc := &model.Schedule{
		TenantID: CurrentTenant, UserID: CurrentUserID, Title: title,
		Type: typ, Status: "待确访", EventDate: evDate, EventTime: tm,
	}
	if err := s.schedRepo.Create(ctx, sc); err != nil {
		return nil, err
	}
	return sc, nil
}

// ---------------- 员工 ----------------

type Employee struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Campus string `json:"campus"`
}

func (s *ScrmService) Employees(ctx context.Context) ([]Employee, error) {
	users, err := s.perfRepo.DistinctUsers(ctx, CurrentTenant)
	if err != nil {
		return nil, err
	}
	out := []Employee{{ID: "U900", Name: "王主管", Role: "supervisor", Campus: "武汉校区"}}
	for _, u := range users {
		role := "employee"
		out = append(out, Employee{ID: u.UserID, Name: u.Name, Role: role, Campus: u.Campus})
	}
	return out, nil
}

// ---------------- 活动 ----------------

type ActivityView struct {
	model.Activity
	Stats map[string]int64 `json:"stats"`
}

func (s *ScrmService) ListActivities(ctx context.Context, f repo.ActivityFilter, p repo.Page) (*repo.PageResult[ActivityView], error) {
	res, err := s.actRepo.List(ctx, CurrentTenant, f, p)
	if err != nil {
		return nil, err
	}
	statMap, _ := s.leadRepo.GroupByActivity(ctx, CurrentTenant)
	views := &repo.PageResult[ActivityView]{Total: res.Total, Page: res.Page, PageSize: res.PageSize}
	for _, a := range res.Items {
		views.Items = append(views.Items, ActivityView{Activity: a, Stats: statMap[a.ID]})
	}
	return views, nil
}

func (s *ScrmService) AllActivities(ctx context.Context, keyword string) ([]model.Activity, error) {
	return s.actRepo.All(ctx, CurrentTenant, keyword)
}

func (s *ScrmService) CreateActivity(ctx context.Context, a *model.Activity) error {
	a.TenantID = CurrentTenant
	if a.Creator == "" {
		a.Creator = CurrentUserName
	}
	if a.ExternalID == "" {
		a.ExternalID = fmt.Sprintf("%d", 3800000+time.Now().Unix()%100000)
	}
	if a.Code == "" {
		a.Code = "HDPC" + time.Now().Format("010215") + fmt.Sprintf("%04d", time.Now().Unix()%10000)
	}
	if a.SourceSystem == "" {
		a.SourceSystem = "SCRM"
	}
	return s.actRepo.Create(ctx, a)
}

func (s *ScrmService) UpdateActivity(ctx context.Context, a *model.Activity) error {
	a.TenantID = CurrentTenant
	return s.actRepo.Update(ctx, a)
}

// ---------------- 线索 ----------------

func (s *ScrmService) ListLeads(ctx context.Context, f repo.LeadFilter, p repo.Page) (*repo.PageResult[model.Lead], error) {
	return s.leadRepo.List(ctx, CurrentTenant, f, p)
}

func (s *ScrmService) ExportLeads(ctx context.Context, f repo.LeadFilter) ([][]string, error) {
	list, err := s.leadRepo.ListAll(ctx, CurrentTenant, f, 10000)
	if err != nil {
		return nil, err
	}
	rows := [][]string{{"姓名", "手机号", "学员标签", "学员号", "线索编码", "线索状态", "进入类型", "跟进次数", "意向等级", "最近跟进时间", "最近跟进内容", "交接时间", "活动名称"}}
	for _, l := range list {
		rows = append(rows, []string{
			l.Name, l.Phone, l.StudentTag, l.StudentNo, l.Code, l.Status, l.EnterType,
			strconv.Itoa(int(l.FollowCount)), l.IntentionLevel,
			formatTime(l.LastFollowAt), l.LastFollowContent, formatTime(l.HandoverTime), l.ActivityName,
		})
	}
	return rows, nil
}

type LeadStat struct {
	Total  int64            `json:"total"`
	Status map[string]int64 `json:"status"`
}

func (s *ScrmService) LeadStat(ctx context.Context, f repo.LeadFilter) (*LeadStat, error) {
	total, err := s.leadRepo.Count(ctx, CurrentTenant, f)
	if err != nil {
		return nil, err
	}
	status, err := s.leadRepo.CountByStatus(ctx, CurrentTenant)
	if err != nil {
		return nil, err
	}
	return &LeadStat{Total: total, Status: status}, nil
}

var adjustableCols = map[string]bool{
	"student_tag": true, "channel_l1": true, "channel_l2": true, "channel_l3": true,
	"intention_grade": true, "intention_dept": true, "intention_project": true,
	"intention_campus": true, "lead_type": true, "intention_level": true, "remark": true,
}

func (s *ScrmService) BatchAdjust(ctx context.Context, ids []int64, fields map[string]string, operator string) error {
	updates := map[string]interface{}{}
	for k, v := range fields {
		if adjustableCols[k] && v != "" {
			updates[k] = v
		}
	}
	if len(updates) == 0 {
		return fmt.Errorf("没有可调整的字段")
	}
	if err := s.leadRepo.BatchAdjust(ctx, CurrentTenant, ids, updates); err != nil {
		return err
	}
	return s.handRepo.Create(ctx, &model.HandoverRecord{
		TenantID: CurrentTenant, BizType: "lead", BizIDs: joinIDs(ids),
		Action: "批量调整信息", Operator: operator, Detail: fmt.Sprintf("%v", updates),
	})
}

func (s *ScrmService) Handover(ctx context.Context, ids []int64, receiverID, receiverName, operator string) error {
	if len(ids) == 0 || receiverName == "" {
		return fmt.Errorf("请选择线索与接收人")
	}
	if err := s.leadRepo.Handover(ctx, CurrentTenant, ids, receiverID, receiverName, time.Now()); err != nil {
		return err
	}
	return s.handRepo.Create(ctx, &model.HandoverRecord{
		TenantID: CurrentTenant, BizType: "lead", BizIDs: joinIDs(ids),
		Action: "线索交接", ToOwner: receiverName, Operator: operator,
	})
}

func (s *ScrmService) BatchArchive(ctx context.Context, ids []int64, rule, way, reason string) error {
	if len(ids) == 0 {
		return fmt.Errorf("请选择线索")
	}
	status := model.StatusInvalidArchive
	if way == "入库归档" {
		status = model.StatusStoredArchive
	}
	if err := s.leadRepo.BatchArchive(ctx, CurrentTenant, ids, status, rule, way, time.Now()); err != nil {
		return err
	}
	return s.handRepo.Create(ctx, &model.HandoverRecord{
		TenantID: CurrentTenant, BizType: "lead", BizIDs: joinIDs(ids),
		Action: "批量归档", Operator: CurrentUserName, Detail: rule + " / " + reason,
	})
}

func (s *ScrmService) Activate(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("请选择线索")
	}
	if err := s.leadRepo.Activate(ctx, CurrentTenant, ids); err != nil {
		return err
	}
	return s.handRepo.Create(ctx, &model.HandoverRecord{
		TenantID: CurrentTenant, BizType: "lead", BizIDs: joinIDs(ids),
		Action: "激活", Operator: CurrentUserName,
	})
}

func joinIDs(ids []int64) string {
	ps := make([]string, len(ids))
	for i, id := range ids {
		ps[i] = strconv.FormatInt(id, 10)
	}
	return strings.Join(ps, ",")
}

func formatTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// ---------------- 异常线索 ----------------

func (s *ScrmService) ListExceptions(ctx context.Context, f repo.ExceptionFilter, p repo.Page) (*repo.PageResult[model.ExceptionLead], error) {
	return s.excRepo.List(ctx, CurrentTenant, f, p)
}

func (s *ScrmService) IgnoreExceptions(ctx context.Context, ids []int64) error {
	if err := s.excRepo.Ignore(ctx, CurrentTenant, ids); err != nil {
		return err
	}
	return s.handRepo.Create(ctx, &model.HandoverRecord{
		TenantID: CurrentTenant, BizType: "exception", BizIDs: joinIDs(ids),
		Action: "忽略异常数据", Operator: CurrentUserName,
	})
}

// AdjustExceptions 修正异常数据字段（修正后可重新入线索流程）
var exceptionAdjustCols = map[string]bool{
	"related_name": true, "related_phone": true, "related_student_no": true,
	"channel_l1": true, "channel_l2": true, "channel_l3": true,
	"intention_dept": true, "intention_project": true, "data_remark": true,
}

func (s *ScrmService) AdjustExceptions(ctx context.Context, ids []int64, fields map[string]string) error {
	updates := map[string]interface{}{}
	for k, v := range fields {
		if exceptionAdjustCols[k] {
			updates[k] = v
		}
	}
	if len(updates) == 0 {
		return fmt.Errorf("没有可调整的字段")
	}
	if err := s.excRepo.Adjust(ctx, CurrentTenant, ids, updates); err != nil {
		return err
	}
	return s.handRepo.Create(ctx, &model.HandoverRecord{
		TenantID: CurrentTenant, BizType: "exception", BizIDs: joinIDs(ids),
		Action: "批量调整", Operator: CurrentUserName, Detail: fmt.Sprintf("%v", updates),
	})
}

// ---------------- 任务模板 / 导入 ----------------

type DedupConfig struct {
	Enable            bool              `json:"enable"`
	FollowedDurDays   int               `json:"followed_dur_days"` // 指定时间段内跟进过去重（天）
	FollowedChecked  bool              `json:"followed_checked"`
	OldStudentChecked bool              `json:"old_student_checked"`
	OldStudent        map[string]string `json:"old_student"`
	PhoneProtected    bool              `json:"phone_protected"`
	OpportunityDedup  bool              `json:"opportunity_dedup"`
	TaskDedup         bool              `json:"task_dedup"`
}

type ImportForm struct {
	TaskName      string
	ActivityID    int64
	ActivityName  string
	ImportMethod  string
	DataSource    string
	ChannelL1     string
	ChannelL2     string
	ChannelL3     string
	SubmitType    string
	AsOpportunity bool
	Dedup         DedupConfig
	FileName      string
}

func (s *ScrmService) SaveTaskTemplate(ctx context.Context, f ImportForm) (*model.ImportTask, error) {
	cfg, _ := json.Marshal(f.Dedup)
	t := &model.ImportTask{
		TenantID: CurrentTenant, TaskName: f.TaskName, IsTemplate: 1,
		ActivityID: f.ActivityID, ActivityName: f.ActivityName,
		ImportMethod: f.ImportMethod, DataSource: f.DataSource,
		ChannelL1: f.ChannelL1, ChannelL2: f.ChannelL2, ChannelL3: f.ChannelL3,
		SubmitType: f.SubmitType, AsOpportunity: b2i(f.AsOpportunity),
		DedupConfig: string(cfg), Operator: CurrentUserName, Status: "模板",
	}
	if err := s.taskRepo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *ScrmService) ListTaskTemplates(ctx context.Context, keyword string) ([]model.ImportTask, error) {
	res, err := s.taskRepo.List(ctx, CurrentTenant, repo.ImportTaskFilter{TaskName: keyword},
		repo.Page{Page: 1, PageSize: 50}, true)
	if err != nil {
		return nil, err
	}
	return res.Items, nil
}

var importHeaders = []string{"姓名", "手机号", "学员号", "所在学校", "在读年级", "区县", "学员标签",
	"意向年级", "意向部门", "意向项目", "意向校区", "归属人", "备注"}

func (s *ScrmService) ImportTemplate(method string) ([]byte, string, error) {
	headers := importHeaders
	if method == "student_no" {
		headers = []string{"姓名", "学员号", "手机号", "所在学校", "在读年级", "区县", "学员标签",
			"意向年级", "意向部门", "意向项目", "意向校区", "归属人", "备注"}
	}
	demo := []string{"张三", "13800138000", "XDF20260001", "武汉一中", "高三", "江汉区", "高意向",
		"高三", "高中部", "一对一", "武汉校区", "美到810", "地推获取"}
	if method == "student_no" {
		demo = []string{"张三", "XDF20260001", "13800138000", "武汉一中", "高三", "江汉区", "高意向",
			"高三", "高中部", "一对一", "武汉校区", "美到810", "地推获取"}
	}
	data, err := xlsxutil.Write([][]string{headers, demo})
	return data, "线索导入模板.xlsx", err
}

type ImportResult struct {
	TaskID int64 `json:"task_id"`
	model.ImportTask
}

func (s *ScrmService) ProcessImport(ctx context.Context, f ImportForm, fileData []byte) (*ImportResult, error) {
	rows, err := xlsxutil.ReadAll(fileData)
	if err != nil {
		return nil, fmt.Errorf("解析Excel失败: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel中没有数据行")
	}
	header := rows[0]
	idx := map[string]int{}
	for i, h := range header {
		idx[strings.TrimSpace(h)] = i
	}
	get := func(row []string, name string) string {
		if i, ok := idx[name]; ok && i < len(row) {
			return strings.TrimSpace(row[i])
		}
		return ""
	}

	now := time.Now()
	task := &model.ImportTask{
		TenantID: CurrentTenant, TaskName: f.TaskName, IsTemplate: 0,
		ActivityID: f.ActivityID, ActivityName: f.ActivityName,
		ImportMethod: f.ImportMethod, DataSource: f.DataSource,
		ChannelL1: f.ChannelL1, ChannelL2: f.ChannelL2, ChannelL3: f.ChannelL3,
		SubmitType: f.SubmitType, AsOpportunity: b2i(f.AsOpportunity),
		Operator: CurrentUserName, FileName: f.FileName,
		SerialNo: "LS" + now.Format("20060102150405"),
		BatchNo:  "PC" + now.Format("20060102") + fmt.Sprintf("%03d", now.Unix()%1000),
		Status:   "已完成", StartedAt: &now, FinishedAt: &now,
	}
	cfg, _ := json.Marshal(f.Dedup)
	task.DedupConfig = string(cfg)

	var leads []model.Lead
	var records []model.DataRecord
	failRows := [][]string{append(append([]string{}, header...), "失败原因")}
	dedupRuleCount := map[string]int{}

	for rowNo := 1; rowNo < len(rows); rowNo++ {
		row := rows[rowNo]
		if isEmptyRow(row) {
			continue
		}
		task.TotalCount++
		name := get(row, "姓名")
		phone := get(row, "手机号")
		studentNo := get(row, "学员号")
		school := get(row, "所在学校")
		grade := get(row, "在读年级")
		district := get(row, "区县")
		tag := get(row, "学员标签")
		if tag == "" {
			tag = "未知"
		}
		intentionGrade := get(row, "意向年级")
		if intentionGrade == "" {
			intentionGrade = grade
		}
		ownerName := get(row, "归属人")
		remark := get(row, "备注")

		failReason := ""
		if f.ImportMethod == "student_no" && studentNo == "" {
			failReason = "无效数据：学员号导入方式下学员号必填"
		} else if !validPhone(phone) {
			failReason = "无效数据：手机号格式不正确"
		}

		var existing *model.Lead
		if failReason == "" {
			if f.ImportMethod == "student_no" {
				if l, ok, _ := s.leadRepo.ExistsByStudentNo(ctx, CurrentTenant, studentNo); ok {
					existing = l
				}
			} else {
				if l, ok, _ := s.leadRepo.ExistsByPhone(ctx, CurrentTenant, phone); ok {
					existing = l
				}
			}
		}

		enterType := "纯新线索"
		isDup := int32(0)
		if existing != nil {
			if existing.Status == model.StatusInvalidArchive || existing.Status == model.StatusStoredArchive || existing.Status == model.StatusRecycled {
				enterType = "激活线索"
				dedupRuleCount["老学生去重-归档学员重新激活"]++
			} else if f.Dedup.Enable {
				// 命中重复
				failReason = fmt.Sprintf("学员【%s】已存在线索（编码%s，状态%s），不符合线索唯一性", name, existing.Code, existing.Status)
				isDup = 1
				dedupRuleCount["手机号/学员号与库内线索重复"]++
				// 学员号同时指向不同人：计入唯一性失败
				if studentNo != "" {
					if other, ok, _ := s.leadRepo.ExistsByStudentNo(ctx, CurrentTenant, studentNo); ok && other.ID != existing.ID {
						task.UniquenessFailCount++
						failReason = "不符合线索唯一性：手机号与学员号指向不同学员"
					}
				}
			} else {
				isDup = 1
			}
		}

		record := model.DataRecord{
			TenantID: CurrentTenant, RelatedName: name, RelatedPhone: phone,
			RelatedStudentNo: studentNo, ActivityStudentName: name, ActivityPhone: phone,
			ImportType: orDefault(f.DataSource, "市场地推"),
			ActivityGrade: grade, RelatedGrade: grade, EnterWay: "Excel上传",
			ChannelL1: f.ChannelL1, ChannelL2: f.ChannelL2, ChannelL3: f.ChannelL3,
			ActivityName: f.ActivityName,
			IntentionDept: get(row, "意向部门"), IntentionProject: get(row, "意向项目"),
			OwnerName: ownerName,
		}

		if failReason != "" {
			task.FailCount++
			if strings.Contains(failReason, "重复") || strings.Contains(failReason, "唯一性") {
				task.DuplicateCount++
				if strings.Contains(failReason, "去重规则") || strings.Contains(failReason, "已存在线索") {
					task.DedupHitCount++
				}
			} else {
				task.InvalidCount++
			}
			record.LeadEnterType = "失败_线索重复"
			if strings.HasPrefix(failReason, "无效数据") {
				record.LeadEnterType = "失败_无效数据"
			}
			record.FailReason = failReason
			failRow := append(append([]string{}, row...), failReason)
			failRows = append(failRows, failRow)
			records = append(records, record)
			continue
		}

		// 分配规则
		status := model.StatusPendingAssign
		resolvedOwnerID := ""
		if f.SubmitType != "none" && ownerName != "" {
			// Excel 归属人列优先：直接分配到所属员工
			resolvedOwnerID = "U_" + ownerName
			status = model.StatusPendingFollow
		} else if f.SubmitType == "supervisor" {
			status = model.StatusPendingHandover
		} else if f.SubmitType == "employee" {
			status = model.StatusPendingAssign
		}
		assignTime := &now
		if status == model.StatusPendingAssign || status == model.StatusPendingHandover {
			assignTime = nil
		}

		code := fmt.Sprintf("XXL%s%06d", now.Format("200601"), int(task.TotalCount)*100+rowNo)
		lead := model.Lead{
			TenantID: CurrentTenant, Code: code, Name: name, Phone: phone, StudentNo: studentNo,
			StudentTag: tag, NewOldType: "新生", School: school, Grade: grade, District: district,
			LeadType: "普通线索", EnterType: enterType, Status: status,
			IntentionGrade: intentionGrade, IntentionLevel: "D",
			IntentionDept: get(row, "意向部门"), IntentionProject: get(row, "意向项目"),
			IntentionCampus: get(row, "意向校区"),
			ChannelL1: f.ChannelL1, ChannelL2: f.ChannelL2, ChannelL3: f.ChannelL3,
			SystemSource: "SCRM", CreateMethod: "Excel上传",
			ActivityID: f.ActivityID, ActivityName: f.ActivityName,
			OwnerID: resolvedOwnerID, OwnerName: ownerName,
			AssignTime: assignTime, IsDuplicate: isDup,
			Creator: CurrentUserName, CreatorDept: "市场部",
			Remark: remark, ImportType: orDefault(f.DataSource, "市场地推"), EnterWay: "Excel上传",
		}
		leads = append(leads, lead)
		record.LeadEnterType = "成功_新建线索"
		if enterType == "激活线索" {
			record.LeadEnterType = "成功_激活线索"
			task.ActiveCount++
		} else {
			task.NewCount++
		}
		task.SuccessCount++
		records = append(records, record)
	}

	// 去重规则明细
	type dedupDetail struct {
		Rule  string `json:"rule"`
		Count int    `json:"count"`
	}
	details := []dedupDetail{}
	for rule, c := range dedupRuleCount {
		details = append(details, dedupDetail{Rule: rule, Count: c})
	}
	sort.Slice(details, func(i, j int) bool { return details[i].Count > details[j].Count })
	dd, _ := json.Marshal(details)
	task.DedupDetails = string(dd)
	task.DedupConditionCount = int32(len(details))
	if task.FailCount > 0 {
		fb, _ := json.Marshal(failRows)
		task.FailRows = string(fb)
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(task).Error; err != nil {
			return err
		}
		for i := range leads {
			leads[i].TaskID = task.ID
			leads[i].TaskName = f.TaskName
		}
		for i := range records {
			records[i].ImportTaskID = task.ID
		}
		if len(leads) > 0 {
			if err := tx.CreateInBatches(leads, 200).Error; err != nil {
				return err
			}
		}
		if len(records) > 0 {
			if err := tx.CreateInBatches(records, 200).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ImportResult{TaskID: task.ID, ImportTask: *task}, nil
}

func isEmptyRow(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}

func validPhone(p string) bool {
	if len(p) != 11 || p[0] != '1' {
		return false
	}
	for _, c := range p {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func b2i(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func orDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

// ---------------- 导入日志 / 数据进入记录 ----------------

func (s *ScrmService) ListImportLogs(ctx context.Context, f repo.ImportTaskFilter, p repo.Page) (*repo.PageResult[model.ImportTask], error) {
	return s.taskRepo.List(ctx, CurrentTenant, f, p, false)
}

type FailFile struct {
	Name    string
	Content []byte
}

func (s *ScrmService) FailFile(ctx context.Context, id int64) (*FailFile, error) {
	t, err := s.taskRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if t.FailRows == "" {
		return nil, fmt.Errorf("该批次没有失败行")
	}
	var rows [][]string
	if err := json.Unmarshal([]byte(t.FailRows), &rows); err != nil {
		return nil, err
	}
	data, err := xlsxutil.Write(rows)
	if err != nil {
		return nil, err
	}
	return &FailFile{Name: "导入失败明细_" + t.SerialNo + ".xlsx", Content: data}, nil
}

func (s *ScrmService) DedupDetails(ctx context.Context, id int64) (string, error) {
	t, err := s.taskRepo.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return t.DedupDetails, nil
}

func (s *ScrmService) ListDataRecords(ctx context.Context, f repo.DataRecordFilter, p repo.Page) (*repo.PageResult[model.DataRecord], error) {
	return s.recRepo.List(ctx, CurrentTenant, f, p)
}

func (s *ScrmService) HandoverDataRecords(ctx context.Context, ids []int64, receiver string) error {
	if err := s.recRepo.Handover(ctx, CurrentTenant, ids, receiver); err != nil {
		return err
	}
	return s.handRepo.Create(ctx, &model.HandoverRecord{
		TenantID: CurrentTenant, BizType: "data_record", BizIDs: joinIDs(ids),
		Action: "交接给校责任人", ToOwner: receiver, Operator: CurrentUserName,
	})
}

var recordAdjustCols = map[string]bool{
	"related_name": true, "related_phone": true, "related_grade": true,
	"channel_l1": true, "channel_l2": true, "channel_l3": true,
	"intention_dept": true, "intention_project": true,
}

func (s *ScrmService) AdjustDataRecords(ctx context.Context, ids []int64, fields map[string]string) error {
	updates := map[string]interface{}{}
	for k, v := range fields {
		if recordAdjustCols[k] {
			updates[k] = v
		}
	}
	if len(updates) == 0 {
		return fmt.Errorf("没有可调整的字段")
	}
	if err := s.recRepo.Adjust(ctx, CurrentTenant, ids, updates); err != nil {
		return err
	}
	return s.handRepo.Create(ctx, &model.HandoverRecord{
		TenantID: CurrentTenant, BizType: "data_record", BizIDs: joinIDs(ids),
		Action: "批量调整信息", Operator: CurrentUserName, Detail: fmt.Sprintf("%v", updates),
	})
}
