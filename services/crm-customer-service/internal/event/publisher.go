package event

const (
	EventTypeCustomerCreated = "crm.CustomerCreated"
	EventTypeCustomerUpdated = "crm.CustomerUpdated"
	EventTypeFollowUpCreated = "crm.FollowUpCreated"
)

type CustomerCreated struct {
	TenantID string `json:"tenant_id"`
	CustomerID int64 `json:"customer_id"`
	Name     string `json:"name"`
}

type FollowUpCreated struct {
	TenantID   string `json:"tenant_id"`
	CustomerID int64  `json:"customer_id"`
	Type       string `json:"type"`
}

type Publisher interface {
	Publish(topic, key string, event interface{}) error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(topic, key string, event interface{}) error { return nil }
