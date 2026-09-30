package model

import "time"

// 审计日志
type AuditLog struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	TenantID      string    `gorm:"index;size:64" json:"tenant_id"`
	UserID        string    `gorm:"size:64;index" json:"user_id"`
	Username      string    `gorm:"size:128" json:"username"`
	Action        string    `gorm:"size:64;index" json:"action"`       // CREATE/UPDATE/DELETE/LOGIN/...
	Module        string    `gorm:"size:64;index" json:"module"`       // user/customer/order/...
	TargetID      string    `gorm:"size:128" json:"target_id"`
	TargetName    string    `gorm:"size:256" json:"target_name"`
	RequestParams string    `gorm:"type:text" json:"request_params"`   // JSON
	Result        string    `gorm:"type:text" json:"result"`
	IP            string    `gorm:"size:64" json:"ip"`
	Status        int32     `gorm:"index" json:"status"`               // 1=成功 2=失败
	DurationMs    int64     `json:"duration_ms"`
	CreatedAt     time.Time `gorm:"index" json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_log" }
