package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

// =====================
// EMBED ASSETS
// =====================

//go:embed assets/*
var embeddedAssets embed.FS

// =====================
// MAIN
// =====================

func main() {
	mux := http.NewServeMux()

	// 1️⃣ Root → index.html aus assets/
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		data, err := embeddedAssets.ReadFile("assets/index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	// 2️⃣ Statische Dateien (JS, CSS)
	fileServer := http.FileServer(http.FS(embeddedAssets))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	// 3️⃣ API: /parse  (nutzt parser.go)
	mux.HandleFunc("/parse", handleParse)

	addr := "127.0.0.1:8080"
	log.Println("Server running on http://" + addr)

	// Browser automatisch öffnen
	go openBrowser("http://" + addr)

	// Server starten
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

// =====================
// API HANDLER
// =====================

func handleParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Input string `json:"input"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	results := ParseEmails(req.Input)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// =====================
// BROWSER AUTO OPEN
// =====================

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	_ = cmd.Start()
}
