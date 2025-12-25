package main

import (
	"log"
	"net/http"
	"os"
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
	default: // linux, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	_ = exec.Command(cmd, args...).Start()
}

func main() {
	port := "8080"
	url := "http://127.0.0.1:" + port

	// 1️⃣ Root → index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// nur exakt "/" hier abfangen
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "assets/index.html")
	})

	// 2️⃣ Statische Assets (/assets/*)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// 3️⃣ API (bereits vorhanden)
	http.HandleFunc("/api/parse", handleParse)

	// 4️⃣ Server starten
	go func() {
		log.Println("Server running on", url)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// 5️⃣ Kurz warten, dann Browser öffnen
	time.Sleep(300 * time.Millisecond)
	openBrowser(url)

	// 6️⃣ Blockieren (EXE soll laufen)
	select {}
}
