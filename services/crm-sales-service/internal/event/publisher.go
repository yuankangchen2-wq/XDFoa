package event

const (
	EventTypeSalesOrderCreated = "crm.SalesOrderCreated"
	EventTypeContractCreated   = "crm.ContractCreated"
	EventTypePaymentReceived   = "crm.PaymentReceived"
)

type PaymentReceived struct {
	TenantID   string  `json:"tenant_id"`
	ContractID int64   `json:"contract_id"`
	Amount     float64 `json:"amount"`
}

type Publisher interface {
	Publish(topic, key string, event interface{}) error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(topic, key string, event interface{}) error { return nil }
