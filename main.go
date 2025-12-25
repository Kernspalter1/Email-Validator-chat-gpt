package main

import (
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

// Öffnet Standard-Browser
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

	// 1. Statische Assets bereitstellen
	fileServer := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	// 2. Root "/" -> index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./assets/index.html")
	})

	// 3. API-Endpunkt
	http.HandleFunc("/api/parse", parseHandler)

	url := "http://127.0.0.1:" + port
	log.Println("Server running on", url)

	// 4. Browser automatisch öffnen
	go openBrowser(url)

	// 5. Server starten
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
