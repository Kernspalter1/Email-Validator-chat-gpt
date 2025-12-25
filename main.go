package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	host = "127.0.0.1"
	port = 8080
)

func main() {
	fmt.Println("=== LocalEmailHealthChecker START ===")

	// 1. HTTP-Server konfigurieren
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
	}

	// 2. Server im Hintergrund starten
	go func() {
		fmt.Println("Starting local server on http://" + host + ":" + fmt.Sprint(port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 3. Kurz warten, damit Server sicher läuft
	time.Sleep(500 * time.Millisecond)

	// 4. Browser öffnen
	openBrowser(fmt.Sprintf("http://%s:%d", host, port))

	// 5. Blockieren, damit EXE nicht sofort endet
	fmt.Println("Server running. Press CTRL+C to exit.")
	select {}
}

// -----------------------------

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Email Health Checker</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			margin: 40px;
		}
		textarea {
			width: 100%;
			height: 200px;
		}
		button {
			margin-top: 10px;
			padding: 10px 20px;
		}
	</style>
</head>
<body>
	<h1>Email Health Checker</h1>
	<p>GUI bootstrap successful.</p>
	<p>Backend is running locally.</p>

	<textarea placeholder="Paste emails here (one per line)"></textarea><br>
	<button disabled>Validate (coming next)</button>
</body>
</html>
`)
}

// -----------------------------

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		fmt.Println("Unsupported OS. Please open manually:", url)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to open browser:", err)
		fmt.Println("Please open manually:", url)
	}
}
