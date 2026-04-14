package binance

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/adishjain1107/tradex/pkg/market-data/models"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

// StartMultiplexStream streams Binance trades and publishes updates to Kafka.
func StartMultiplexStream(kafkaBroker string, symbols []string) {
	if len(symbols) == 0 {
		log.Println("No market symbols configured. Stream not started.")
		return
	}

	// Setup the Kafka Writer (Producer)
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaBroker),
		Topic:        "market.prices",
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond, // Low latency for trading data
	}
	defer writer.Close()

	// Format the Binance Multiplex URL. E.g., converting ["BTCUSDT", "ETHUSDT"] into "btcusdt@trade/ethusdt@trade"
	var streams []string
	for _, sym := range symbols {
		streams = append(streams, strings.ToLower(sym)+"@trade")
	}
	combinedStreams := strings.Join(streams, "/")
	binanceURL := "wss://stream.binance.com:9443/stream?streams=" + combinedStreams

	log.Printf("Starting Binance Multiplex Stream: %s", binanceURL)

	const reconnectDelay = 2 * time.Second

	for {
		conn, _, err := websocket.DefaultDialer.Dial(binanceURL, nil)
		if err != nil {
			log.Printf("Failed to dial Binance: %v", err)
			time.Sleep(reconnectDelay)
			continue
		}

		log.Println("Connected to Binance. Ingesting live trades...")

		for {
			// Wait for Binance to push a message down the tunnel
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Binance connection dropped: %v", err)
				_ = conn.Close()
				break
			}

			// Unmarshal the outer Multiplex wrapper
			var payload models.BinanceMultiplexPayload
			if err := json.Unmarshal(message, &payload); err != nil {
				log.Printf("Failed to parse multiplex payload: %v", err)
				continue
			}

			// Convert string price to float64 for TradeX internal math
			price, err := strconv.ParseFloat(payload.Data.Price, 64)
			if err != nil {
				log.Printf("Failed to parse price %q: %v", payload.Data.Price, err)
				continue
			}

			// Map to our pristine TradeX standard model
			tradeUpdate := models.TradeXPriceUpdate{
				Symbol:    payload.Data.Symbol,
				Price:     price,
				Timestamp: payload.Data.EventTime,
			}

			// Marshal it back to JSON for Kafka
			kafkaPayload, err := json.Marshal(tradeUpdate)
			if err != nil {
				log.Printf("Failed to marshal TradeX update: %v", err)
				continue
			}

			// Push to Kafka, using the Symbol as the Key to guarantee chronological order
			err = writer.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte(tradeUpdate.Symbol),
					Value: kafkaPayload,
				},
			)

			if err != nil {
				log.Printf("Failed to write to Kafka: %v", err)
			}
		}

		log.Printf("Reconnecting to Binance in %s...", reconnectDelay)
		time.Sleep(reconnectDelay)
	}
}
