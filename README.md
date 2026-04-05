# TradeX

> A high-throughput, distributed paper-trading engine built in Go. Designed with an event-driven microservices architecture, TradeX ingests live market ticks from external exchanges, buffers the data stream through Kafka to prevent bottlenecking, and executes trades with millisecond latency using gRPC and Go's concurrency model.

## Tech Stack

| Component | Technology |
| --- | --- |
| **Language** | Go (Golang) |
| **API Framework** | Gin (REST API routing) |
| **Real-Time Communication** | WebSockets (`gorilla/websocket`) |
| **Internal Microservice Communication** | gRPC & Protocol Buffers (Protobuf) |
| **Event Streaming & Message Broker** | Apache Kafka |
| **Database** | MongoDB (for users, wallets, and unstructured trade receipts) |

---

## Microservices Architecture

TradeX is structured as a monorepo containing four independent microservices that communicate via gRPC and Kafka.

### 1. Auth & Wallet Service (`pkg/auth`)

> The gatekeeper and virtual bank of the platform.

- **Responsibilities:** Handles user registration, login, and JWT-based authentication via a REST API. It also manages the users' "paper trading" virtual wallet balances, storing them in MongoDB.
- **Internal Role:** Runs an internal gRPC server so the execution engine can instantly verify and deduct a user's balance in milliseconds before allowing a trade.

### 2. Market Data Service (`pkg/market_data`)

> The high-frequency data ingestion firehose.

- **Responsibilities:** Opens a persistent WebSocket connection to the Binance public API to stream live, real-world cryptocurrency price ticks.
- **Internal Role:** Acts purely as a Kafka Producer. It uses goroutines to process the massive stream of incoming raw JSON price ticks and fires them into a Kafka topic asynchronously, feeding the rest of the system without dropping data.

### 3. Order Matcher Service (`pkg/order_matcher`)

> The core trade execution engine.

- **Responsibilities:** Receives incoming "Buy" or "Sell" requests via gRPC. It constantly consumes live prices from Kafka to know the exact asset price at the precise millisecond of the trade request.
- **Internal Role:** Uses gRPC to ping the Auth Service to verify funds. Upon successful validation, it executes the trade, saves the transaction receipt to MongoDB, and publishes an "order executed" event back into Kafka.

### 4. Notifier Service (`pkg/notifier`)

> The real-time client communication hub.

- **Responsibilities:** Pushes live state updates directly to the user's browser UI so charts and balances update instantly without page refreshes.
- **Internal Role:** Acts as a Kafka Consumer. It listens for both live price ticks and executed trade events from the Kafka cluster, broadcasting those updates to connected frontend clients via WebSockets.

## The "Buy" Execution Flow

When a user clicks "Buy Bitcoin" on the frontend, the system handles the request entirely asynchronously:

| Step | Flow |
| --- | --- |
| **1. The Price Tick** | `market_data` streams the live Bitcoin price from Binance and pushes it into the `market.prices` Kafka topic. |
| **2. The State Update** | `notifier` consumes this price from Kafka and pushes it to the user's browser via WebSocket. The user sees the live price and clicks "Buy". |
| **3. The Trade Request** | The frontend sends a gRPC "Buy" request directly to the `order_matcher` service. |
| **4. Fund Verification** | `order_matcher` checks the live price it holds in memory (from Kafka), then fires a lightning-fast internal gRPC call to the `auth` service: *"Deduct $65,000 from this user's wallet."* |
| **5. Execution** | `auth` successfully deducts the virtual funds in MongoDB. `order_matcher` logs the executed trade in MongoDB and fires a `trade.success` event back into Kafka. |
| **6. The Notification** | `notifier` reads the `trade.success` event from Kafka and instantly pushes a "Trade Completed" notification to the user's browser via WebSocket. |