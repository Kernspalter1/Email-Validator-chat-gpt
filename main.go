package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// =======================
// Embedded GUI Assets
// =======================

//go:embed assets/*
var embeddedAssets embed.FS

// =======================
// Plausibility Enum
// =======================

type Plausibility string

const (
	ACCEPTED  Plausibility = "ACCEPTED"
	UNCERTAIN Plausibility = "UNCERTAIN"
	REJECTED  Plausibility = "REJECTED"
)

// =======================
// Models
// =======================

type ParsedEmail struct {
	Email       string       `json:"email"`
	Plausibility Plausibility `json:"plausibility"`
	Duplicate   bool         `json:"duplicate"`
}

// =======================
// Parser (Backend-Wahrheit)
// =======================

func ParseEmails(input string) []ParsedEmail {
	lines := strings.Split(input, "\n")
	seen := make(map[string]bool)
	var result []ParsedEmail

	for _, line := range lines {
		email := strings.TrimSpace(strings.ToLower(line))
		if email == "" {
			continue
		}

		duplicate := seen[email]
		if !duplicate {
			seen[email] = true
		}

		pl := ACCEPTED
		if !strings.Contains(email, "@") {
			pl = REJECTED
		}

		result = append(result, ParsedEmail{
			Email:       email,
			Plausibility: pl,
			Duplicate:   duplicate,
		})
	}

	return result
}

// =======================
// HTTP Handlers
// =======================

func parseHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Input string `json:"input"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	parsed := ParseEmails(payload.Input)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsed)
}

// =======================
// Main
// =======================

func main() {
	mux := http.NewServeMux()

	// API
	mux.HandleFunc("/api/parse", parseHandler)

	// GUI
	fileServer := http.FileServer(http.FS(embeddedAssets))
	mux.Handle("/assets/", http.StripPrefix("/", fileServer))

	// ROOT → GUI
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/", http.StatusFound)
	})

	log.Println("Server running on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
