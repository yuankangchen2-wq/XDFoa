package model

import "time"

// BOM 物料清单主表
type Bom struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	TenantID   string    `gorm:"index;size:64" json:"tenant_id"`
	ProductSku string    `gorm:"size:64" json:"product_sku"`
	BomName    string    `gorm:"size:128" json:"bom_name"`
	Status     int32     `gorm:"default:1" json:"status"` // 1=草稿 2=已启用
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Items      []BomItem `gorm:"foreignKey:BomID" json:"items,omitempty"`
}

func (Bom) TableName() string { return "prd_bom" }

// BOM 明细
type BomItem struct {
	ID           int64  `gorm:"primaryKey" json:"id"`
	BomID        int64  `gorm:"index" json:"bom_id"`
	ComponentSku string `gorm:"size:64" json:"component_sku"`
	Quantity     int64  `json:"quantity"`
	Level        int32  `json:"level"`
}

func (BomItem) TableName() string { return "prd_bom_item" }

// 生产工单
type WorkOrder struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	OrderNo     string    `gorm:"size:64;uniqueIndex" json:"order_no"`
	ProductSku  string    `gorm:"size:64" json:"product_sku"`
	Quantity    int64     `json:"quantity"`
	WarehouseID int64     `json:"warehouse_id"`
	Status      int32     `gorm:"default:1" json:"status"` // 1=待领料 2=生产中 3=已完工 9=已取消
	ProducedQty int64     `gorm:"default:0" json:"produced_qty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (WorkOrder) TableName() string { return "prd_work_order" }

// 领料单
type MaterialIssue struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	IssueNo     string    `gorm:"size:64;uniqueIndex" json:"issue_no"`
	WorkOrderID int64     `gorm:"index" json:"work_order_id"`
	Status      int32     `gorm:"default:1" json:"status"` // 1=待出库 2=已出库
	CreatedAt   time.Time `json:"created_at"`
	Items       []MaterialIssueItem `gorm:"foreignKey:IssueID" json:"items,omitempty"`
}

func (MaterialIssue) TableName() string { return "prd_material_issue" }

// 领料明细
type MaterialIssueItem struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	IssueID  int64  `gorm:"index" json:"issue_id"`
	SkuCode  string `gorm:"size:64" json:"sku_code"`
	Quantity int64  `json:"quantity"`
}

func (MaterialIssueItem) TableName() string { return "prd_material_issue_item" }

// 报工记录
type ProductionReport struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	TenantID    string    `gorm:"index;size:64" json:"tenant_id"`
	WorkOrderID int64     `gorm:"index" json:"work_order_id"`
	ProducedQty int64     `json:"produced_qty"`
	ReportNo    string    `gorm:"size:64" json:"report_no"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ProductionReport) TableName() string { return "prd_production_report" }
