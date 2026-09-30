package model

import "time"

// 采购订单
type PurchaseOrder struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	TenantID     string    `gorm:"index;size:64" json:"tenant_id"`
	OrderNo      string    `gorm:"size:64;uniqueIndex" json:"order_no"`
	SupplierCode string    `gorm:"size:64" json:"supplier_code"`
	WarehouseID  int64     `json:"warehouse_id"`
	Status       int32     `gorm:"default:1" json:"status"` // 1=草稿 2=已审批 3=部分到货 4=全部到货 5=已结算 9=已取消
	TotalAmount  float64   `gorm:"type:decimal(18,2)" json:"total_amount"`
	Remark       string    `gorm:"size:512" json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Items        []PurchaseOrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (PurchaseOrder) TableName() string { return "pur_order" }

// 采购订单明细
type PurchaseOrderItem struct {
	ID           int64   `gorm:"primaryKey" json:"id"`
	OrderID      int64   `gorm:"index" json:"order_id"`
	SkuCode      string  `gorm:"size:64" json:"sku_code"`
	Quantity     int64   `json:"quantity"`
	UnitPrice    float64 `gorm:"type:decimal(18,2)" json:"unit_price"`
	ReceivedQty  int64   `gorm:"default:0" json:"received_qty"`
}

func (PurchaseOrderItem) TableName() string { return "pur_order_item" }

// 到货单
type PurchaseReceipt struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	TenantID  string    `gorm:"index;size:64" json:"tenant_id"`
	ReceiptNo string    `gorm:"size:64;uniqueIndex" json:"receipt_no"`
	OrderID   int64     `gorm:"index" json:"order_id"`
	Status    int32     `gorm:"default:1" json:"status"` // 1=待入库 2=已入库
	CreatedAt time.Time `json:"created_at"`
	Items     []PurchaseReceiptItem `gorm:"foreignKey:ReceiptID" json:"items,omitempty"`
}

func (PurchaseReceipt) TableName() string { return "pur_receipt" }

// 到货明细
type PurchaseReceiptItem struct {
	ID        int64  `gorm:"primaryKey" json:"id"`
	ReceiptID int64  `gorm:"index" json:"receipt_id"`
	SkuCode   string `gorm:"size:64" json:"sku_code"`
	Quantity  int64  `json:"quantity"`
}

func (PurchaseReceiptItem) TableName() string { return "pur_receipt_item" }
