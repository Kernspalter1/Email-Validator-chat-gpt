package main

import (
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

// ===============================
// EMBED assets
// ===============================

//go:embed assets/*
var embeddedAssets embed.FS

// ===============================
// Browser öffnen
// ===============================
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

// ===============================
// API: /api/parse
// ===============================
func handleParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	parsed := ParseEmails(string(body))

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(parsed)
}

// ===============================
// MAIN
// ===============================
func main() {
	port := "8080"

	assetsFS, err := fs.Sub(embeddedAssets, "assets")
	if err != nil {
		log.Fatal(err)
	}

	// Statische Assets
	http.Handle("/assets/", http.StripPrefix("/assets/",
		http.FileServer(http.FS(assetsFS)),
	))

	// Root → index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := assetsFS.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "text/html")
		_, _ = io.Copy(w, f)
	})

	// API
	http.HandleFunc("/api/parse", handleParse)

	url := "http://127.0.0.1:" + port
	log.Println("Server running on", url)

	go openBrowser(url)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
