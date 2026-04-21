# Go Rate Limiter

A simple HTTP rate limiter service built with Go and Redis.

## Steps to Run the Project

### Prerequisites

- Go 1.21+
- Docker and Docker Compose

### Quick Start

```bash
# Start Redis and run the application
make run
```

### Manual Setup

```bash
# Start Redis
make redis-start

# Run the application
go run main.go

# Or build and run
make build
./rate-limiter
```

### Configuration

| Environment Variable | Default         | Description       |
|---------------------|-----------------|-------------------|
| `REDIS_ADDR`        | `localhost:6379`| Redis server address |

### API Endpoints

**POST /request** - Submit a rate-limited request
```bash
curl -X POST http://localhost:8080/request \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user123", "payload": {"data": "test"}}'
```

**GET /stats** - Get request counts per user
```bash
curl http://localhost:8080/stats
```

### Makefile Commands

| Command          | Description                    |
|-----------------|--------------------------------|
| `make run`      | Start Redis + run application  |
| `make build`    | Build the binary               |
| `make redis-start` | Start Redis container       |
| `make redis-stop`  | Stop Redis container        |
| `make clean`    | Remove build artifacts         |

## Design Decisions

### Fixed Window Rate Limiting

Chose a **fixed window** algorithm using Redis `INCR` and `EXPIRE`:

- **Simplicity**: Single atomic `INCR` operation per request
- **Performance**: O(1) time complexity, minimal Redis overhead
- **Memory efficient**: One key per user with automatic TTL cleanup

The rate limit is set to **5 requests per 60-second window**.

### Redis as Backend

- **Atomic operations**: `INCR` is atomic, preventing race conditions in concurrent environments
- **Built-in TTL**: Keys automatically expire, no manual cleanup needed
- **Horizontal scaling**: Multiple app instances can share the same Redis for distributed rate limiting

### Project Structure

```
├── main.go          # Application entry point
├── db/              # Redis store layer
├── service/         # Rate limiting business logic
├── router/          # HTTP handlers
├── docker-compose.yml
└── Makefile
```

Separated concerns into layers:
- **db**: Data access (Redis operations)
- **service**: Business logic (rate limit algorithm)
- **router**: HTTP interface

## What I Would Improve With More Time

### Sliding Window Algorithm

Replace fixed window with sliding window log or sliding window counter for smoother rate limiting at window boundaries. Fixed window can allow up to 2x burst at window edges.

### Configuration

- Make `MaxRequests` and `WindowSize` configurable via environment variables
- Add per-user or per-endpoint rate limit configuration
- Support different rate limit tiers

### Observability

- Add Prometheus metrics (request counts, latency, rate limit hits)
- Structured logging with levels
- Health check endpoint

### Testing

- Unit tests for service and db layers
- Integration tests with Redis
- Load testing to verify rate limiting under high concurrency

### Production Readiness

- Graceful shutdown handling
- Redis connection pooling and retry logic
- Circuit breaker for Redis failures
- Rate limit headers in responses (`X-RateLimit-Remaining`, `X-RateLimit-Reset`)

### Security

- API authentication/authorization
- Input validation and sanitization
- Rate limiting by IP address as fallback
