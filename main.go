package main

import (
	"embed"
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

func main() {
	port := "8080"

	// Sub-FS nur für assets
	assetsFS, err := fs.Sub(embeddedAssets, "assets")
	if err != nil {
		log.Fatal(err)
	}

	// 1️⃣ Statische Dateien (JS, CSS, etc.)
	fileServer := http.FileServer(http.FS(assetsFS))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	// 2️⃣ Root → index.html aus embed
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := assetsFS.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		w.Header().Set("Content-Type", "text/html")
		_, _ = w.ReadFrom(f)
	})

	// 3️⃣ API
	http.HandleFunc("/api/parse", handleParse)

	url := "http://127.0.0.1:" + port
	log.Println("Server running on", url)

	// 4️⃣ Browser automatisch öffnen
	go openBrowser(url)

	// 5️⃣ Server starten
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
