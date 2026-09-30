package event

const (
	EventTypeTicketCreated = "crm.TicketCreated"
	EventTypeTicketResolved = "crm.TicketResolved"
)

type Publisher interface {
	Publish(topic, key string, event interface{}) error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(topic, key string, event interface{}) error { return nil }
