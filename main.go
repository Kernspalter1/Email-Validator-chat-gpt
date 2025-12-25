package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

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

	// 1. Statische Assets unter /assets/ bereitstellen
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// 2. Root "/" explizit auf assets/index.html mappen
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/index.html")
	})

	// 3. API-Endpunkte
	http.HandleFunc("/api/parse", handleParse)

	url := "http://127.0.0.1:" + port
	log.Println("Server running on", url)

	// 4. Browser automatisch öffnen
	go openBrowser(url)

	// 5. Server starten
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
