package repo

import (
	"context"
	"strings"
	"time"

	"github.com/oa-portal/crm-scrm-service/internal/model"
	"gorm.io/gorm"
)

// ---------------- 线索二期：回收 ----------------

func (r *LeadRepo) Recycle(ctx context.Context, tenantID string, ids []int64, reason, operator string, now time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Lead{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{
			"status":         model.StatusRecycled,
			"recycle_reason": reason,
			"recycle_at":     now,
		}).Error
}

func (r *LeadRepo) BatchArchiveWithReason(ctx context.Context, tenantID string, ids []int64, way, rule, reason, operator string, now time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Lead{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{
			"status":           way,
			"archive_way":      way,
			"archive_rule":     rule,
			"archive_reason":   reason,
			"archive_operator": operator,
			"archive_time":     now,
		}).Error
}

// TaskSnapshotRow 任务列表快照
type TaskSnapshotRow struct {
	TaskID     int64  `json:"task_id"`
	TaskName   string `json:"task_name"`
	OwnerID    string `json:"owner_id"`
	OwnerName  string `json:"owner_name"`
	OwnerDept  string `json:"owner_dept"`
	Total      int64  `json:"total"`
	CreatedMax string `json:"created_max"`
}

func (r *LeadRepo) TaskSnapshot(ctx context.Context, tenantID, scope, userID string) ([]TaskSnapshotRow, error) {
	var rows []TaskSnapshotRow
	q := r.db.WithContext(ctx).Model(&model.Lead{}).
		Select("task_id, MAX(task_name) AS task_name, owner_id, MAX(owner_name) AS owner_name, MAX(owner_dept) AS owner_dept, COUNT(*) AS total, MAX(created_at) AS created_max").
		Where("tenant_id = ? AND task_id > 0", tenantID)
	if scope == "mine" {
		q = q.Where("owner_id = ?", userID)
	}
	err := q.Group("task_id, owner_id").Order("task_id DESC").Limit(200).Find(&rows).Error
	return rows, err
}

// ---------------- 商机 ----------------

type OppFilter struct {
	Keyword      string
	Name         string
	Phone        string
	Code         string
	StudentNo    string
	Status       string
	Statuses     []string
	Level        string
	Grade2       string
	EnterType    string
	OppType      string
	Category     string
	NewOld       string
	Grade        string
	District     string
	School       string
	IntentionDept    string
	IntentionProject string
	IntentionCampus  string
	ContactPoint string
	ChannelL1    string
	ChannelL2    string
	ChannelL3    string
	ActivityName string
	OwnerID      string
	OwnerName    string
	Creator      string
	Pushed91     string
	ResourceNo91 string
	IsIntl       string
	AddWechat    string
	Scope        string // mine / dept
	Quick        string // pending / today / recycle3 / promiseweek / hightoday / visitnotconv
	ScoreMin     int
	ScoreMax     int
	InitMin      int
	InitMax      int
	MaxMin       int
	MaxMax       int
	CreatedStart *time.Time
	CreatedEnd   *time.Time
	FollowStart  *time.Time
	FollowEnd    *time.Time
	AssignStart  *time.Time
	AssignEnd    *time.Time
	VisitStart   *time.Time
	VisitEnd     *time.Time
	NextStart    *time.Time
	NextEnd      *time.Time
	ArchiveStart *time.Time
	ArchiveEnd   *time.Time
}

