package model

import "time"

// 销售快照（按期间/产品/客户维度）
type SalesSnapshot struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	Period      string    `gorm:"size:32;index" json:"period"` // yyyy-MM
	ProductID   int64     `gorm:"index" json:"product_id"`
	CustomerID  int64     `gorm:"index" json:"customer_id"`
	TotalAmount float64   `gorm:"type:decimal(18,2)" json:"total_amount"`
	OrderCount  int32     `json:"order_count"`
	AvgAmount   float64   `gorm:"type:decimal(18,2)" json:"avg_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

func (SalesSnapshot) TableName() string { return "ana_sales_snapshot" }

// 销售漏斗阶段
type FunnelStage struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	TenantID         string    `gorm:"index;size:64" json:"tenant_id"`
	Period           string    `gorm:"size:32;index" json:"period"`
	Stage            string    `gorm:"size:64" json:"stage"`
	OpportunityCount int32     `json:"opportunity_count"`
	Amount           float64   `gorm:"type:decimal(18,2)" json:"amount"`
	ConversionRate   float64   `gorm:"type:decimal(8,4)" json:"conversion_rate"`
	CreatedAt        time.Time `json:"created_at"`
}

func (FunnelStage) TableName() string { return "ana_funnel_stage" }

// 客户统计
type CustomerStat struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	TenantID         string    `gorm:"index;size:64" json:"tenant_id"`
	Period           string    `gorm:"size:32;index" json:"period"`
	NewCustomers     int32     `json:"new_customers"`
	ActiveCustomers  int32     `json:"active_customers"`
	TotalContribution float64  `gorm:"type:decimal(18,2)" json:"total_contribution"`
	TopCustomerID    int64     `json:"top_customer_id"`
	TopContribution  float64   `gorm:"type:decimal(18,2)" json:"top_contribution"`
	CreatedAt        time.Time `json:"created_at"`
}

func (CustomerStat) TableName() string { return "ana_customer_stat" }
