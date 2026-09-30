package event

const (
	EventTypeOpportunityCreated = "crm.OpportunityCreated"
	EventTypeOpportunityWon     = "crm.OpportunityWon"
	EventTypeOpportunityLost    = "crm.OpportunityLost"
)

type OpportunityWon struct {
	TenantID     string  `json:"tenant_id"`
	OpportunityID int64  `json:"opportunity_id"`
	CustomerID   int64   `json:"customer_id"`
	Amount       float64 `json:"amount"`
}

type Publisher interface {
	Publish(topic, key string, event interface{}) error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(topic, key string, event interface{}) error { return nil }
