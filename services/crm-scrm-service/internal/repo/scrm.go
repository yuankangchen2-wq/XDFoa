package repo

import (
	"context"
	"strings"
	"time"

	"github.com/oa-portal/crm-scrm-service/internal/model"
	"gorm.io/gorm"
)

type Page struct {
	Page     int
	PageSize int
}

type PageResult[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

func Paginate(p Page) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.Page <= 0 {
			p.Page = 1
		}
		if p.PageSize <= 0 || p.PageSize > 200 {
			p.PageSize = 20
		}
		return db.Offset((p.Page - 1) * p.PageSize).Limit(p.PageSize)
	}
}

func like(s string) string { return "%" + s + "%" }

// ---------------- 活动 ----------------

type ActivityFilter struct {
	Name         string
	Type         string
	SourceSystem string
	Creator      string
	CreatorDept  string
	Code         string
	ExternalID   string
	LeadStatus   string
	CreatedStart *time.Time
	CreatedEnd   *time.Time
}

type ActivityRepo struct{ db *gorm.DB }

func NewActivityRepo(db *gorm.DB) *ActivityRepo { return &ActivityRepo{db: db} }

func (r *ActivityRepo) Create(ctx context.Context, a *model.Activity) error {
	return r.db.WithContext(ctx).Create(a).Error
}
func (r *ActivityRepo) Update(ctx context.Context, a *model.Activity) error {
	return r.db.WithContext(ctx).Save(a).Error
}
func (r *ActivityRepo) Get(ctx context.Context, id int64) (*model.Activity, error) {
	var a model.Activity
	err := r.db.WithContext(ctx).First(&a, id).Error
	return &a, err
}
func (r *ActivityRepo) All(ctx context.Context, tenantID, keyword string) ([]model.Activity, error) {
	q := r.db.WithContext(ctx).Model(&model.Activity{}).Where("tenant_id = ?", tenantID)
	if keyword != "" {
		q = q.Where("name LIKE ?", like(keyword))
	}
	var list []model.Activity
	err := q.Order("id DESC").Limit(50).Find(&list).Error
	return list, err
}

func (r *ActivityRepo) List(ctx context.Context, tenantID string, f ActivityFilter, p Page) (*PageResult[model.Activity], error) {
	q := r.db.WithContext(ctx).Model(&model.Activity{}).Where("tenant_id = ?", tenantID)
	if f.Name != "" {
		q = q.Where("name LIKE ?", like(f.Name))
	}
	if f.Type != "" {
		q = q.Where("type = ?", f.Type)
	}
	if f.SourceSystem != "" {
		q = q.Where("source_system = ?", f.SourceSystem)
	}
	if f.Creator != "" {
		q = q.Where("creator LIKE ?", like(f.Creator))
	}
	if f.CreatorDept != "" {
		q = q.Where("creator_dept = ?", f.CreatorDept)
	}
	if f.Code != "" {
		q = q.Where("code LIKE ?", like(f.Code))
	}
	if f.ExternalID != "" {
		q = q.Where("external_id LIKE ?", like(f.ExternalID))
	}
	if f.CreatedStart != nil {
		q = q.Where("created_at >= ?", *f.CreatedStart)
	}
	if f.CreatedEnd != nil {
		q = q.Where("created_at <= ?", *f.CreatedEnd)
	}
	if f.LeadStatus != "" {
		q = q.Where("id IN (SELECT activity_id FROM scrm_lead WHERE status = ?)", f.LeadStatus)
	}
	var res PageResult[model.Activity]
	if err := q.Count(&res.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil {
		return nil, err
	}
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}

func (r *ActivityRepo) CountAll(ctx context.Context, tenantID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.Activity{}).Where("tenant_id = ?", tenantID).Count(&n).Error
	return n, err
}

// ---------------- 线索 ----------------

