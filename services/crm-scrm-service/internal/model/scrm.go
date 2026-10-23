package model

import "time"

// 线索状态常量
const (
	StatusPendingAssign = "待分配"
	StatusPendingFollow = "待跟进"
	StatusFollowing     = "跟进中"
	StatusPendingHandover = "待交接"
	StatusHandedOver    = "已交接"
	StatusConverted     = "已转商"
	StatusInvalidArchive = "无效归档"
	StatusStoredArchive = "入库归档"
	StatusRecycled      = "已回收"
)

// Activity 活动（线索的批量来源载体）
type Activity struct {
	ID              int64      `gorm:"primaryKey" json:"id"`
	TenantID        string     `gorm:"index;size:64" json:"tenant_id"`
	Name            string     `gorm:"size:256;index" json:"name"`
	SourceSystem    string     `gorm:"size:64" json:"source_system"`     // 来源系统 SCRM
	Category        string     `gorm:"size:64" json:"category"`          // 活动类别
	Type            string     `gorm:"size:64" json:"type"`              // 活动类型
	TargetObject    string     `gorm:"size:64" json:"target_object"`     // 活动对象
	IntentionDept   string     `gorm:"size:128" json:"intention_dept"`
	IntentionProject string    `gorm:"size:128" json:"intention_project"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	Remark          string     `gorm:"size:1024" json:"remark"`
	ExternalID      string     `gorm:"size:64;index" json:"external_id"` // 活动Id 如 3829506
	Code            string     `gorm:"size:64;index" json:"code"`         // 活动编码 HDPC...
	Status          int32      `gorm:"default:0" json:"status"`           // 0未开始 1进行中 2已结束
	Creator         string     `gorm:"size:64" json:"creator"`
	CreatorDept     string     `gorm:"size:128" json:"creator_dept"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (Activity) TableName() string { return "scrm_activity" }

// Lead 线索
type Lead struct {
	ID               int64      `gorm:"primaryKey" json:"id"`
	TenantID         string     `gorm:"index;size:64" json:"tenant_id"`
	Code             string     `gorm:"size:64;index" json:"code"` // 线索编码 XXL...
	Name             string     `gorm:"size:64;index" json:"name"`
	Phone            string     `gorm:"size:32;index" json:"phone"`
	StudentNo        string     `gorm:"size:64;index" json:"student_no"`
	StudentID        string     `gorm:"size:64" json:"student_id"`
	StudentTag       string     `gorm:"size:64" json:"student_tag"` // 学员标签
	NewOldType       string     `gorm:"size:32" json:"new_old_type"`
	School           string     `gorm:"size:128" json:"school"`
	Grade            string     `gorm:"size:32" json:"grade"`            // 在读年级
	District         string     `gorm:"size:64" json:"district"`
	CampusCode       string     `gorm:"size:64" json:"campus_code"`
	LeadType         string     `gorm:"size:64" json:"lead_type"`        // 线索类型
	EnterType        string     `gorm:"size:32" json:"enter_type"`       // 纯新线索 / 激活线索
	Status           string     `gorm:"size:32;index" json:"status"`
	WithdrawStatus   string     `gorm:"size:32" json:"withdraw_status"`
	IntentionLevel   string     `gorm:"size:32" json:"intention_level"`
	IntentionGrade   string     `gorm:"size:32" json:"intention_grade"`
	IntentionDept    string     `gorm:"size:128" json:"intention_dept"`
	IntentionProject string     `gorm:"size:128" json:"intention_project"`
	IntentionCampus  string     `gorm:"size:128" json:"intention_campus"`
	IntentionSubject string     `gorm:"size:64" json:"intention_subject"`
	ChannelL1        string     `gorm:"size:64" json:"channel_l1"`
	ChannelL2        string     `gorm:"size:64" json:"channel_l2"`
	ChannelL3        string     `gorm:"size:64" json:"channel_l3"`
	BusinessChannel  string     `gorm:"size:64" json:"business_channel"`
	SystemSource     string     `gorm:"size:64" json:"system_source"`
	CreateMethod     string     `gorm:"size:64" json:"create_method"`
	ActivityID       int64      `gorm:"index" json:"activity_id"`
	ActivityName     string     `gorm:"size:256" json:"activity_name"`
	TaskID           int64      `gorm:"index" json:"task_id"`
	TaskName         string     `gorm:"size:256" json:"task_name"`
	TaskType         string     `gorm:"size:64" json:"task_type"`
	OwnerID          string     `gorm:"size:64;index" json:"owner_id"`
	OwnerName        string     `gorm:"size:64" json:"owner_name"`
	OwnerDept        string     `gorm:"size:128" json:"owner_dept"`
	AssignTime       *time.Time `json:"assign_time"`
	HandoverTime     *time.Time `json:"handover_time"`
	LastFollowAt     *time.Time `json:"last_follow_at"`
	LastFollowContent string    `gorm:"size:1024" json:"last_follow_content"`
	FollowCount      int32      `json:"follow_count"`
	CallStatus       string     `gorm:"size:32" json:"call_status"`         // 接线状态：未拨打/已接通/未接通
	LastCallAt       *time.Time `json:"last_call_at"`                        // 最近外呼时间
	SelfCallCount    int32      `json:"self_call_count"`                     // 自外呼次数
	TotalCallCount   int32      `json:"total_call_count"`                    // 总外呼次数
	AiCallCount      int32      `json:"ai_call_count"`                       // AI外呼次数
	AiCallAt         *time.Time `json:"ai_call_at"`                          // AI外呼时间
	TotalFollowCount int32      `json:"total_follow_count"`                  // 总跟进次数（全员）
	DeductWechat     int32      `json:"deduct_wechat"`                       // 是否扣微
	NextFollowAt     *time.Time `gorm:"index" json:"next_follow_at"`
	ArchiveTime      *time.Time `json:"archive_time"`
	ArchiveRule      string     `gorm:"size:128" json:"archive_rule"`
	ArchiveWay       string     `gorm:"size:64" json:"archive_way"`
	ArchiveReason    string     `gorm:"size:256" json:"archive_reason"`
	ArchiveOperator  string     `gorm:"size:64" json:"archive_operator"`
	RecycleReason    string     `gorm:"size:256" json:"recycle_reason"`
	RecycleAt        *time.Time `json:"recycle_at"`
	BusinessChannelL2 string    `gorm:"size:64" json:"business_channel_l2"`
	BusinessChannelL3 string    `gorm:"size:64" json:"business_channel_l3"`
	IsDuplicate      int32      `gorm:"default:0;index" json:"is_duplicate"`
	Creator          string     `gorm:"size:64" json:"creator"`
	CreatorDept      string     `gorm:"size:128" json:"creator_dept"`
	Recommender      string     `gorm:"size:64" json:"recommender"`
	Remark           string     `gorm:"size:1024" json:"remark"`
	ImportType       string     `gorm:"size:64" json:"import_type"`
	EnterWay         string     `gorm:"size:64" json:"enter_way"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (Lead) TableName() string { return "scrm_lead" }

// ImportTask 导入任务/批次（一个活动下可多次导入）
type ImportTask struct {
	ID                  int64      `gorm:"primaryKey" json:"id"`
	TenantID            string     `gorm:"index;size:64" json:"tenant_id"`
	TaskName            string     `gorm:"size:256;index" json:"task_name"`
	IsTemplate          int32      `gorm:"default:0" json:"is_template"` // 保存的可复用任务
	ActivityID          int64      `gorm:"index" json:"activity_id"`
	ActivityName        string     `gorm:"size:256" json:"activity_name"`
	ImportMethod        string     `gorm:"size:32" json:"import_method"` // phone / student_no
	DataSource          string     `gorm:"size:64" json:"data_source"`
	ChannelL1           string     `gorm:"size:64" json:"channel_l1"`
	ChannelL2           string     `gorm:"size:64" json:"channel_l2"`
	ChannelL3           string     `gorm:"size:64" json:"channel_l3"`
	SubmitType          string     `gorm:"size:32" json:"submit_type"` // supervisor / employee / none
	AsOpportunity       int32      `json:"as_opportunity"`
	DedupConfig         string     `gorm:"type:text" json:"dedup_config"`
	Operator            string     `gorm:"size:64" json:"operator"`
	SerialNo            string     `gorm:"size:64;index" json:"serial_no"` // 导入流水号
	BatchNo             string     `gorm:"size:64;index" json:"batch_no"`
	Status              string     `gorm:"size:32" json:"status"`          // 进行中 / 已完成
	FileName            string     `gorm:"size:256" json:"file_name"`
	TotalCount          int32      `json:"total_count"`
	SuccessCount        int32      `json:"success_count"`
	NewCount            int32      `json:"new_count"`
	ActiveCount         int32      `json:"active_count"`
	FailCount           int32      `json:"fail_count"`
	DuplicateCount      int32      `json:"duplicate_count"`
	InvalidCount        int32      `json:"invalid_count"`
	UniquenessFailCount int32      `json:"uniqueness_fail_count"`
	DedupHitCount       int32      `json:"dedup_hit_count"`
	DedupConditionCount int32      `json:"dedup_condition_count"`
	DedupDetails        string     `gorm:"type:text" json:"dedup_details"`
	FailRows            string     `gorm:"type:longtext" json:"fail_rows"`
	Deadline            *time.Time `json:"deadline"` // 任务失效时间
	StartedAt           *time.Time `json:"started_at"`
	FinishedAt          *time.Time `json:"finished_at"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (ImportTask) TableName() string { return "scrm_import_task" }

// DataRecord 数据进入记录（逐条数据级别的导入结果追踪）
type DataRecord struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	TenantID         string    `gorm:"index;size:64" json:"tenant_id"`
	ImportTaskID     int64     `gorm:"index" json:"import_task_id"`
	RelatedName      string    `gorm:"size:64;index" json:"related_name"`
	RelatedPhone     string    `gorm:"size:32;index" json:"related_phone"`
	RelatedStudentNo string    `gorm:"size:64" json:"related_student_no"`
	ActivityStudentName string `gorm:"size:64" json:"activity_student_name"`
	ActivityPhone    string    `gorm:"size:32" json:"activity_phone"`
	ImportType       string    `gorm:"size:64" json:"import_type"`
	ActivityGrade    string    `gorm:"size:32" json:"activity_grade"`
	RelatedGrade     string    `gorm:"size:32" json:"related_grade"`
	EnterWay         string    `gorm:"size:64" json:"enter_way"`
	LeadEnterType    string    `gorm:"size:64" json:"lead_enter_type"` // 成功_新建线索 / 失败_线索重复
	FailReason       string    `gorm:"size:1024" json:"fail_reason"`
	ChannelL1        string    `gorm:"size:64" json:"channel_l1"`
	ChannelL2        string    `gorm:"size:64" json:"channel_l2"`
	ChannelL3        string    `gorm:"size:64" json:"channel_l3"`
	ActivityName     string    `gorm:"size:256" json:"activity_name"`
	IntentionDept    string    `gorm:"size:128" json:"intention_dept"`
	IntentionProject string    `gorm:"size:128" json:"intention_project"`
	OwnerName        string    `gorm:"size:64" json:"owner_name"`
	Handled          int32     `gorm:"default:0" json:"handled"`
	CreatedAt        time.Time `json:"created_at"`
}

func (DataRecord) TableName() string { return "scrm_data_record" }

// ExceptionLead 异常线索（数据冲突或信息不完整）
type ExceptionLead struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	TenantID         string    `gorm:"index;size:64" json:"tenant_id"`
	RelatedName      string    `gorm:"size:64;index" json:"related_name"`
	RelatedPhone     string    `gorm:"size:32;index" json:"related_phone"`
	RelatedStudentNo string    `gorm:"size:64" json:"related_student_no"`
	ActivityStudentName string `gorm:"size:64" json:"activity_student_name"`
	ActivityPhone    string    `gorm:"size:32" json:"activity_phone"`
	ChannelL1        string    `gorm:"size:64" json:"channel_l1"`
	ChannelL2        string    `gorm:"size:64" json:"channel_l2"`
	ChannelL3        string    `gorm:"size:64" json:"channel_l3"`
	SystemSourceL1   string    `gorm:"size:64" json:"system_source_l1"`
	SystemSourceL2   string    `gorm:"size:64" json:"system_source_l2"`
	ActivityName     string    `gorm:"size:256" json:"activity_name"`
	DataRemark       string    `gorm:"size:1024" json:"data_remark"`
	IntentionDept    string    `gorm:"size:128" json:"intention_dept"`
	IntentionProject string    `gorm:"size:128" json:"intention_project"`
	Status           string    `gorm:"size:32;index" json:"status"` // 待处理 / 已忽略
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (ExceptionLead) TableName() string { return "scrm_exception_lead" }

// HandoverRecord 交接/批量操作审计记录
type HandoverRecord struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TenantID  string    `gorm:"index;size:64" json:"tenant_id"`
	BizType   string    `gorm:"size:32" json:"biz_type"` // lead / exception / data_record
	BizIDs    string    `gorm:"size:512" json:"biz_ids"`
	Action    string    `gorm:"size:64" json:"action"`
	FromOwner string    `gorm:"size:64" json:"from_owner"`
	ToOwner   string    `gorm:"size:64" json:"to_owner"`
	Operator  string    `gorm:"size:64" json:"operator"`
	Detail    string    `gorm:"size:1024" json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

func (HandoverRecord) TableName() string { return "scrm_handover_record" }

// Schedule 顾问日程
type Schedule struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TenantID  string    `gorm:"index;size:64" json:"tenant_id"`
	UserID    string    `gorm:"index;size:64" json:"user_id"`
	Title     string    `gorm:"size:256" json:"title"`
	Type      string    `gorm:"size:32" json:"type"`       // 日程 / 体验课
	Status    string    `gorm:"size:32" json:"status"`     // 待确访 / 已确访 / 到访未转化
	EventDate time.Time `gorm:"index" json:"event_date"`
	EventTime string    `gorm:"size:16" json:"event_time"`
	CreatedAt time.Time `json:"created_at"`
}

func (Schedule) TableName() string { return "scrm_schedule" }

// ConsultantPerf 顾问成单业绩（金额使用 decimal，禁止 float）
type ConsultantPerf struct {
	ID          int64   `gorm:"primaryKey" json:"id"`
	TenantID    string  `gorm:"index;size:64" json:"tenant_id"`
	Month       string  `gorm:"size:7;index" json:"month"` // 2026-10
	Campus      string  `gorm:"size:64;index" json:"campus"`
	UserID      string  `gorm:"size:64" json:"user_id"`
	Name        string  `gorm:"size:64" json:"name"`
	DealAmount  string  `gorm:"type:decimal(18,2)" json:"deal_amount"`
	DealPersons int32   `json:"deal_persons"`
	DealVisits  int32   `json:"deal_visits"`
}

func (ConsultantPerf) TableName() string { return "scrm_consultant_perf" }

// TodoSeed 待办数量种子（非线索/日程类待办的数量来源）
type TodoSeed struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	TenantID string `gorm:"index;size:64" json:"tenant_id"`
	Role     string `gorm:"size:32;index" json:"role"`
	TodoKey  string `gorm:"size:64;index" json:"todo_key"`
	Count    int32  `json:"count"`
}

func (TodoSeed) TableName() string { return "scrm_todo_seed" }

// FollowLog 线索跟进日志（点击外呼、通话记录均落此表）
type FollowLog struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	LeadID      int64     `gorm:"index" json:"lead_id"`
	LeadName    string    `gorm:"size:64" json:"lead_name"`
	Type        string    `gorm:"size:32;index" json:"type"` // 双呼 / 电话外呼 / 跟进备注
	Content     string    `gorm:"size:1024" json:"content"`
	Result      string    `gorm:"size:32" json:"result"` // 已接通 / 未接通
	DurationSec int32     `json:"duration_sec"`
	CustomerPhone string  `gorm:"size:32" json:"customer_phone"`
	Operator    string    `gorm:"size:64" json:"operator"`
	OperatorID  string    `gorm:"size:64" json:"operator_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (FollowLog) TableName() string { return "scrm_follow_log" }

// Notification 站内通知
type Notification struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TenantID  string    `gorm:"index;size:64" json:"tenant_id"`
	UserID    string    `gorm:"index;size:64" json:"user_id"`
	Type      string    `gorm:"size:32" json:"type"` // todo / handover / task / system
	Title     string    `gorm:"size:128" json:"title"`
	Content   string    `gorm:"size:512" json:"content"`
	BizType   string    `gorm:"size:32" json:"biz_type"` // lead / home / activity ...
	BizID     string    `gorm:"size:64" json:"biz_id"`
	Route     string    `gorm:"size:256" json:"route"`   // 点击跳转路由
	IsRead    int32     `gorm:"default:0;index" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

func (Notification) TableName() string { return "scrm_notification" }
