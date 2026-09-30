package model

import "time"

// 成本中心
type CostCenter struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TenantID  string    `gorm:"index;size:64" json:"tenant_id"`
	Code      string    `gorm:"size:64;uniqueIndex:idx_cc_code" json:"code"`
	Name      string    `gorm:"size:128" json:"name"`
	Type      int32     `json:"type"` // 1=部门 2=车间 3=其他
	Remark    string    `gorm:"size:512" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (CostCenter) TableName() string { return "cost_center" }

// 产品成本
type ProductCost struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	ProductID    int64     `gorm:"index" json:"product_id"`
	MaterialCost float64   `gorm:"type:decimal(18,4);default:0" json:"material_cost"`
	LaborCost    float64   `gorm:"type:decimal(18,4);default:0" json:"labor_cost"`
	OverheadCost float64   `gorm:"type:decimal(18,4);default:0" json:"overhead_cost"`
	TotalCost    float64   `gorm:"type:decimal(18,4);default:0" json:"total_cost"`
	Quantity     float64   `gorm:"type:decimal(18,4)" json:"quantity"`
	UnitCost     float64   `gorm:"type:decimal(18,4);default:0" json:"unit_cost"`
	Period       string    `gorm:"size:32" json:"period"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (ProductCost) TableName() string { return "product_cost" }

// 成本记录（按工单/产品归集）
type CostRecord struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	ProductID    int64     `gorm:"index" json:"product_id"`
	WorkOrderID  int64     `gorm:"index" json:"work_order_id"`
	CostCenterID int64     `gorm:"index" json:"cost_center_id"`
	Element      int32     `json:"element"` // 1=物料 2=人工 3=制造费用
	Amount       float64   `gorm:"type:decimal(18,4)" json:"amount"`
	Period       string    `gorm:"size:32" json:"period"`
	Status       int32     `gorm:"default:1" json:"status"` // 1=待归集 2=已归集
	Remark       string    `gorm:"size:512" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (CostRecord) TableName() string { return "cost_record" }
