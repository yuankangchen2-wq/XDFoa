package model

import "time"

// 订单
type Order struct {
	ID           int64      `gorm:"primaryKey" json:"id"`
	TenantID     string     `gorm:"index;size:64" json:"tenant_id"`
	OrderNo      string     `gorm:"size:64;uniqueIndex" json:"order_no"`
	UserID       int64      `gorm:"index" json:"user_id"`
	Status       int32      `gorm:"default:1" json:"status"` // 1=待支付 2=已支付 3=已发货 4=已完成 9=已取消
	TotalAmount  float64    `gorm:"type:decimal(18,2)" json:"total_amount"`
	WarehouseID  int64      `json:"warehouse_id"`
	Remark       string     `gorm:"size:512" json:"remark"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Items        []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (Order) TableName() string { return "ord_order" }

// 订单明细
type OrderItem struct {
	ID        int64   `gorm:"primaryKey" json:"id"`
	OrderID   int64   `gorm:"index" json:"order_id"`
	SkuCode   string  `gorm:"size:64" json:"sku_code"`
	Quantity  int64   `json:"quantity"`
	UnitPrice float64 `gorm:"type:decimal(18,2)" json:"unit_price"`
}

func (OrderItem) TableName() string { return "ord_order_item" }

// 订单状态日志
type OrderLog struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	OrderID    int64     `gorm:"index" json:"order_id"`
	FromStatus int32     `json:"from_status"`
	ToStatus   int32     `json:"to_status"`
	Remark     string    `gorm:"size:512" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}

func (OrderLog) TableName() string { return "ord_order_log" }
