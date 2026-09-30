package event

const (
	EventTypePurchaseOrderCreated = "purchase.OrderCreated"
	EventTypeGoodsReceived        = "purchase.GoodsReceived"
	EventTypePurchaseSettled      = "purchase.Settled"
)

// PurchaseOrderCreated 采购订单创建事件
type PurchaseOrderCreated struct {
	TenantID     string  `json:"tenant_id"`
	OrderNo      string  `json:"order_no"`
	SupplierCode string  `json:"supplier_code"`
	WarehouseID  int64   `json:"warehouse_id"`
	TotalAmount  float64 `json:"total_amount"`
}

// GoodsReceived 到货事件
type GoodsReceived struct {
	TenantID    string `json:"tenant_id"`
	ReceiptNo   string `json:"receipt_no"`
	OrderID     int64  `json:"order_id"`
	WarehouseID int64  `json:"warehouse_id"`
	SkuCode     string `json:"sku_code"`
	Quantity    int64  `json:"quantity"`
}

type Publisher interface {
	Publish(topic, key string, event interface{}) error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(topic, key string, event interface{}) error { return nil }