func (f OppFilter) apply(q *gorm.DB) *gorm.DB {
	like := func(s string) string { return "%" + s + "%" }
	if f.Keyword != "" {
		q = q.Where("name LIKE ? OR phone LIKE ? OR code LIKE ? OR student_no LIKE ?", like(f.Keyword), like(f.Keyword), like(f.Keyword), like(f.Keyword))
	}
	if f.Name != "" { q = q.Where("name LIKE ?", like(f.Name)) }
	if f.Phone != "" { q = q.Where("phone LIKE ?", like(f.Phone)) }
	if f.Code != "" { q = q.Where("code LIKE ?", like(f.Code)) }
	if f.StudentNo != "" { q = q.Where("student_no LIKE ?", like(f.StudentNo)) }
	if f.Status != "" { q = q.Where("status = ?", f.Status) }
	if len(f.Statuses) > 0 { q = q.Where("status IN ?", f.Statuses) }
	if f.Level != "" { q = q.Where("level = ?", f.Level) }
	if f.Grade2 != "" { q = q.Where("grade2 = ?", f.Grade2) }
	if f.EnterType != "" { q = q.Where("enter_type = ?", f.EnterType) }
	if f.OppType != "" { q = q.Where("opp_type = ?", f.OppType) }
	if f.Category != "" { q = q.Where("category = ?", f.Category) }
	if f.NewOld != "" { q = q.Where("new_old_type = ?", f.NewOld) }
	if f.Grade != "" { q = q.Where("grade = ?", f.Grade) }
	if f.District != "" { q = q.Where("district = ?", f.District) }
	if f.School != "" { q = q.Where("school LIKE ?", like(f.School)) }
	if f.IntentionDept != "" { q = q.Where("intention_dept = ?", f.IntentionDept) }
	if f.IntentionProject != "" { q = q.Where("intention_project = ?", f.IntentionProject) }
	if f.IntentionCampus != "" { q = q.Where("intention_campus = ?", f.IntentionCampus) }
	if f.ContactPoint != "" { q = q.Where("contact_point = ?", f.ContactPoint) }
	if f.ChannelL1 != "" { q = q.Where("channel_l1 = ?", f.ChannelL1) }
	if f.ChannelL2 != "" { q = q.Where("channel_l2 = ?", f.ChannelL2) }
	if f.ChannelL3 != "" { q = q.Where("channel_l3 = ?", f.ChannelL3) }
	if f.ActivityName != "" { q = q.Where("activity_name LIKE ?", like(f.ActivityName)) }
	if f.OwnerName != "" { q = q.Where("owner_name LIKE ?", like(f.OwnerName)) }
	if f.Creator != "" { q = q.Where("creator LIKE ?", like(f.Creator)) }
	if f.Pushed91 != "" {
		if f.Pushed91 == "1" { q = q.Where("pushed_91 = 1") } else { q = q.Where("pushed_91 = 0") }
	}
	if f.ResourceNo91 != "" { q = q.Where("resource_no_91 LIKE ?", like(f.ResourceNo91)) }
	if f.IsIntl == "1" { q = q.Where("is_intl_travel = 1") }
	if f.AddWechat == "1" { q = q.Where("add_wechat = 1") } else if f.AddWechat == "0" { q = q.Where("add_wechat = 0") }
	if f.ScoreMin > 0 { q = q.Where("score >= ?", f.ScoreMin) }
	if f.ScoreMax > 0 { q = q.Where("score <= ?", f.ScoreMax) }
	if f.InitMin > 0 { q = q.Where("init_score >= ?", f.InitMin) }
	if f.InitMax > 0 { q = q.Where("init_score <= ?", f.InitMax) }
	if f.MaxMin > 0 { q = q.Where("max_score >= ?", f.MaxMin) }
	if f.MaxMax > 0 { q = q.Where("max_score <= ?", f.MaxMax) }
	if f.CreatedStart != nil { q = q.Where("created_at >= ?", *f.CreatedStart) }
	if f.CreatedEnd != nil { q = q.Where("created_at <= ?", *f.CreatedEnd) }
	if f.FollowStart != nil { q = q.Where("last_follow_at >= ?", *f.FollowStart) }
	if f.FollowEnd != nil { q = q.Where("last_follow_at <= ?", *f.FollowEnd) }
	if f.AssignStart != nil { q = q.Where("assign_at >= ?", *f.AssignStart) }
	if f.AssignEnd != nil { q = q.Where("assign_at <= ?", *f.AssignEnd) }
	if f.VisitStart != nil { q = q.Where("promise_visit_at >= ?", *f.VisitStart) }
	if f.VisitEnd != nil { q = q.Where("promise_visit_at <= ?", *f.VisitEnd) }
	if f.NextStart != nil { q = q.Where("next_follow_at >= ?", *f.NextStart) }
	if f.NextEnd != nil { q = q.Where("next_follow_at <= ?", *f.NextEnd) }
	if f.ArchiveStart != nil { q = q.Where("archive_at >= ?", *f.ArchiveStart) }
	if f.ArchiveEnd != nil { q = q.Where("archive_at <= ?", *f.ArchiveEnd) }
	now := time.Now()
	switch f.Quick {
	case "pending":
		q = q.Where("status IN ?", []string{model.OppStatusPendingAssign, model.OppStatusPendingFollow, model.OppStatusFollowing})
	case "today":
		q = q.Where("next_follow_at IS NOT NULL AND DATE(next_follow_at) = ?", now.Format("2006-01-02"))
	case "recycle3":
		q = q.Where("planned_recycle_at IS NOT NULL AND planned_recycle_at BETWEEN ? AND ?", now, now.AddDate(0, 0, 3))
	case "promiseweek":
		wd := int(now.Weekday()); if wd == 0 { wd = 7 }
		mon := time.Date(now.Year(), now.Month(), now.Day()-wd+1, 0, 0, 0, 0, time.Local)
		sun := mon.AddDate(0, 0, 7).Add(-time.Second)
		q = q.Where("promise_visit_at BETWEEN ? AND ? AND status = ?", mon, sun, model.OppStatusFollowing)
	case "hightoday":
		q = q.Where("level = ? AND DATE(last_follow_at) = ?", "高", now.Format("2006-01-02"))
	case "visitnotconv":
		wd := int(now.Weekday()); if wd == 0 { wd = 7 }
		mon := time.Date(now.Year(), now.Month(), now.Day()-wd+1, 0, 0, 0, 0, time.Local)
		sun := mon.AddDate(0, 0, 7).Add(-time.Second)
		q = q.Where("promise_visit_at BETWEEN ? AND ? AND visit_converted = 0 AND status = ?", mon, sun, model.OppStatusFollowing)
	}
	return q
}

