package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// 事件类型常量
const (
	EventTypeEntitySynced = "mdm.EntitySynced"
	EventVersionV1        = "v1"
)

// 事件统一结构
type Event struct {
	EventID     string      `json:"event_id"`
	EventType   string      `json:"event_type"`
	EventVersion string     `json:"event_version"`
	Source      string      `json:"source"`
	TenantID    string      `json:"tenant_id"`
	OccurredAt  int64       `json:"occurred_at"`
	TraceID     string      `json:"trace_id"`
	EntityType  string      `json:"entity_type"` // customer/product/supplier/organization
	Action      string      `json:"action"`      // created/updated/deleted
	EntityID    string      `json:"entity_id"`
	Payload     interface{} `json:"payload"`
}

// Publisher 事件发布接口（可替换为 Kafka 实现）
type Publisher interface {
	Publish(ctx context.Context, event *Event) error
}

// NoopPublisher 默认空实现（无 Kafka 时使用）
type NoopPublisher struct{}

func (NoopPublisher) Publish(ctx context.Context, event *Event) error {
	return nil
}

// NewEntitySyncedEvent 构造 MDM 实体同步事件
func NewEntitySyncedEvent(tenantID, entityType, action string, entityID string, payload interface{}) *Event {
	return &Event{
		EventID:      uuid.NewString(),
		EventType:    EventTypeEntitySynced,
		EventVersion: EventVersionV1,
		Source:       "mdm-service",
		TenantID:     tenantID,
		OccurredAt:   time.Now().UnixMilli(),
		EntityType:   entityType,
		Action:       action,
		EntityID:     entityID,
		Payload:      payload,
	}
}

// Marshal 序列化为 JSON
func (e *Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// String 用于日志
func (e *Event) String() string {
	return fmt.Sprintf("%s/%s/%s/%s", e.EventType, e.EntityType, e.Action, e.EntityID)
}
