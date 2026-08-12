package notify

// --- Event Definitions ---

type EventType string

const (
	EventStockUpdate EventType = "STOCK_UPDATE"
)

type Event struct {
	Type    EventType
	Message string
	Data    map[string]interface{}
}
