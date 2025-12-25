package main

import (
	"embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
)

//go:embed assets/*
var embeddedAssets embed.FS

func main() {
	mux := http.NewServeMux()

	// 1) Root soll automatisch GUI öffnen
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// wenn genau "/" -> redirect
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/assets/", http.StatusFound)
			return
		}
		// alles andere, was wir nicht kennen:
		http.NotFound(w, r)
	})

	// 2) assets/ statisch ausliefern
	//    /assets/ -> assets/index.html
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// directory index
		reqPath := r.URL.Path
		if reqPath == "" || strings.HasSuffix(reqPath, "/") {
			reqPath = path.Join(reqPath, "index.html")
		}

		filePath := path.Join("assets", reqPath)
		b, err := embeddedAssets.ReadFile(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		// mini content-type handling
		switch {
		case strings.HasSuffix(filePath, ".html"):
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case strings.HasSuffix(filePath, ".js"):
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case strings.HasSuffix(filePath, ".css"):
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		default:
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	})))

	// 3) API: parse + validate (minimal)
	mux.HandleFunc("/api/parse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req ParseRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		parsed := ParseEmails(req.Text, req.KeepDuplicates)
		validated := ValidateParsed(parsed)

		// Listen für TXT-Export vorbereiten:
		accepted := make([]ParsedEmail, 0)
		acceptedUncertain := make([]ParsedEmail, 0)

		seenUnique := make(map[string]bool)

		for _, it := range validated {
			// Unique-Listen nur einmal, egal ob Duplikate angezeigt werden
			if seenUnique[it.Email] {
				// skip for unique output lists
			} else {
				seenUnique[it.Email] = true

				if it.Status == StatusAccepted {
					accepted = append(accepted, it)
					acceptedUncertain = append(acceptedUncertain, it)
				} else if it.Status == StatusUncertain {
					acceptedUncertain = append(acceptedUncertain, it)
				}
			}
		}

		resp := ParseResponse{
			Items:            validated,
			AcceptedSorted:    sortByWeightDesc(accepted),
			AcceptedUncertain: sortByWeightDesc(acceptedUncertain),
		}

		// Default-UniqueSorted = accepted (so wie dein Default)
		resp.UniqueSorted = resp.AcceptedSorted

		writeJSON(w, resp)
	})

	addr := "127.0.0.1:8080"
	log.Printf("Server running on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}
