package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-rate-limiter/db"
	"github.com/go-rate-limiter/router"
	"github.com/go-rate-limiter/service"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	store, err := db.New(redisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer store.Close()

	svc := service.New(store)
	r := router.Setup(svc)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
