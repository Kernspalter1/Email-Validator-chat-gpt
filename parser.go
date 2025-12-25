package main

import (
	"strings"
)

type EmailEntry struct {
	Email          string
	IsDuplicate    bool
	FirstSeenIndex int
	OriginalLine   int
}

// ParseEmailsFromText
// - takes raw TXT input (one email per line)
// - normalizes emails (trim + lowercase)
// - detects duplicates
// - DOES NOT delete anything
func ParseEmailsFromText(input string) []EmailEntry {
	lines := strings.Split(input, "\n")

	seen := make(map[string]int) // email -> first index
	result := make([]EmailEntry, 0, len(lines))

	for i, raw := range lines {
		email := strings.TrimSpace(raw)
		if email == "" {
			continue
		}

		email = strings.ToLower(email)

		entry := EmailEntry{
			Email:        email,
			OriginalLine: i,
		}

		if firstIndex, exists := seen[email]; exists {
			entry.IsDuplicate = true
			entry.FirstSeenIndex = firstIndex
		} else {
			seen[email] = len(result)
			entry.IsDuplicate = false
			entry.FirstSeenIndex = len(result)
		}

		result = append(result, entry)
	}

	return result
}
