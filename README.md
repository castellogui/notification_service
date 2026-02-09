# Notification Service

A high-throughput, event-driven notification service built with Go. It decouples notification ingestion from processing and delivery through Kafka, persists data in ScyllaDB, and is designed to scale horizontally.

## Architecture

The system is split into two independent services:

```
                                  ┌──────────────┐
  HTTP Request ──► API Service ──►│    Kafka      │──► Pusher Service ──► ScyllaDB
                   (gin)          │  (broker)     │    (consumer)         (persistence)
                                  └──────────────┘
```

### API Service

Exposes a REST API for notification management. When a notification is created, the API publishes an event to Kafka and returns immediately — the caller never waits for downstream processing.

Endpoints:

| Method  | Route                                | Description              |
|---------|--------------------------------------|--------------------------|
| POST    | `/v1/notifications/`                 | Create a notification    |
| GET     | `/v1/notifications/:user_id/:id`     | Get a notification       |
| PATCH   | `/v1/notifications/:user_id/:id`     | Partial update           |
| DELETE  | `/v1/notifications/:id`              | Delete a notification    |

### Pusher Service

A Kafka consumer that processes notification events. For each message it:

1. Deserializes the envelope
2. Decodes the kind-specific payload via the registry
3. Builds a user-facing view model (title, body, deep link, etc.)
4. Persists the notification to ScyllaDB

The pusher is stateless and can be scaled horizontally by adding more consumers to the group.

## Project Structure

```
cmd/
  api/main.go              # API entrypoint
  pusher/main.go           # Kafka consumer entrypoint

internal/
  api/
    router.go              # Gin router setup and dependency wiring
    handlers/notification/  # HTTP handlers (create, get, update, delete)
    routes/                 # Route registration

  pusher/
    handler.go             # Kafka message processing
    domain/                # Core models (Envelope, ViewModel, NotificationDB)
    interfaces/            # Writer interface for DB abstraction
    registry/              # Decoder/builder registry for notification kinds
    kinds/                 # Kind implementations (e.g. status)
    setup/                 # Registry initialization

  infra/
    scylla_db.go                  # ScyllaDB session management
    scylla_notification_writer.go # Writer interface implementation
```

## Registry Pattern

The notification system supports multiple notification kinds, each with versioned payloads. Adding a new kind requires no changes to the core pipeline:

1. Define a payload struct in `internal/pusher/kinds/`
2. Implement a decoder (JSON to struct) and a builder (struct to ViewModel)
3. Register it in `internal/pusher/setup/registry.go`

The registry maps `(Kind, Version)` to its decoder and builder, making the pipeline fully polymorphic. The pusher doesn't need to know the specifics of any notification type — it just calls `Decode` and `Build`.

## Tech Stack

### Go

Compiles to a single binary with no runtime dependencies. Low memory footprint and goroutine-based concurrency make it a natural fit for a service that needs to handle high message throughput with minimal resource consumption.

### Kafka

Decouples the API from downstream processing entirely. The API writes to a topic and returns — the pusher consumes at its own pace. This gives us:

- **Back-pressure handling**: if the pusher falls behind, messages queue in Kafka rather than failing at the API.
- **Horizontal scaling**: adding more consumers to the group distributes partitions automatically.
- **Durability**: messages are persisted to disk and replicated across brokers.

### ScyllaDB

A Cassandra-compatible database written in C++, designed for low-latency at high throughput. Chosen because:

- **Partition-based access**: notifications are partitioned by `user_id`, which matches the primary access pattern (fetch all notifications for a user).
- **Linear scalability**: adding nodes increases throughput proportionally with no coordination overhead.
- **Tunable consistency**: allows trading consistency for latency per-query depending on the use case.

### Gin

A lightweight HTTP framework with minimal overhead. Provides routing, middleware, and JSON binding without pulling in a large dependency tree.

## Running

### Infrastructure

```bash
docker-compose up -d
```

This starts Kafka, Zookeeper, and the Kafka UI tools.

### Services

```bash
# API (port 8080)
go run cmd/api/main.go

# Pusher (Kafka consumer)
go run cmd/pusher/main.go
```

ScyllaDB must be running separately on `localhost:9042` with keyspace `notification_service`.
