package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func parseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	results := ParseEmails(string(body))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
