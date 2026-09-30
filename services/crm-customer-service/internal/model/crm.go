package model

import "time"

// 客户
type Customer struct {
	ID           int64      `gorm:"primaryKey" json:"id"`
	TenantID     string     `gorm:"index;size:64" json:"tenant_id"`
	Name         string     `gorm:"size:128" json:"name"`
	ContactPerson string    `gorm:"size:64" json:"contact_person"`
	Phone        string     `gorm:"size:32" json:"phone"`
	Email        string     `gorm:"size:128" json:"email"`
	Address      string     `gorm:"size:256" json:"address"`
	Level        string     `gorm:"size:8" json:"level"`   // A/B/C/D
	Industry     string     `gorm:"size:64" json:"industry"`
	OwnerID      string     `gorm:"size:64" json:"owner_id"`
	Status       int32      `gorm:"default:1" json:"status"` // 1=潜在 2=意向 3=成交 4=流失
	Remark       string     `gorm:"size:512" json:"remark"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Contacts     []Contact  `gorm:"foreignKey:CustomerID" json:"contacts,omitempty"`
	FollowUps    []FollowUp `gorm:"foreignKey:CustomerID" json:"follow_ups,omitempty"`
}

func (Customer) TableName() string { return "crm_customer" }

// 联系人
type Contact struct {
	ID         int64  `gorm:"primaryKey" json:"id"`
	CustomerID int64  `gorm:"index" json:"customer_id"`
	Name       string `gorm:"size:64" json:"name"`
	Position   string `gorm:"size:64" json:"position"`
	Phone      string `gorm:"size:32" json:"phone"`
	Email      string `gorm:"size:128" json:"email"`
	IsPrimary  int32  `gorm:"default:0" json:"is_primary"`
}

func (Contact) TableName() string { return "crm_contact" }

// 跟进记录
type FollowUp struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	CustomerID int64     `gorm:"index" json:"customer_id"`
	Type       string    `gorm:"size:32" json:"type"` // 电话/拜访/邮件/微信
	Content    string    `gorm:"type:text" json:"content"`
	UserID     string    `gorm:"size:64" json:"user_id"`
	FollowUpAt time.Time `json:"follow_up_at"`
}

func (FollowUp) TableName() string { return "crm_follow_up" }