type LeadFilter struct {
	Keyword        string // 姓名/手机/学员号 模糊
	Name           string
	Phone          string
	StudentNo      string
	StudentID      string
	Status         string
	Statuses       []string
	Dup            string // "" / "1" 仅重复 / "0" 仅非重复
	EnterType      string
	LeadType       string
	StudentTag     string
	NewOldType     string
	School         string
	Grade          string
	District       string
	CampusCode     string
	IntentionLevel string
	IntentionGrade string
	IntentionDept  string
	IntentionProject string
	IntentionCampus string
	IntentionSubject string
	ChannelL1      string
	ChannelL2      string
	ChannelL3      string
	BusinessChannel string
	SystemSource   string
	CreateMethod   string
	ActivityID     int64
	ActivityName   string
	TaskID         int64
	OwnerID        string
	OwnerName      string
	Creator        string
	CreatorDept    string
	Recommender    string
	ArchiveRule    string
	ArchiveWay     string
	Remark         string
	FollowMin      int
	FollowMax      int
	NextFollowStart *time.Time
	NextFollowEnd   *time.Time
	CreatedStart   *time.Time
	CreatedEnd     *time.Time
	AssignStart    *time.Time
	AssignEnd      *time.Time
	HandoverStart  *time.Time
	HandoverEnd    *time.Time
	FollowStart    *time.Time
	FollowEnd      *time.Time
	ArchiveStart   *time.Time
	ArchiveEnd     *time.Time
	TaskCreatedStart *time.Time
	TaskCreatedEnd   *time.Time
	View           string // all / handover
}

