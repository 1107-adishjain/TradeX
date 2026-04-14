package models

// BinanceTradePayload is the raw Binance trade payload.
type BinanceTradePayload struct {
	EventType string `json:"e"` // Event type (e.g., "trade")
	EventTime int64  `json:"E"` // Event time (Unix timestamp)
	Symbol    string `json:"s"` // Symbol (e.g., "BTCUSDT")
	Price     string `json:"p"` // Price
	Quantity  string `json:"q"` // Quantity traded
}

// BinanceMultiplexPayload wraps trade data for multiplex streams.
type BinanceMultiplexPayload struct {
	Stream string              `json:"stream"` // The stream name, e.g., "btcusdt@trade"
	Data   BinanceTradePayload `json:"data"`   // The actual trade data nested inside
}

// TradeXPriceUpdate is our internal normalized price event.That is pushed to kafka.
type TradeXPriceUpdate struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}
