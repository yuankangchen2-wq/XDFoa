package model

import "time"

// 商机状态
const (
	OppStatusPendingAssign = "待分配"
	OppStatusPendingFollow = "待跟进"
	OppStatusFollowing     = "跟进中"
	OppStatusRecycled      = "已回收"
	OppStatusDeadArchive   = "死单归档"
	OppStatusDealt         = "已成单"
)

// Opportunity 商机（第三期 51 列全集）
type Opportunity struct {
	ID               int64      `gorm:"primaryKey" json:"id"`
	TenantID         string     `gorm:"index;size:64" json:"tenant_id"`
	Code             string     `gorm:"size:64;index" json:"code"` // XSJ...
	Name             string     `gorm:"size:64;index" json:"name"`
	Phone            string     `gorm:"size:32;index" json:"phone"`
	Phone2           string     `gorm:"size:32" json:"phone2"`
	Phone3           string     `gorm:"size:32" json:"phone3"`
	StudentNo        string     `gorm:"size:64;index" json:"student_no"`
	StudentID        string     `gorm:"size:64" json:"student_id"`
	NewOldType       string     `gorm:"size:32" json:"new_old_type"`
	Grade            string     `gorm:"size:32" json:"grade"`
	School           string     `gorm:"size:128" json:"school"`
	District         string     `gorm:"size:64" json:"district"`
	AddWechat        int32      `json:"add_wechat"`
	DeductWechat     int32      `json:"deduct_wechat"`
	Status           string     `gorm:"size:32;index" json:"status"`
	EnterType        string     `gorm:"size:64" json:"enter_type"` // 激活商机/直接创建商机/优秀跟进创建
	OppType          string     `gorm:"size:64" json:"opp_type"`
	Category         string     `gorm:"size:64" json:"category"` // 商机类别（国游专线用）
	IsIntlTravel     int32      `json:"is_intl_travel"`
	// 评分体系
	Level       string `gorm:"size:16" json:"level"` // 意向等级 高/中/低
	FirstLevel  string `gorm:"size:16" json:"first_level"`
	Grade2      string `gorm:"size:4;index" json:"grade2"` // S/A/B/C
	Score       int32  `json:"score"`
	InitScore   int32  `json:"init_score"`
	MaxScore    int32  `json:"max_score"`
	AiDiagnosis string `gorm:"size:256" json:"ai_diagnosis"`
	// 意向
	IntentionDept    string `gorm:"size:128" json:"intention_dept"`
	IntentionProject string `gorm:"size:128" json:"intention_project"`
	IntentionCampus  string `gorm:"size:128" json:"intention_campus"`
	IntentionGrade   string `gorm:"size:32" json:"intention_grade"`
	IntentionSubject string `gorm:"size:64" json:"intention_subject"`
	IntentionClassSize string `gorm:"size:32" json:"intention_class_size"`
	IntentionYear    string `gorm:"size:16" json:"intention_year"`
	IntentionQuarter string `gorm:"size:16" json:"intention_quarter"`
	// 渠道
	ContactPoint    string `gorm:"size:32" json:"contact_point"` // 接触点 呼入/呼出
	ActivityID      int64  `gorm:"index" json:"activity_id"`
	ActivityName    string `gorm:"size:256" json:"activity_name"`
	ChannelL1       string `gorm:"size:64" json:"channel_l1"`
	ChannelL2       string `gorm:"size:64" json:"channel_l2"`
	ChannelL3       string `gorm:"size:64" json:"channel_l3"`
	BusinessChannel string `gorm:"size:64" json:"business_channel"`
	SystemSource    string `gorm:"size:64" json:"system_source"`
	CreateMethod    string `gorm:"size:64" json:"create_method"` // 搜索手工创建/导入
	LeadCode        string `gorm:"size:64;index" json:"lead_code"`
	// 归属与跟进
	OwnerID           string     `gorm:"size:64;index" json:"owner_id"`
	OwnerName         string     `gorm:"size:64" json:"owner_name"`
	OwnerDept         string     `gorm:"size:128" json:"owner_dept"`
	AssignAt          *time.Time `json:"assign_at"`
	FollowCount       int32      `json:"follow_count"`
	LastFollowAt      *time.Time `json:"last_follow_at"`
	LastFollowContent string     `gorm:"size:1024" json:"last_follow_content"`
	LastFollowBy      string     `gorm:"size:64" json:"last_follow_by"`
	NextFollowAt      *time.Time `gorm:"index" json:"next_follow_at"`
	PromiseVisitAt    *time.Time `gorm:"index" json:"promise_visit_at"` // 诺访时间
	VisitConverted    int32      `json:"visit_converted"`
	OverdueDays       int32      `json:"overdue_days"`
	FriendStaff       int32      `json:"friend_staff"`
	// 归档/回收
	ArchiveAt         *time.Time `json:"archive_at"`
	ArchiveRule       string     `gorm:"size:128" json:"archive_rule"`
	PlannedArchiveAt  *time.Time `json:"planned_archive_at"`
	RecycleRule       string     `gorm:"size:128" json:"recycle_rule"`
	PlannedRecycleAt  *time.Time `gorm:"index" json:"planned_recycle_at"`
	RecycleReason     string     `gorm:"size:256" json:"recycle_reason"`
	RecycleAt         *time.Time `json:"recycle_at"`
	DealAt            *time.Time `json:"deal_at"`
	Remark            string     `gorm:"size:1024" json:"remark"`
	Creator           string     `gorm:"size:64" json:"creator"`
	CreatorDept       string     `gorm:"size:128" json:"creator_dept"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	// 外部推送：91国游 / 刺猬
	Pushed91      int32      `json:"pushed_91"`
	Push91Result  string     `gorm:"size:64" json:"push_91_result"`
	ResourceNo91  string     `gorm:"size:64" json:"resource_no_91"`
	Push91At      *time.Time `json:"push_91_at"`
	Follow91      string     `gorm:"size:64" json:"follow_91"`
	Advancer91    string     `gorm:"size:64" json:"advancer_91"`
	PushedHedgehog int32     `json:"pushed_hedgehog"`
}

func (Opportunity) TableName() string { return "scrm_opportunity" }

// OppAssignRecord 商机分配流水（每次流转一条）
type OppAssignRecord struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	OppID        int64     `gorm:"index" json:"opp_id"`
	Code         string    `gorm:"size:64" json:"code"`
	Name         string    `gorm:"size:64" json:"name"`
	StudentNo    string    `gorm:"size:64" json:"student_no"`
	Grade        string    `gorm:"size:32" json:"grade"`
	School       string    `gorm:"size:128" json:"school"`
	Remark       string    `gorm:"size:512" json:"remark"`
	Operator     string    `gorm:"size:64" json:"operator"`
	OperatorDept string    `gorm:"size:128" json:"operator_dept"`
	AutoAssign   int32     `json:"auto_assign"`
	OpAt         time.Time `gorm:"index" json:"op_at"`
	Receiver     string    `gorm:"size:64;index" json:"receiver"`
	ReceiverDept string    `gorm:"size:128" json:"receiver_dept"`
	Status       string    `gorm:"size:32" json:"status"`
	LastFollowAt *time.Time `json:"last_follow_at"`
}

func (OppAssignRecord) TableName() string { return "scrm_opp_assign_record" }

// EmployeeTask 员工任务（任务 × 员工 二维快照）
type EmployeeTask struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	BizKind      string    `gorm:"size:16" json:"biz_kind"` // lead / opp
	TaskID       int64     `gorm:"index" json:"task_id"`
	TaskName     string    `gorm:"size:256" json:"task_name"`
	TaskType     string    `gorm:"size:64" json:"task_type"`
	OwnerID      string    `gorm:"size:64;index" json:"owner_id"`
	OwnerName    string    `gorm:"size:64" json:"owner_name"`
	OwnerDept    string    `gorm:"size:128" json:"owner_dept"`
	Supervisor   string    `gorm:"size:64" json:"supervisor"`
	StartAt      *time.Time `json:"start_at"`
	EndAt        *time.Time `json:"end_at"`
	AutoArchive  int32     `json:"auto_archive"`
	HandoverAt   *time.Time `json:"handover_at"`
	TotalCount   int32     `json:"total_count"`
	StatusCounts string    `gorm:"type:text" json:"status_counts"` // JSON map
}

func (EmployeeTask) TableName() string { return "scrm_employee_task" }

// Appeal 业绩申诉
type Appeal struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	BizLine      string    `gorm:"size:32;index" json:"biz_line"` // 班课/高端/会员卡/课时包
	OppCode      string    `gorm:"size:64;index" json:"opp_code"`
	CertNo       string    `gorm:"size:64" json:"cert_no"` // 听课证号/合同编号/交易编码/交易编号
	OwnerName    string    `gorm:"size:64" json:"owner_name"`
	OwnerDept    string    `gorm:"size:128" json:"owner_dept"`
	Status       string    `gorm:"size:32;index" json:"status"` // 待审批/已通过/已驳回
	Reason       string    `gorm:"size:1024" json:"reason"`
	ImageURL     string    `gorm:"size:512" json:"image_url"`
	AdjustMode   string    `gorm:"size:64" json:"adjust_mode"`
	Applicant    string    `gorm:"size:64" json:"applicant"`
	ApplicantDept string   `gorm:"size:128" json:"applicant_dept"`
	HandledBy    string    `gorm:"size:64" json:"handled_by"`
	CreatedAt    time.Time `json:"created_at"`
	HandledAt    *time.Time `json:"handled_at"`
}

func (Appeal) TableName() string { return "scrm_appeal" }

// DealApproval 申请成单待审
type DealApproval struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	OppCode      string    `gorm:"size:64;index" json:"opp_code"`
	Status       string    `gorm:"size:32;index" json:"status"` // 待审核/已通过/已驳回
	Approver     string    `gorm:"size:64" json:"approver"`
	Applicant    string    `gorm:"size:64" json:"applicant"`
	ApplicantDept string   `gorm:"size:128" json:"applicant_dept"`
	BizType      string    `gorm:"size:32" json:"biz_type"`
	CertNo       string    `gorm:"size:64" json:"cert_no"`
	Amount       string    `gorm:"type:decimal(18,2)" json:"amount"`
	Reason       string    `gorm:"size:1024" json:"reason"`
	ImageURL     string    `gorm:"size:512" json:"image_url"`
	SameStudent  string    `gorm:"size:64" json:"same_student"`
	OrderCreator string    `gorm:"size:64" json:"order_creator"`
	CreatedAt    time.Time `json:"created_at"`
	HandledAt    *time.Time `json:"handled_at"`
}

func (DealApproval) TableName() string { return "scrm_deal_approval" }

// DealLog 成单日志（班课/高端/会员卡/课时包/国游）
type DealLog struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	BizLine      string    `gorm:"size:32;index" json:"biz_line"`
	OppCode      string    `gorm:"size:64;index" json:"opp_code"`
	OwnerName    string    `gorm:"size:64" json:"owner_name"`
	OwnerDept    string    `gorm:"size:128" json:"owner_dept"`
	OrderPerson  string    `gorm:"size:64" json:"order_person"`
	StudentNo    string    `gorm:"size:64" json:"student_no"`
	StudentName  string    `gorm:"size:64" json:"student_name"`
	ProductName  string    `gorm:"size:256" json:"product_name"` // 班级名称/产品
	Amount       string    `gorm:"type:decimal(18,2)" json:"amount"`
	CertNo       string    `gorm:"size:64" json:"cert_no"`
	OrderNo      string    `gorm:"size:64" json:"order_no"`
	ResourceNo   string    `gorm:"size:64" json:"resource_no"`
	DealType     string    `gorm:"size:64" json:"deal_type"`
	PayAt        *time.Time `json:"pay_at"`
	MatchAt      *time.Time `gorm:"index" json:"match_at"` // 订单匹配商机时间
	CreatedAt    time.Time  `json:"created_at"`
}

func (DealLog) TableName() string { return "scrm_deal_log" }

// BizTicket 退调转工单
type BizTicket struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	TenantID      string    `gorm:"index;size:64" json:"tenant_id"`
	Code          string    `gorm:"size:64;index" json:"code"`
	SubmitAccount string    `gorm:"size:64" json:"submit_account"`
	SubmitName    string    `gorm:"size:64" json:"submit_name"`
	SubmitAt      time.Time `gorm:"index" json:"submit_at"`
	BizType       string    `gorm:"size:32" json:"biz_type"` // 退班/调班/转班
	StudentName   string    `gorm:"size:64" json:"student_name"`
	StudentNo     string    `gorm:"size:64" json:"student_no"`
	Status        string    `gorm:"size:32;index" json:"status"` // 待审核/待客服处理/已完成/已驳回
	Source        string    `gorm:"size:64" json:"source"`
	AutoHandle    int32     `json:"auto_handle"`
	CsAccount     string    `gorm:"size:64" json:"cs_account"`
	CreatedAt     time.Time `json:"created_at"`
}

func (BizTicket) TableName() string { return "scrm_biz_ticket" }

// MallOrder 积分商城订单（我的订单）
type MallOrder struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	StudentCode  string    `gorm:"size:64" json:"student_code"`
	StudentName  string    `gorm:"size:64" json:"student_name"`
	ProductType  string    `gorm:"size:64" json:"product_type"`
	SaleType     string    `gorm:"size:32" json:"sale_type"`
	ProductCode  string    `gorm:"size:64" json:"product_code"`
	ProductName  string    `gorm:"size:256" json:"product_name"`
	OrderNo      string    `gorm:"size:64;index" json:"order_no"`
	BatchNo      string    `gorm:"size:64" json:"batch_no"`
	BizType      string    `gorm:"size:32" json:"biz_type"`
	Status       string    `gorm:"size:32" json:"status"`
	OrderAccount string    `gorm:"size:64" json:"order_account"`
	FlagDept     string    `gorm:"size:128" json:"flag_dept"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}

func (MallOrder) TableName() string { return "scrm_mall_order" }

// GridConfig 用户级列/筛选配置
type GridConfig struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TenantID  string    `gorm:"index;size:64" json:"tenant_id"`
	UserID    string    `gorm:"size:64;index:idx_grid,priority:2" json:"user_id"`
	PageKey   string    `gorm:"size:64;index:idx_grid,priority:1" json:"page_key"`
	CfgType   string    `gorm:"size:16;index:idx_grid,priority:3" json:"cfg_type"` // column / filter
	Config    string    `gorm:"type:text" json:"config"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (GridConfig) TableName() string { return "scrm_grid_config" }
