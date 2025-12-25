package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

//go:embed assets/*
var embedded embed.FS

const addr = "127.0.0.1:8080"

func main() {
	mux := http.NewServeMux()

	// Serve embedded frontend
	mux.Handle("/", http.FileServer(http.FS(embedded)))

	// Backend endpoint
	mux.HandleFunc("/validate", handleValidate)

	// Start server
	go func() {
		log.Println("Server running at http://" + addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatal(err)
		}
	}()

	// Give server a moment to start
	time.Sleep(400 * time.Millisecond)

	// Open browser automatically
	openBrowser("http://" + addr)

	// Keep process alive
	select {}
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	entries := ParseEmailsFromText(req.Text)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entries)
}

func openBrowser(url string) {
	if runtime.GOOS == "windows" {
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	}
}
