package main

import (
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	_ = exec.Command(cmd, args...).Start()
}

func main() {
	port := "8080"
	url := "http://127.0.0.1:" + port

	// Root → GUI
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "assets/index.html")
	})

	// Static assets
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Server starten
	go func() {
		log.Println("Server running on", url)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Browser automatisch öffnen
	time.Sleep(300 * time.Millisecond)
	openBrowser(url)

	select {}
}
