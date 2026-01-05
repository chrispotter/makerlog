package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSON encodes the given data as JSON and writes it to the response writer
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
