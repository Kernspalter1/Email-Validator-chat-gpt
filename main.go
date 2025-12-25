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

// Embed everything under assets/
//go:embed assets/*
var embedded embed.FS

const addr = "127.0.0.1:8080"

type ValidateRequest struct {
	Emails            []string `json:"emails"`
	IncludeDuplicates bool     `json:"includeDuplicates"`
}

type ValidateResult struct {
	Email  string `json:"email"`
	Status string `json:"status"` // accepted | uncertain | rejected
	Reason string `json:"reason"`
}

func main() {
	mux := http.NewServeMux()

	// Serve embedded static files
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

	// Give server a moment
	time.Sleep(400 * time.Millisecond)

	// Open browser
	openBrowser("http://" + addr)

	// Keep running
	select {}
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	var req ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	seen := make(map[string]bool)
	out := make([]ValidateResult, 0, len(req.Emails))

	for _, raw := range req.Emails {
		email := strings.TrimSpace(strings.ToLower(raw))
		if email == "" {
			continue
		}

		// Remove duplicates by default, unless checkbox is enabled
		if !req.IncludeDuplicates {
			if seen[email] {
				continue
			}
			seen[email] = true
		}

		out = append(out, dummyValidate(email))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

// Dummy logic ONLY for Step 1 (UI + TXT workflow)
// Later we replace this with real MX/SMTP checks.
func dummyValidate(email string) ValidateResult {
	if strings.HasSuffix(email, "@gmail.com") {
		return ValidateResult{Email: email, Status: "accepted", Reason: "SMTP accepted (dummy)"}
	}
	if strings.IndexAny(email, "0123456789") >= 0 {
		return ValidateResult{Email: email, Status: "uncertain", Reason: "SMTP unclear (dummy)"}
	}
	return ValidateResult{Email: email, Status: "rejected", Reason: "No SMTP (dummy)"}
}

func openBrowser(url string) {
	fmt.Println("Opening:", url)

	if runtime.GOOS == "windows" {
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		return
	}
	// Other OS not needed (Windows only), but harmless:
	_ = exec.Command("xdg-open", url).Start()
}
