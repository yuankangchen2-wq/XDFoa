package model

import "time"

// 商机
type Opportunity struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	TenantID         string    `gorm:"index;size:64" json:"tenant_id"`
	Name             string    `gorm:"size:128" json:"name"`
	CustomerID       int64     `gorm:"index" json:"customer_id"`
	Stage            int32     `gorm:"default:1" json:"stage"` // 1=线索 2=意向 3=报价 4=谈判 5=成交 6=输单
	Amount           float64   `gorm:"type:decimal(18,2)" json:"amount"`
	WinRate          int32     `gorm:"default:10" json:"win_rate"` // 赢率 0-100
	ExpectedCloseDate string   `gorm:"size:32" json:"expected_close_date"`
	OwnerID          string    `gorm:"size:64" json:"owner_id"`
	Remark           string    `gorm:"size:512" json:"remark"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (Opportunity) TableName() string { return "crm_opportunity" }
