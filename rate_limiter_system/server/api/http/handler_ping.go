package http

import (
	"encoding/json"
	"net/http"

	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/structs"
)

func HandlerPing(w http.ResponseWriter, r *http.Request) {

	message := structs.Message{
		Status: "Successful",
		Body:   "Hi! You've reached the API.",
	}

	// Directly encode and write the message to the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
