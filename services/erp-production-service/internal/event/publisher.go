package event

const (
	EventTypeWorkOrderCreated   = "production.WorkOrderCreated"
	EventTypeMaterialIssued     = "production.MaterialIssued"
	EventTypeProductionCompleted = "production.ProductionCompleted"
)

type WorkOrderCreated struct {
	TenantID   string `json:"tenant_id"`
	OrderNo    string `json:"order_no"`
	ProductSku string `json:"product_sku"`
	Quantity   int64  `json:"quantity"`
}

type MaterialIssued struct {
	TenantID   string `json:"tenant_id"`
	IssueNo    string `json:"issue_no"`
	WorkOrderID int64 `json:"work_order_id"`
}

type ProductionCompleted struct {
	TenantID   string `json:"tenant_id"`
	OrderNo    string `json:"order_no"`
	ProductSku string `json:"product_sku"`
	Quantity   int64  `json:"quantity"`
}

type Publisher interface {
	Publish(topic, key string, event interface{}) error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(topic, key string, event interface{}) error { return nil }
