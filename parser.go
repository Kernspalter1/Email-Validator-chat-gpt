package main

import "strings"

func ParseEmails(input string) []ParsedEmail {
	lines := strings.Split(input, "\n")

	seen := make(map[string]bool)
	results := make([]ParsedEmail, 0)

	for _, line := range lines {
		email := strings.TrimSpace(strings.ToLower(line))
		if email == "" {
			continue
		}

		pe := ParsedEmail{Email: email}

		if seen[email] {
			pe.Duplicate = true
			pe.Status = REJECTED
			pe.Reason = "duplicate"
		} else {
			seen[email] = true
			pe.Duplicate = false
			pe.Status = UNCERTAIN
			pe.Reason = "not yet validated"
		}

		results = append(results, pe)
	}

	return results
}
