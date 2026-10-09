package model

import "time"

// PortalApproval 门户审批
type PortalApproval struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:200;not null" json:"title"`         // 审批标题
	Type        string    `gorm:"size:32;index;not null" json:"type"`     // EHR / 财务 / OA
	SerialNo    string    `gorm:"size:64;index" json:"serial_no"`         // 流水号
	Applicant   string    `gorm:"size:64" json:"applicant"`               // 申请人
	ApplicantID string    `gorm:"size:64;index" json:"applicant_id"`
	Status      int32     `gorm:"default:1;index" json:"status"`           // 1=待办 2=已办
	CurrentNode string    `gorm:"size:128" json:"current_node"`           // 当前节点
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PortalApproval) TableName() string { return "portal_approval" }

// PortalSchedule 我的日程
type PortalSchedule struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"size:64;index" json:"user_id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	EventDate time.Time `gorm:"type:date;index" json:"event_date"`
	CreatedAt time.Time `json:"created_at"`
}

func (PortalSchedule) TableName() string { return "portal_schedule" }

// PortalFavorite 我的常用
type PortalFavorite struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"size:64;index" json:"user_id"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Icon      string    `gorm:"size:64" json:"icon"`
	Path      string    `gorm:"size:256" json:"path"`
	Sort      int       `gorm:"default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
}

func (PortalFavorite) TableName() string { return "portal_favorite" }

// IAMUser iam 库只读视图（员工搜索）
type IAMUser struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int32  `json:"status"`
}

func (IAMUser) TableName() string { return "iam_user" }
