package model

import "time"

// 应收账款
type Receivable struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	TenantID       string    `gorm:"index;size:64" json:"tenant_id"`
	CustomerID     int64     `gorm:"index" json:"customer_id"`
	Amount         float64   `gorm:"type:decimal(18,2)" json:"amount"`
	ReceivedAmount float64   `gorm:"type:decimal(18,2);default:0" json:"received_amount"`
	Status         int32     `gorm:"default:1" json:"status"` // 1=未收 2=部分已收 3=已收清
	DueDate        string    `gorm:"size:32" json:"due_date"`
	Remark         string    `gorm:"size:512" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Receivable) TableName() string { return "fin_receivable" }

// 应付账款
type Payable struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	TenantID   string    `gorm:"index;size:64" json:"tenant_id"`
	SupplierID int64     `gorm:"index" json:"supplier_id"`
	Amount     float64   `gorm:"type:decimal(18,2)" json:"amount"`
	PaidAmount float64   `gorm:"type:decimal(18,2);default:0" json:"paid_amount"`
	Status     int32     `gorm:"default:1" json:"status"` // 1=未付 2=部分已付 3=已付清
	DueDate    string    `gorm:"size:32" json:"due_date"`
	Remark     string    `gorm:"size:512" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Payable) TableName() string { return "fin_payable" }

// 收付款记录
type FinancePayment struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TenantID  string    `gorm:"index;size:64" json:"tenant_id"`
	Type      int32     `json:"type"` // 1=收款 2=付款
	RefID     int64     `gorm:"index" json:"ref_id"`
	Amount    float64   `gorm:"type:decimal(18,2)" json:"amount"`
	PayDate   string    `gorm:"size:32" json:"pay_date"`
	PayMethod string    `gorm:"size:32" json:"pay_method"`
	Remark    string    `gorm:"size:512" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

func (FinancePayment) TableName() string { return "fin_payment" }
