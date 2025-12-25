package main

import (
	"strings"
)

// ParseEmails nimmt rohen Text (eine Mail pro Zeile)
// und gibt eine einfache Struktur zurück.
// DAS ist nur der Parser – keine Validation hier!
func ParseEmails(input string) []ParsedEmail {
	lines := strings.Split(input, "\n")

	seen := make(map[string]bool)
	results := []ParsedEmail{}

	for _, line := range lines {
		email := strings.TrimSpace(line)
		if email == "" {
			continue
		}

		duplicate := false
		if seen[email] {
			duplicate = true
		}
		seen[email] = true

		results = append(results, ParsedEmail{
			Email:     email,
			Duplicate: duplicate,
		})
	}

	return results
}