func (f LeadFilter) apply(q *gorm.DB) *gorm.DB {
	if f.Keyword != "" {
		k := like(f.Keyword)
		q = q.Where("(name LIKE ? OR phone LIKE ? OR student_no LIKE ?)", k, k, k)
	}
	if f.Name != "" {
		q = q.Where("name LIKE ?", like(f.Name))
	}
	if f.Phone != "" {
		q = q.Where("phone LIKE ?", like(f.Phone))
	}
	if f.StudentNo != "" {
		q = q.Where("student_no LIKE ?", like(f.StudentNo))
	}
	if f.StudentID != "" {
		q = q.Where("student_id LIKE ?", like(f.StudentID))
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if len(f.Statuses) > 0 {
		q = q.Where("status IN ?", f.Statuses)
	}
	if f.Dup == "1" {
		q = q.Where("is_duplicate = 1")
	} else if f.Dup == "0" {
		q = q.Where("is_duplicate = 0")
	}
	if f.EnterType != "" {
		q = q.Where("enter_type = ?", f.EnterType)
	}
	if f.LeadType != "" {
		q = q.Where("lead_type = ?", f.LeadType)
	}
	if f.StudentTag != "" {
		q = q.Where("student_tag = ?", f.StudentTag)
	}
	if f.NewOldType != "" {
		q = q.Where("new_old_type = ?", f.NewOldType)
	}
	if f.School != "" {
		q = q.Where("school LIKE ?", like(f.School))
	}
	if f.Grade != "" {
		q = q.Where("grade = ?", f.Grade)
	}
	if f.District != "" {
		q = q.Where("district = ?", f.District)
	}
	if f.CampusCode != "" {
		q = q.Where("campus_code = ?", f.CampusCode)
	}
	if f.IntentionLevel != "" {
		q = q.Where("intention_level = ?", f.IntentionLevel)
	}
	if f.IntentionGrade != "" {
		q = q.Where("intention_grade = ?", f.IntentionGrade)
	}
	if f.IntentionDept != "" {
		q = q.Where("intention_dept = ?", f.IntentionDept)
	}
	if f.IntentionProject != "" {
		q = q.Where("intention_project = ?", f.IntentionProject)
	}
	if f.IntentionCampus != "" {
		q = q.Where("intention_campus = ?", f.IntentionCampus)
	}
	if f.IntentionSubject != "" {
		q = q.Where("intention_subject = ?", f.IntentionSubject)
	}
	if f.ChannelL1 != "" {
		q = q.Where("channel_l1 = ?", f.ChannelL1)
	}
	if f.ChannelL2 != "" {
		q = q.Where("channel_l2 = ?", f.ChannelL2)
	}
	if f.ChannelL3 != "" {
		q = q.Where("channel_l3 = ?", f.ChannelL3)
	}
	if f.BusinessChannel != "" {
		q = q.Where("business_channel = ?", f.BusinessChannel)
	}
	if f.SystemSource != "" {
		q = q.Where("system_source = ?", f.SystemSource)
	}
	if f.CreateMethod != "" {
		q = q.Where("create_method = ?", f.CreateMethod)
	}
	if f.ActivityID > 0 {
		q = q.Where("activity_id = ?", f.ActivityID)
	}
	if f.ActivityName != "" {
		q = q.Where("activity_name LIKE ?", like(f.ActivityName))
	}
	if f.TaskID > 0 {
		q = q.Where("task_id = ?", f.TaskID)
	}
	if f.OwnerID != "" {
		q = q.Where("owner_id = ?", f.OwnerID)
	}
	if f.OwnerName != "" {
		q = q.Where("owner_name LIKE ?", like(f.OwnerName))
	}
	if f.Creator != "" {
		q = q.Where("creator LIKE ?", like(f.Creator))
	}
	if f.CreatorDept != "" {
		q = q.Where("creator_dept = ?", f.CreatorDept)
	}
	if f.Recommender != "" {
		q = q.Where("recommender LIKE ?", like(f.Recommender))
	}
	if f.ArchiveRule != "" {
		q = q.Where("archive_rule = ?", f.ArchiveRule)
	}
	if f.ArchiveWay != "" {
		q = q.Where("archive_way = ?", f.ArchiveWay)
	}
	if f.Remark != "" {
		q = q.Where("remark LIKE ?", like(f.Remark))
	}
	if f.FollowMin > 0 {
		q = q.Where("follow_count >= ?", f.FollowMin)
	}
	if f.FollowMax > 0 {
		q = q.Where("follow_count <= ?", f.FollowMax)
	}
	addRange := func(col string, start, end *time.Time) {
		if start != nil {
			q = q.Where(col+" >= ?", *start)
		}
		if end != nil {
			q = q.Where(col+" <= ?", *end)
		}
	}
	addRange("next_follow_at", f.NextFollowStart, f.NextFollowEnd)
	addRange("created_at", f.CreatedStart, f.CreatedEnd)
	addRange("assign_time", f.AssignStart, f.AssignEnd)
	addRange("handover_time", f.HandoverStart, f.HandoverEnd)
	addRange("last_follow_at", f.FollowStart, f.FollowEnd)
	addRange("archive_time", f.ArchiveStart, f.ArchiveEnd)
	return q
}

type LeadRepo struct{ db *gorm.DB }

func NewLeadRepo(db *gorm.DB) *LeadRepo { return &LeadRepo{db: db} }

func (r *LeadRepo) CreateBatch(ctx context.Context, leads []model.Lead) error {
	if len(leads) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(leads, 200).Error
}
func (r *LeadRepo) Update(ctx context.Context, l *model.Lead) error {
	return r.db.WithContext(ctx).Save(l).Error
}
func (r *LeadRepo) Get(ctx context.Context, id int64) (*model.Lead, error) {
	var l model.Lead
	err := r.db.WithContext(ctx).First(&l, id).Error
	return &l, err
}

// Search 全字段检索：姓名 / 明文电话 / 密文电话（130****3868）/ 线索编码 / 学员号
func (r *LeadRepo) Search(ctx context.Context, tenantID, keyword string, limit int) ([]model.Lead, error) {
	kw := strings.TrimSpace(keyword)
	if kw == "" {
		return nil, nil
	}
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Model(&model.Lead{}).Where("tenant_id = ?", tenantID)
	if strings.Contains(kw, "*") || strings.Contains(kw, "＊") {
		// 密文电话：把星号转成 SQL 通配符做模糊匹配
		pat := strings.NewReplacer("＊", "%", "*", "%").Replace(kw)
		q = q.Where("phone LIKE ?", pat)
	} else {
		q = q.Where("name LIKE ? OR phone LIKE ? OR code LIKE ? OR student_no LIKE ?",
			like(kw), like(kw), like(kw), like(kw))
	}
	var list []model.Lead
	err := q.Order("id DESC").Limit(limit).Find(&list).Error
	return list, err
}

// ApplyCallResult 通话结束后回写线索通话维度字段
func (r *LeadRepo) ApplyCallResult(ctx context.Context, tenantID string, leadID int64, connected bool, at time.Time) error {
	status := "未接通"
	if connected {
		status = "已接通"
	}
	return r.db.WithContext(ctx).Model(&model.Lead{}).
		Where("tenant_id = ? AND id = ?", tenantID, leadID).
		Updates(map[string]interface{}{
			"call_status":     status,
			"last_call_at":    at,
			"self_call_count": gorm.Expr("self_call_count + 1"),
			"total_call_count": gorm.Expr("total_call_count + 1"),
		}).Error
}
func (r *LeadRepo) List(ctx context.Context, tenantID string, f LeadFilter, p Page) (*PageResult[model.Lead], error) {
	q := r.db.WithContext(ctx).Model(&model.Lead{}).Where("tenant_id = ?", tenantID)
	q = f.apply(q)
	if f.View == "handover" {
		q = q.Where("status = ?", model.StatusPendingHandover)
	}
	var res PageResult[model.Lead]
	if err := q.Count(&res.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil {
		return nil, err
	}
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}
func (r *LeadRepo) ListAll(ctx context.Context, tenantID string, f LeadFilter, limit int) ([]model.Lead, error) {
	q := r.db.WithContext(ctx).Model(&model.Lead{}).Where("tenant_id = ?", tenantID)
	q = f.apply(q)
	if limit <= 0 || limit > 10000 {
		limit = 10000
	}
	var list []model.Lead
	err := q.Order("id DESC").Limit(limit).Find(&list).Error
	return list, err
}
func (r *LeadRepo) Count(ctx context.Context, tenantID string, f LeadFilter) (int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Lead{}).Where("tenant_id = ?", tenantID)
	q = f.apply(q)
	var n int64
	err := q.Count(&n).Error
	return n, err
}

func (r *LeadRepo) CountByStatus(ctx context.Context, tenantID string) (map[string]int64, error) {
	type row struct {
		Status string
		Cnt    int64
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&model.Lead{}).
		Select("status, COUNT(*) AS cnt").Where("tenant_id = ?", tenantID).
		Group("status").Scan(&rows).Error
	out := map[string]int64{}
	if err != nil {
		return out, err
	}
	for _, x := range rows {
		out[x.Status] = x.Cnt
	}
	return out, nil
}

func (r *LeadRepo) ExistsByPhone(ctx context.Context, tenantID, phone string) (*model.Lead, bool, error) {
	var l model.Lead
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND phone = ?", tenantID, phone).Order("id DESC").First(&l).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &l, true, nil
}
func (r *LeadRepo) ExistsByStudentNo(ctx context.Context, tenantID, no string) (*model.Lead, bool, error) {
	var l model.Lead
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND student_no = ?", tenantID, no).Order("id DESC").First(&l).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &l, true, nil
}

// Handover 批量交接
func (r *LeadRepo) Handover(ctx context.Context, tenantID string, ids []int64, ownerID, ownerName string, now time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Lead{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{
			"owner_id":      ownerID,
			"owner_name":    ownerName,
			"status":        model.StatusHandedOver,
			"handover_time": now,
		}).Error
}

// BatchAdjust 批量调整信息（白名单字段）
func (r *LeadRepo) BatchAdjust(ctx context.Context, tenantID string, ids []int64, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&model.Lead{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(fields).Error
}

// BatchArchive 批量归档
func (r *LeadRepo) BatchArchive(ctx context.Context, tenantID string, ids []int64, status, rule, way string, now time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Lead{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{
			"status":       status,
			"archive_rule": rule,
			"archive_way":  way,
			"archive_time": now,
		}).Error
}

// Activate 激活已归档/沉睡线索
func (r *LeadRepo) Activate(ctx context.Context, tenantID string, ids []int64) error {
	return r.db.WithContext(ctx).Model(&model.Lead{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{
			"status":       model.StatusPendingAssign,
			"enter_type":   "激活线索",
			"archive_rule": "",
			"archive_way":  "",
			"archive_time": nil,
		}).Error
}

// ---------------- 活动线索统计 ----------------

type ActivityLeadStatRow struct {
	ActivityID int64
	Status     string
	EnterType  string
	Cnt        int64
}

func (r *LeadRepo) GroupByActivity(ctx context.Context, tenantID string) (map[int64]map[string]int64, error) {
	var rows []ActivityLeadStatRow
	err := r.db.WithContext(ctx).Model(&model.Lead{}).
		Select("activity_id, status, enter_type, COUNT(*) AS cnt").
		Where("tenant_id = ? AND activity_id > 0", tenantID).
		Group("activity_id, status, enter_type").
		Scan(&rows).Error
	out := map[int64]map[string]int64{}
	if err != nil {
		return out, err
	}
	for _, x := range rows {
		m := out[x.ActivityID]
		if m == nil {
			m = map[string]int64{}
			out[x.ActivityID] = m
		}
		m["total"] += x.Cnt
		m["status:"+x.Status] += x.Cnt
		m["enter:"+x.EnterType] += x.Cnt
	}
	return out, nil
}

// GradeDistRow 成单年级分布
type GradeDistRow struct {
	Grade string
	Cnt   int64
}

func (r *LeadRepo) GradeDistribution(ctx context.Context, tenantID, campus string, start, end time.Time) ([]GradeDistRow, error) {
	var rows []GradeDistRow
	q := r.db.WithContext(ctx).Model(&model.Lead{}).
		Select("intention_grade AS grade, COUNT(*) AS cnt").
		Where("tenant_id = ? AND status = ? AND created_at >= ? AND created_at <= ?", tenantID, model.StatusConverted, start, end)
	if campus != "" {
		q = q.Where("intention_campus = ?", campus)
	}
	err := q.Group("intention_grade").Order("cnt DESC").Scan(&rows).Error
	return rows, err
}

// ---------------- 导入任务 ----------------

type ImportTaskFilter struct {
	ActivityName string
	TaskName     string
	Operator     string
	SerialNo     string
	BatchNo      string
	Start        *time.Time
	End          *time.Time
}

type ImportTaskRepo struct{ db *gorm.DB }

func NewImportTaskRepo(db *gorm.DB) *ImportTaskRepo { return &ImportTaskRepo{db: db} }

func (r *ImportTaskRepo) Create(ctx context.Context, t *model.ImportTask) error {
	return r.db.WithContext(ctx).Create(t).Error
}
func (r *ImportTaskRepo) Update(ctx context.Context, t *model.ImportTask) error {
	return r.db.WithContext(ctx).Save(t).Error
}
func (r *ImportTaskRepo) Get(ctx context.Context, id int64) (*model.ImportTask, error) {
	var t model.ImportTask
	err := r.db.WithContext(ctx).First(&t, id).Error
	return &t, err
}
func (r *ImportTaskRepo) List(ctx context.Context, tenantID string, f ImportTaskFilter, p Page, templatesOnly bool) (*PageResult[model.ImportTask], error) {
	q := r.db.WithContext(ctx).Model(&model.ImportTask{}).Where("tenant_id = ?", tenantID)
	if templatesOnly {
		q = q.Where("is_template = 1")
	} else {
		q = q.Where("is_template = 0")
	}
	if f.ActivityName != "" {
		q = q.Where("activity_name LIKE ?", like(f.ActivityName))
	}
	if f.TaskName != "" {
		q = q.Where("task_name LIKE ?", like(f.TaskName))
	}
	if f.Operator != "" {
		q = q.Where("operator LIKE ?", like(f.Operator))
	}
	if f.SerialNo != "" {
		q = q.Where("serial_no LIKE ?", like(f.SerialNo))
	}
	if f.BatchNo != "" {
		q = q.Where("batch_no LIKE ?", like(f.BatchNo))
	}
	if f.Start != nil {
		q = q.Where("created_at >= ?", *f.Start)
	}
	if f.End != nil {
		q = q.Where("created_at <= ?", *f.End)
	}
	var res PageResult[model.ImportTask]
	if err := q.Count(&res.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil {
		return nil, err
	}
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}

// ---------------- 数据进入记录 ----------------

type DataRecordFilter struct {
	Name      string
	Phone     string
	Grade     string
	School    string
	StudentNo string
	Success   string // success / fail
}

type DataRecordRepo struct{ db *gorm.DB }

func NewDataRecordRepo(db *gorm.DB) *DataRecordRepo { return &DataRecordRepo{db: db} }

func (r *DataRecordRepo) CreateBatch(ctx context.Context, rows []model.DataRecord) error {
	if len(rows) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}

func (r *DataRecordRepo) List(ctx context.Context, tenantID string, f DataRecordFilter, p Page) (*PageResult[model.DataRecord], error) {
	q := r.db.WithContext(ctx).Model(&model.DataRecord{}).Where("tenant_id = ?", tenantID)
	if f.Name != "" {
		q = q.Where("related_name LIKE ? OR activity_student_name LIKE ?", like(f.Name), like(f.Name))
	}
	if f.Phone != "" {
		q = q.Where("related_phone LIKE ? OR activity_phone LIKE ?", like(f.Phone), like(f.Phone))
	}
	if f.Grade != "" {
		q = q.Where("related_grade = ? OR activity_grade = ?", f.Grade, f.Grade)
	}
	if f.School != "" {
		q = q.Where("activity_name LIKE ?", like(f.School))
	}
	if f.StudentNo != "" {
		q = q.Where("related_student_no LIKE ?", like(f.StudentNo))
	}
	if f.Success == "success" {
		q = q.Where("lead_enter_type LIKE ?", "成功%")
	} else if f.Success == "fail" {
		q = q.Where("lead_enter_type LIKE ?", "失败%")
	}
	var res PageResult[model.DataRecord]
	if err := q.Count(&res.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil {
		return nil, err
	}
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}

func (r *DataRecordRepo) Handover(ctx context.Context, tenantID string, ids []int64, owner string) error {
	return r.db.WithContext(ctx).Model(&model.DataRecord{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(map[string]interface{}{"owner_name": owner, "handled": 1}).Error
}

func (r *DataRecordRepo) Adjust(ctx context.Context, tenantID string, ids []int64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.DataRecord{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(fields).Error
}

// ---------------- 异常线索 ----------------

type ExceptionFilter struct {
	RelatedName  string
	RelatedPhone string
	ActivityName string
	ActivityPhone string
	StudentNo    string
}

type ExceptionRepo struct{ db *gorm.DB }

func NewExceptionRepo(db *gorm.DB) *ExceptionRepo { return &ExceptionRepo{db: db} }

func (r *ExceptionRepo) Create(ctx context.Context, e *model.ExceptionLead) error {
	return r.db.WithContext(ctx).Create(e).Error
}
func (r *ExceptionRepo) List(ctx context.Context, tenantID string, f ExceptionFilter, p Page) (*PageResult[model.ExceptionLead], error) {
	q := r.db.WithContext(ctx).Model(&model.ExceptionLead{}).Where("tenant_id = ?", tenantID)
	if f.RelatedName != "" {
		q = q.Where("related_name LIKE ?", like(f.RelatedName))
	}
	if f.RelatedPhone != "" {
		q = q.Where("related_phone LIKE ?", like(f.RelatedPhone))
	}
	if f.ActivityName != "" {
		q = q.Where("activity_student_name LIKE ?", like(f.ActivityName))
	}
	if f.ActivityPhone != "" {
		q = q.Where("activity_phone LIKE ?", like(f.ActivityPhone))
	}
	if f.StudentNo != "" {
		q = q.Where("related_student_no LIKE ?", like(f.StudentNo))
	}
	var res PageResult[model.ExceptionLead]
	if err := q.Count(&res.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil {
		return nil, err
	}
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}
func (r *ExceptionRepo) CountPending(ctx context.Context, tenantID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.ExceptionLead{}).
		Where("tenant_id = ? AND status = ?", tenantID, "待处理").Count(&n).Error
	return n, err
}
func (r *ExceptionRepo) Ignore(ctx context.Context, tenantID string, ids []int64) error {
	return r.db.WithContext(ctx).Model(&model.ExceptionLead{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Update("status", "已忽略").Error
}

func (r *ExceptionRepo) Adjust(ctx context.Context, tenantID string, ids []int64, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.ExceptionLead{}).
		Where("tenant_id = ? AND id IN ?", tenantID, ids).
		Updates(fields).Error
}

// ---------------- 日程 ----------------

type ScheduleRepo struct{ db *gorm.DB }

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo { return &ScheduleRepo{db: db} }

func (r *ScheduleRepo) Create(ctx context.Context, s *model.Schedule) error {
	return r.db.WithContext(ctx).Create(s).Error
}
func (r *ScheduleRepo) ListByMonth(ctx context.Context, tenantID, userID, month string) ([]model.Schedule, error) {
	var list []model.Schedule
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ? AND DATE_FORMAT(event_date, '%Y-%m') = ?", tenantID, userID, month).
		Order("event_date ASC, event_time ASC").Find(&list).Error
	return list, err
}
func (r *ScheduleRepo) CountByStatus(ctx context.Context, tenantID, userID, status string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.Schedule{}).
		Where("tenant_id = ? AND user_id = ? AND status = ?", tenantID, userID, status).Count(&n).Error
	return n, err
}

// ---------------- 顾问业绩 ----------------

type PerfRepo struct{ db *gorm.DB }

func NewPerfRepo(db *gorm.DB) *PerfRepo { return &PerfRepo{db: db} }

func (r *PerfRepo) Ranking(ctx context.Context, tenantID, month, campus string) ([]model.ConsultantPerf, error) {
	q := r.db.WithContext(ctx).Model(&model.ConsultantPerf{}).
		Where("tenant_id = ? AND month = ?", tenantID, month)
	if campus != "" && campus != "全部校区" {
		q = q.Where("campus = ?", campus)
	}
	var list []model.ConsultantPerf
	err := q.Order("CAST(deal_amount AS DECIMAL(18,2)) DESC").Limit(50).Find(&list).Error
	return list, err
}
func (r *PerfRepo) DistinctUsers(ctx context.Context, tenantID string) ([]model.ConsultantPerf, error) {
	var list []model.ConsultantPerf
	err := r.db.WithContext(ctx).Model(&model.ConsultantPerf{}).
		Select("user_id, MAX(name) AS name, MAX(campus) AS campus").
		Where("tenant_id = ?", tenantID).Group("user_id").Find(&list).Error
	return list, err
}
func (r *PerfRepo) CountAny(ctx context.Context, tenantID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.ConsultantPerf{}).Where("tenant_id = ?", tenantID).Count(&n).Error
	return n, err
}

// ---------------- 待办种子 ----------------

type TodoSeedRepo struct{ db *gorm.DB }

func NewTodoSeedRepo(db *gorm.DB) *TodoSeedRepo { return &TodoSeedRepo{db: db} }

func (r *TodoSeedRepo) ListByRole(ctx context.Context, tenantID, role string) ([]model.TodoSeed, error) {
	var list []model.TodoSeed
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND role = ?", tenantID, role).Find(&list).Error
	return list, err
}
func (r *TodoSeedRepo) CountAny(ctx context.Context, tenantID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.TodoSeed{}).Where("tenant_id = ?", tenantID).Count(&n).Error
	return n, err
}
func (r *TodoSeedRepo) CreateBatch(ctx context.Context, rows []model.TodoSeed) error {
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}

// ---------------- 交接审计 ----------------

type HandoverRepo struct{ db *gorm.DB }

func NewHandoverRepo(db *gorm.DB) *HandoverRepo { return &HandoverRepo{db: db} }

func (r *HandoverRepo) Create(ctx context.Context, rec *model.HandoverRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

// ---------------- 跟进日志（外呼/通话） ----------------

type FollowLogRepo struct{ db *gorm.DB }

func NewFollowLogRepo(db *gorm.DB) *FollowLogRepo { return &FollowLogRepo{db: db} }

func (r *FollowLogRepo) Create(ctx context.Context, l *model.FollowLog) error {
	return r.db.WithContext(ctx).Create(l).Error
}
func (r *FollowLogRepo) Update(ctx context.Context, l *model.FollowLog) error {
	return r.db.WithContext(ctx).Save(l).Error
}
func (r *FollowLogRepo) Get(ctx context.Context, id int64) (*model.FollowLog, error) {
	var l model.FollowLog
	err := r.db.WithContext(ctx).First(&l, id).Error
	return &l, err
}
func (r *FollowLogRepo) ListByLead(ctx context.Context, tenantID string, leadID int64, limit int) ([]model.FollowLog, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var list []model.FollowLog
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND lead_id = ?", tenantID, leadID).
		Order("id DESC").Limit(limit).Find(&list).Error
	return list, err
}
func (r *FollowLogRepo) ListRecent(ctx context.Context, tenantID, operatorID string, limit int) ([]model.FollowLog, error) {
	if limit <= 0 || limit > 200 {
		limit = 20
	}
	var list []model.FollowLog
	q := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if operatorID != "" {
		q = q.Where("operator_id = ?", operatorID)
	}
	err := q.Order("id DESC").Limit(limit).Find(&list).Error
	return list, err
}
func (r *FollowLogRepo) CountAny(ctx context.Context, tenantID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.FollowLog{}).Where("tenant_id = ?", tenantID).Count(&n).Error
	return n, err
}
func (r *FollowLogRepo) CreateBatch(ctx context.Context, rows []model.FollowLog) error {
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}

// ---------------- 站内通知 ----------------

type NotificationFilter struct {
	UserID string
	Unread string // "" 全部 / "1" 未读 / "0" 已读
}

type NotificationRepo struct{ db *gorm.DB }

func NewNotificationRepo(db *gorm.DB) *NotificationRepo { return &NotificationRepo{db: db} }

func (r *NotificationRepo) CreateBatch(ctx context.Context, rows []model.Notification) error {
	return r.db.WithContext(ctx).CreateInBatches(rows, 200).Error
}
func (r *NotificationRepo) List(ctx context.Context, tenantID string, f NotificationFilter, p Page) (*PageResult[model.Notification], error) {
	q := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, f.UserID)
	if f.Unread == "1" {
		q = q.Where("is_read = 0")
	} else if f.Unread == "0" {
		q = q.Where("is_read = 1")
	}
	var res PageResult[model.Notification]
	if err := q.Count(&res.Total).Error; err != nil {
		return nil, err
	}
	if err := q.Order("id DESC").Scopes(Paginate(p)).Find(&res.Items).Error; err != nil {
		return nil, err
	}
	res.Page, res.PageSize = p.Page, p.PageSize
	return &res, nil
}
func (r *NotificationRepo) UnreadCount(ctx context.Context, tenantID, userID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("tenant_id = ? AND user_id = ? AND is_read = 0", tenantID, userID).Count(&n).Error
	return n, err
}
func (r *NotificationRepo) MarkRead(ctx context.Context, tenantID, userID string, ids []int64, all bool) error {
	q := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID)
	if !all {
		q = q.Where("id IN ?", ids)
	}
	return q.Update("is_read", 1).Error
}
func (r *NotificationRepo) CountAny(ctx context.Context, tenantID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.Notification{}).Where("tenant_id = ?", tenantID).Count(&n).Error
	return n, err
}
