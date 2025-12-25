package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.Handle("/assets/",
		http.StripPrefix("/assets/",
			http.FileServer(http.Dir("./assets"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/", http.StatusFound)
	})

	http.HandleFunc("/api/parse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload struct {
			Input string `json:"input"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		parsed := ParseEmails(payload.Input)
		validated := ValidateEmails(parsed)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(validated)
	})

	log.Println("Server running on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
