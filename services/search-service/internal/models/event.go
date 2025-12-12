package models

type EventType string

const (
	EventProductCreated EventType = "product.created"
	EventProductUpdated EventType = "product.updated"
	EventProductDeleted EventType = "product.deleted"
	EventCategoryCreated EventType = "category.created"
	EventCategoryUpdated EventType = "category.updated"
	EventCategoryDeleted EventType = "category.deleted"
)

type Event struct {
	Type      EventType `json:"type"`
	EntityID  int64     `json:"entity_id"`
	Timestamp int64     `json:"timestamp"`
}
