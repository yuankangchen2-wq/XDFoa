package event

const EventTypeAudit = "audit.LogCreated"

type Publisher interface {
	Publish(topic, key string, event interface{}) error
}

type NoopPublisher struct{}

func (NoopPublisher) Publish(topic, key string, event interface{}) error { return nil }
