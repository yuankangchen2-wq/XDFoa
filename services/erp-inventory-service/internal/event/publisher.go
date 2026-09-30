package event

// 事件类型
const (
	EventTypeStockChanged    = "inventory.StockChanged"
	EventTypeStockShortage   = "inventory.StockShortage"
	EventTypeStockTransferred = "inventory.StockTransferred"
)

// StockChanged 库存变动事件
type StockChanged struct {
	TenantID     string `json:"tenant_id"`
	SkuCode      string `json:"sku_code"`
	WarehouseID  int64  `json:"warehouse_id"`
	BeforeQty    int64  `json:"before_qty"`
	AfterQty     int64  `json:"after_qty"`
	ChangeType   string `json:"change_type"`
	ReferenceID  string `json:"reference_id"`
}

// StockShortage 库存不足事件
type StockShortage struct {
	TenantID     string `json:"tenant_id"`
	SkuCode      string `json:"sku_code"`
	WarehouseID  int64  `json:"warehouse_id"`
	RequiredQty  int64  `json:"required_qty"`
	AvailableQty int64  `json:"available_qty"`
	ReferenceID  string `json:"reference_id"`
}

// StockTransferred 调拨事件
type StockTransferred struct {
	TenantID        string `json:"tenant_id"`
	TransferNo      string `json:"transfer_no"`
	FromWarehouseID int64  `json:"from_warehouse_id"`
	ToWarehouseID   int64  `json:"to_warehouse_id"`
	SkuCode         string `json:"sku_code"`
	Quantity        int64  `json:"quantity"`
}

type Publisher interface {
	Publish(topic, key string, event interface{}) error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(topic, key string, event interface{}) error { return nil }
