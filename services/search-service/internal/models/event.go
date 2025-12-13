package models

type EventType string

const (
	EventProductCreated  EventType = "product.created"
	EventProductUpdated  EventType = "product.updated"
	EventProductDeleted  EventType = "product.deleted"
	EventCategoryCreated EventType = "category.created"
	EventCategoryUpdated EventType = "category.updated"
	EventCategoryDeleted EventType = "category.deleted"
	EventBidPlaced       EventType = "bid.placed"       // New event for bid updates
	EventBidderUpdated   EventType = "bidder.updated"   // New event for bidder info updates
)

type Event struct {
	Type      EventType              `json:"type"`
	EntityID  int64                  `json:"entity_id"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"` // Additional event data
}

// BidEventData represents data for bid events
type BidEventData struct {
	ProductID       int64       `json:"product_id"`
	CurrentPrice    float64     `json:"current_price"`
	CurrentBidCount int         `json:"current_bid_count"`
	BidderInfo      *BidderInfo `json:"bidder_info,omitempty"`
}
