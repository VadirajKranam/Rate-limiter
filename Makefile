.PHONY: run build test clean redis redis-stop redis-logs

# Go parameters
BINARY_NAME=rate-limiter
REDIS_ADDR?=localhost:6379

# Build the application
build:
	go build -o $(BINARY_NAME) .

# Run the application
run: redis-start
	REDIS_ADDR=$(REDIS_ADDR) go run main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Start Redis using Docker Compose
redis-start:
	docker compose up -d redis
	@echo "Waiting for Redis to be ready..."
	@sleep 2

# Stop Redis
redis-stop:
	docker compose down

# View Redis logs
redis-logs:
	docker compose logs -f redis

# Start everything (Redis + app)
up: redis-start run

# Stop everything
down: redis-stop