type OppRepo struct{ db *gorm.DB }

func NewOppRepo(db *gorm.DB) *OppRepo { return &OppRepo{db: db} }

func (r *OppRepo) base(tenantID, scope, userID string) *gorm.DB {
	q := r.db.Model(&model.Opportunity{}).Where("tenant_id = ?", tenantID)
	if scope == "mine" {
		q = q.Where("owner_id = ?", userID)
	}
	return q
}

func (r *OppRepo) Create(ctx context.Context, o *model.Opportunity) error {
	return r.db.WithContext(ctx).Create(o).Error
}
func (r *OppRepo) Update(ctx context.Context, o *model.Opportunity) error {
	return r.db.WithContext(ctx).Save(o).Error
}
func (r *OppRepo) Get(ctx context.Context, id int64) (*model.Opportunity, error) {
	var o model.Opportunity
	err := r.db.WithContext(ctx).First(&o, id).Error
	return &o, err
}

func (r *OppRepo) List(ctx context.Context, tenantID, scope, userID string, f OppFilter, p Page) (*PageResult[model.Opportunity], error) {
	q := r.base(tenantID, scope, userID)
	if scope == "mine" && f.OwnerID != "" {
		q = q.Where("owner_id = ?", f.OwnerID)
	}
	q = f.apply(q)
	var res PageResult[model.Opportunity]
	if err := q.Count(&res.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil {
		return nil, err
	}
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}

func (r *OppRepo) ListAll(ctx context.Context, tenantID string, f OppFilter, limit int) ([]model.Opportunity, error) {
	if limit <= 0 || limit > 5000 { limit = 500 }
	var list []model.Opportunity
	err := f.apply(r.db.WithContext(ctx).Model(&model.Opportunity{}).Where("tenant_id = ?", tenantID)).
		Order("id DESC").Limit(limit).Find(&list).Error
	return list, err
}

func (r *OppRepo) CountByGrade(ctx context.Context, tenantID string) (map[string]int64, error) {
	type row struct {
		Grade2 string
		Cnt    int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&model.Opportunity{}).
		Select("grade2, COUNT(*) cnt").Where("tenant_id = ?", tenantID).Group("grade2").Find(&rows).Error
	m := map[string]int64{}
	for _, x := range rows { m[x.Grade2] = x.Cnt }
	return m, err
}

func (r *OppRepo) Assign(ctx context.Context, tenantID string, ids []int64, receiverID, receiverName, receiverDept string, auto int32, now time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Opportunity{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{
			"owner_id": receiverID, "owner_name": receiverName, "owner_dept": receiverDept,
			"assign_at": now, "status": model.OppStatusPendingFollow,
		}).Error
}
func (r *OppRepo) Recycle(ctx context.Context, tenantID string, ids []int64, reason string, now time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Opportunity{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{"status": model.OppStatusRecycled, "recycle_reason": reason, "recycle_at": now}).Error
}
func (r *OppRepo) BatchArchive(ctx context.Context, tenantID string, ids []int64, rule string, now time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Opportunity{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{"status": model.OppStatusDeadArchive, "archive_rule": rule, "archive_at": now}).Error
}
func (r *OppRepo) Activate(ctx context.Context, tenantID string, ids []int64) error {
	return r.db.WithContext(ctx).Model(&model.Opportunity{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{"status": model.OppStatusPendingFollow, "archive_at": nil, "recycle_at": nil}).Error
}
func (r *OppRepo) MarkPush91(ctx context.Context, tenantID string, ids []int64, result, resourceNo, follower, advancer string, now time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Opportunity{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{
			"pushed_91": 1, "push_91_result": result, "resource_no_91": resourceNo,
			"push_91_at": now, "follow_91": follower, "advancer_91": advancer,
		}).Error
}
func (r *OppRepo) MarkPushHedgehog(ctx context.Context, tenantID string, ids []int64) error {
	return r.db.WithContext(ctx).Model(&model.Opportunity{}).
		Where("tenant_id = ? AND id IN ?", tenantID).Update("pushed_hedgehog", 1).Error
}

// ---------------- 分配流水 ----------------

type AssignRecordRepo struct{ db *gorm.DB }

func NewAssignRecordRepo(db *gorm.DB) *AssignRecordRepo { return &AssignRecordRepo{db: db} }

type AssignFilter struct {
	Keyword  string
	Receiver string
	Operator string
	Auto     string
	Status   string
	Start    *time.Time
	End      *time.Time
}

func (r *AssignRecordRepo) List(ctx context.Context, tenantID, scope, userID string, f AssignFilter, p Page) (*PageResult[model.OppAssignRecord], error) {
	q := r.db.WithContext(ctx).Model(&model.OppAssignRecord{}).Where("tenant_id = ?", tenantID)
	if scope == "mine" { q = q.Where("receiver = ? OR operator = ?", userID, userID) }
	if f.Keyword != "" {
		q = q.Where("code LIKE ? OR name LIKE ? OR student_no LIKE ?", "%"+f.Keyword+"%", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
	}
	if f.Receiver != "" { q = q.Where("receiver LIKE ?", "%"+f.Receiver+"%") }
	if f.Operator != "" { q = q.Where("operator LIKE ?", "%"+f.Operator+"%") }
	if f.Auto == "1" { q = q.Where("auto_assign = 1") } else if f.Auto == "0" { q = q.Where("auto_assign = 0") }
	if f.Status != "" { q = q.Where("status = ?", f.Status) }
	if f.Start != nil { q = q.Where("op_at >= ?", *f.Start) }
	if f.End != nil { q = q.Where("op_at <= ?", *f.End) }
	var res PageResult[model.OppAssignRecord]
	if err := q.Count(&res.Total).Error; err != nil { return nil, err }
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil { return nil, err }
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}
func (r *AssignRecordRepo) CreateBatch(ctx context.Context, rows []model.OppAssignRecord) error {
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}

// ---------------- 员工任务 ----------------

type EmployeeTaskRepo struct{ db *gorm.DB }

func NewEmployeeTaskRepo(db *gorm.DB) *EmployeeTaskRepo { return &EmployeeTaskRepo{db: db} }

type EmployeeTaskFilter struct {
	OwnerName string
	TaskName  string
	TaskType  string
	BizKind   string
	Start     *time.Time
	End       *time.Time
	Status    string
}

func (r *EmployeeTaskRepo) List(ctx context.Context, tenantID string, f EmployeeTaskFilter, p Page) (*PageResult[model.EmployeeTask], error) {
	q := r.db.WithContext(ctx).Model(&model.EmployeeTask{}).Where("tenant_id = ?", tenantID)
	if f.BizKind != "" { q = q.Where("biz_kind = ?", f.BizKind) }
	if f.OwnerName != "" { q = q.Where("owner_name LIKE ?", "%"+f.OwnerName+"%") }
	if f.TaskName != "" { q = q.Where("task_name LIKE ?", "%"+f.TaskName+"%") }
	if f.TaskType != "" { q = q.Where("task_type = ?", f.TaskType) }
	if f.Start != nil { q = q.Where("start_at >= ?", *f.Start) }
	if f.End != nil { q = q.Where("end_at <= ?", *f.End) }
	var res PageResult[model.EmployeeTask]
	if err := q.Count(&res.Total).Error; err != nil { return nil, err }
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil { return nil, err }
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}
func (r *EmployeeTaskRepo) ReplaceBatch(ctx context.Context, rows []model.EmployeeTask) error {
	if len(rows) == 0 { return nil }
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND biz_kind = ?", rows[0].TenantID, rows[0].BizKind).Delete(&model.EmployeeTask{}).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}

// ---------------- 申诉 / 审批 / 成单 / 工单 / 订单 / 配置 ----------------

type AppealRepo struct{ db *gorm.DB }

func NewAppealRepo(db *gorm.DB) *AppealRepo { return &AppealRepo{db: db} }

func (r *AppealRepo) Create(ctx context.Context, a *model.Appeal) error { return r.db.WithContext(ctx).Create(a).Error }
func (r *AppealRepo) CreateBatch(ctx context.Context, rows []model.Appeal) error {
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}
func (r *AppealRepo) List(ctx context.Context, tenantID, bizLine, keyword, status string, p Page) (*PageResult[model.Appeal], error) {
	q := r.db.WithContext(ctx).Model(&model.Appeal{}).Where("tenant_id = ?", tenantID)
	if bizLine != "" { q = q.Where("biz_line = ?", bizLine) }
	if status != "" { q = q.Where("status = ?", status) }
	if keyword != "" { q = q.Where("opp_code LIKE ? OR cert_no LIKE ?", "%"+keyword+"%", "%"+keyword+"%") }
	var res PageResult[model.Appeal]
	if err := q.Count(&res.Total).Error; err != nil { return nil, err }
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil { return nil, err }
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}
func (r *AppealRepo) Handle(ctx context.Context, tenantID string, ids []int64, pass bool, by string, now time.Time) error {
	st := "已驳回"
	if pass { st = "已通过" }
	return r.db.WithContext(ctx).Model(&model.Appeal{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{"status": st, "handled_by": by, "handled_at": now}).Error
}

type DealApprovalRepo struct{ db *gorm.DB }

func NewDealApprovalRepo(db *gorm.DB) *DealApprovalRepo { return &DealApprovalRepo{db: db} }

func (r *DealApprovalRepo) Create(ctx context.Context, a *model.DealApproval) error { return r.db.WithContext(ctx).Create(a).Error }
func (r *DealApprovalRepo) CreateBatch(ctx context.Context, rows []model.DealApproval) error {
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}
func (r *DealApprovalRepo) List(ctx context.Context, tenantID, keyword, status, applicant, bizType string, p Page) (*PageResult[model.DealApproval], error) {
	q := r.db.WithContext(ctx).Model(&model.DealApproval{}).Where("tenant_id = ?", tenantID)
	if status != "" { q = q.Where("status = ?", status) }
	if keyword != "" { q = q.Where("opp_code LIKE ? OR cert_no LIKE ?", "%"+keyword+"%", "%"+keyword+"%") }
	if applicant != "" { q = q.Where("applicant LIKE ?", "%"+applicant+"%") }
	if bizType != "" { q = q.Where("biz_type = ?", bizType) }
	var res PageResult[model.DealApproval]
	if err := q.Count(&res.Total).Error; err != nil { return nil, err }
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil { return nil, err }
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}
func (r *DealApprovalRepo) Handle(ctx context.Context, tenantID string, ids []int64, pass bool, by string, now time.Time) error {
	st := "已驳回"
	if pass { st = "已通过" }
	return r.db.WithContext(ctx).Model(&model.DealApproval{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{"status": st, "approver": by, "handled_at": now}).Error
}

type DealRepo struct{ db *gorm.DB }

func NewDealRepo(db *gorm.DB) *DealRepo { return &DealRepo{db: db} }

type DealFilter struct {
	BizLine   string
	Keyword   string
	OwnerName string
	Start     *time.Time
	End       *time.Time
}

func (r *DealRepo) List(ctx context.Context, tenantID string, f DealFilter, p Page) (*PageResult[model.DealLog], error) {
	q := r.db.WithContext(ctx).Model(&model.DealLog{}).Where("tenant_id = ?", tenantID)
	if f.BizLine != "" { q = q.Where("biz_line = ?", f.BizLine) }
	if f.Keyword != "" {
		q = q.Where("opp_code LIKE ? OR student_name LIKE ? OR student_no LIKE ? OR cert_no LIKE ? OR order_no LIKE ?",
			"%"+f.Keyword+"%", "%"+f.Keyword+"%", "%"+f.Keyword+"%", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
	}
	if f.OwnerName != "" { q = q.Where("owner_name LIKE ?", "%"+f.OwnerName+"%") }
	if f.Start != nil { q = q.Where("match_at >= ?", *f.Start) }
	if f.End != nil { q = q.Where("match_at <= ?", *f.End) }
	var res PageResult[model.DealLog]
	if err := q.Count(&res.Total).Error; err != nil { return nil, err }
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil { return nil, err }
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}
func (r *DealRepo) Create(ctx context.Context, d *model.DealLog) error { return r.db.WithContext(ctx).Create(d).Error }
func (r *DealRepo) CreateBatch(ctx context.Context, rows []model.DealLog) error {
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}
func (r *DealRepo) SumAmount(ctx context.Context, tenantID, ownerName string) (string, error) {
	var v string
	err := r.db.WithContext(ctx).Model(&model.DealLog{}).
		Where("tenant_id = ? AND (owner_name = ? OR ? = '')", tenantID, ownerName, ownerName).
		Select("COALESCE(SUM(CAST(amount AS DECIMAL(18,2))),0)").Scan(&v).Error
	return v, err
}

type TicketRepo struct{ db *gorm.DB }

func NewTicketRepo(db *gorm.DB) *TicketRepo { return &TicketRepo{db: db} }

func (r *TicketRepo) List(ctx context.Context, tenantID, scope, status, bizType, keyword string, p Page) (*PageResult[model.BizTicket], error) {
	q := r.db.WithContext(ctx).Model(&model.BizTicket{}).Where("tenant_id = ?", tenantID)
	if status != "" { q = q.Where("status = ?", status) }
	if bizType != "" { q = q.Where("biz_type = ?", bizType) }
	if keyword != "" {
		q = q.Where("code LIKE ? OR student_name LIKE ? OR student_no LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if scope == "mine" { q = q.Where("submit_account = ?", "chenyuankang") }
	var res PageResult[model.BizTicket]
	if err := q.Count(&res.Total).Error; err != nil { return nil, err }
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil { return nil, err }
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}
func (r *TicketRepo) CreateBatch(ctx context.Context, rows []model.BizTicket) error {
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}

type OrderRepo struct{ db *gorm.DB }

func NewOrderRepo(db *gorm.DB) *OrderRepo { return &OrderRepo{db: db} }

type OrderFilter struct {
	Keyword     string
	ProductType string
	Status      string
	BizType     string
	SaleType    string
	Start       *time.Time
	End         *time.Time
}

func (r *OrderRepo) List(ctx context.Context, tenantID string, f OrderFilter, p Page) (*PageResult[model.MallOrder], error) {
	q := r.db.WithContext(ctx).Model(&model.MallOrder{}).Where("tenant_id = ? AND order_account = ?", tenantID, "chenyuankang")
	if f.Keyword != "" {
		q = q.Where("student_name LIKE ? OR student_code LIKE ? OR product_code LIKE ? OR order_no LIKE ? OR batch_no LIKE ?",
			"%"+f.Keyword+"%", "%"+f.Keyword+"%", "%"+f.Keyword+"%", "%"+f.Keyword+"%", "%"+f.Keyword+"%")
	}
	if f.ProductType != "" { q = q.Where("product_type = ?", f.ProductType) }
	if f.Status != "" { q = q.Where("status = ?", f.Status) }
	if f.BizType != "" { q = q.Where("biz_type = ?", f.BizType) }
	if f.SaleType != "" { q = q.Where("sale_type = ?", f.SaleType) }
	if f.Start != nil { q = q.Where("created_at >= ?", *f.Start) }
	if f.End != nil { q = q.Where("created_at <= ?", *f.End) }
	var res PageResult[model.MallOrder]
	if err := q.Count(&res.Total).Error; err != nil { return nil, err }
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil { return nil, err }
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}
func (r *OrderRepo) CreateBatch(ctx context.Context, rows []model.MallOrder) error {
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}

type GridConfigRepo struct{ db *gorm.DB }

func NewGridConfigRepo(db *gorm.DB) *GridConfigRepo { return &GridConfigRepo{db: db} }

func (r *GridConfigRepo) Get(ctx context.Context, tenantID, userID, pageKey, cfgType string) (string, error) {
	var g model.GridConfig
	err := r.db.WithContext(ctx).
		Where("tenant_id=? AND user_id=? AND page_key=? AND cfg_type=?", tenantID, userID, pageKey, cfgType).
		First(&g).Error
	if err != nil { return "", err }
	return g.Config, nil
}
func (r *GridConfigRepo) Upsert(ctx context.Context, tenantID, userID, pageKey, cfgType, config string) error {
	var g model.GridConfig
	err := r.db.WithContext(ctx).
		Where("tenant_id=? AND user_id=? AND page_key=? AND cfg_type=?", tenantID, userID, pageKey, cfgType).
		First(&g).Error
	now := time.Now()
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(&model.GridConfig{
			TenantID: tenantID, UserID: userID, PageKey: pageKey, CfgType: cfgType, Config: config, UpdatedAt: now,
		}).Error
	}
	if err != nil { return err }
	return r.db.WithContext(ctx).Model(&g).Updates(map[string]interface{}{"config": config, "updated_at": now}).Error
}

var _ = strings.TrimSpace
