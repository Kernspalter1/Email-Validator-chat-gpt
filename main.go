package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// ----------------------------
// Main
// ----------------------------
func main() {
	mux := http.NewServeMux()

	// 1️⃣ Statische Assets (CSS, JS)
	// Zugriff z.B. auf /assets/app.js
	fileServer := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	// 2️⃣ Root -> index.html
	// Damit http://127.0.0.1:8080 direkt die GUI lädt
	mux.HandleFunc("/", serveIndex)

	// 3️⃣ API Endpoint
	mux.HandleFunc("/validate", validateHandler)

	log.Println("Server running on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// ----------------------------
// Handlers
// ----------------------------

// serveIndex liefert IMMER die GUI aus
func serveIndex(w http.ResponseWriter, r *http.Request) {
	// Optional: nur GET erlauben
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "./assets/index.html")
}

// validateHandler verarbeitet die Eingabe aus dem Frontend
func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Erwartet JSON vom Frontend
	// z.B. { "input": "a@gmail.com\nb@gmail.com" }
	var req struct {
		Input string `json:"input"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 👉 HIER kommt später dein Parser rein
	// z.B.:
	// results := ParseEmails(req.Input)
	// validated := Validate(results)

	// Platzhalter-Antwort (damit Frontend funktioniert)
	resp := map[string]interface{}{
		"status":  "ok",
		"message": "backend wired successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
