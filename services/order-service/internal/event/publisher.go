package event

const (
	EventTypeOrderCreated   = "order.OrderCreated"
	EventTypeOrderPaid      = "order.OrderPaid"
	EventTypeOrderShipped   = "order.OrderShipped"
	EventTypeOrderCancelled = "order.OrderCancelled"
	EventTypeOrderCompleted = "order.OrderCompleted"
)

type OrderPaid struct {
	TenantID string `json:"tenant_id"`
	OrderNo  string `json:"order_no"`
	OrderID  int64  `json:"order_id"`
	UserID   int64  `json:"user_id"`
}

type OrderCancelled struct {
	TenantID string `json:"tenant_id"`
	OrderNo  string `json:"order_no"`
	OrderID  int64  `json:"order_id"`
}

type Publisher interface {
	Publish(topic, key string, event interface{}) error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(topic, key string, event interface{}) error { return nil }
