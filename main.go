package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

//go:embed assets/*
var assets embed.FS

const addr = "127.0.0.1:8080"

type Request struct {
	Emails           []string `json:"emails"`
	IncludeDuplicates bool     `json:"includeDuplicates"`
}

type Result struct {
	Email  string `json:"email"`
	Status string `json:"status"`
	Reason string `json:"reason"`
	Rank   int    `json:"-"`
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.FS(assets)))
	mux.HandleFunc("/validate", validateHandler)

	go func() {
		log.Println("Server running at http://" + addr)
		log.Fatal(http.ListenAndServe(addr, mux))
	}()

	time.Sleep(400 * time.Millisecond)
	openBrowser("http://" + addr)
	select {}
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	json.NewDecoder(r.Body).Decode(&req)

	seen := map[string]bool{}
	results := []Result{}

	for _, raw := range req.Emails {
		email := strings.TrimSpace(strings.ToLower(raw))
		if email == "" {
			continue
		}

		if !req.IncludeDuplicates {
			if seen[email] {
				continue
			}
			seen[email] = true
		}

		res := dummyValidate(email)
		results = append(results, res)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func dummyValidate(email string) Result {
	if strings.HasSuffix(email, "@gmail.com") {
		return Result{email, "accepted", "SMTP accepted", 1}
	}
	if strings.IndexAny(email, "0123456789") >= 0 {
		return Result{email, "uncertain", "SMTP unclear", 2}
	}
	return Result{email, "rejected", "No SMTP", 3}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	cmd.Start()
}
