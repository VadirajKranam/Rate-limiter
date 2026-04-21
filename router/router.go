package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-rate-limiter/service"
)

type Request struct {
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload"`
}

func Setup(svc *service.RateLimiter) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
			return
		}

		if !svc.Allow(req.UserID) {
			writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "Rate limit exceeded"})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"message": "Request accepted"})
	})

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, svc.GetStats())
	})

	return mux
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
