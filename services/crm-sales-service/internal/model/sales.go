package model

import "time"

// 销售订单
type SalesOrder struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	OrderNo     string    `gorm:"size:64;uniqueIndex" json:"order_no"`
	CustomerID  int64     `gorm:"index" json:"customer_id"`
	Status      int32     `gorm:"default:1" json:"status"` // 1=待审批 2=已审批 3=已发货 4=已完成 9=已取消
	TotalAmount float64   `gorm:"type:decimal(18,2)" json:"total_amount"`
	Remark      string    `gorm:"size:512" json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (SalesOrder) TableName() string { return "crm_sales_order" }

// 合同
type Contract struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	TenantID   string    `gorm:"index;size:64" json:"tenant_id"`
	ContractNo string    `gorm:"size:64;uniqueIndex" json:"contract_no"`
	CustomerID int64     `gorm:"index" json:"customer_id"`
	Name       string    `gorm:"size:128" json:"name"`
	Amount     float64   `gorm:"type:decimal(18,2)" json:"amount"`
	PaidAmount float64   `gorm:"type:decimal(18,2);default:0" json:"paid_amount"`
	Status     int32     `gorm:"default:1" json:"status"` // 1=草稿 2=审批中 3=生效 4=归档
	SignDate   string    `gorm:"size:32" json:"sign_date"`
	Remark     string    `gorm:"size:512" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Contract) TableName() string { return "crm_contract" }

// 回款
type Payment struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	TenantID   string    `gorm:"index;size:64" json:"tenant_id"`
	ContractID int64     `gorm:"index" json:"contract_id"`
	Amount     float64   `gorm:"type:decimal(18,2)" json:"amount"`
	PayDate    string    `gorm:"size:32" json:"pay_date"`
	PayMethod  string    `gorm:"size:32" json:"pay_method"`
	Remark     string    `gorm:"size:512" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Payment) TableName() string { return "crm_payment" }
