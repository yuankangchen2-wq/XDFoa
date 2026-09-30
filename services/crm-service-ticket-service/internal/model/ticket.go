package model

import "time"

// 工单
type ServiceTicket struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	TicketNo    string    `gorm:"size:64;uniqueIndex" json:"ticket_no"`
	CustomerID  int64     `gorm:"index" json:"customer_id"`
	Title       string    `gorm:"size:256" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Status      int32     `gorm:"default:1" json:"status"`   // 1=待处理 2=处理中 3=已解决 4=已关闭
	Priority    int32     `gorm:"default:2" json:"priority"` // 1=低 2=中 3=高 4=紧急
	AssigneeID  string    `gorm:"size:64" json:"assignee_id"`
	Satisfaction int32    `gorm:"default:0" json:"satisfaction"` // 满意度 1-5
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Replies     []TicketReply `gorm:"foreignKey:TicketID" json:"replies,omitempty"`
}

func (ServiceTicket) TableName() string { return "crm_service_ticket" }

// 工单回复
type TicketReply struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TicketID  int64     `gorm:"index" json:"ticket_id"`
	UserID    string    `gorm:"size:64" json:"user_id"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (TicketReply) TableName() string { return "crm_ticket_reply" }
