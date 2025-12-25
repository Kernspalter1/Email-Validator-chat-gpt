package main

import "strings"

func ValidateEmails(items []ParsedEmail) []ParsedEmail {
	for i := range items {
		if items[i].Duplicate {
			items[i].Status = REJECTED
			items[i].Reason = "duplicate"
			continue
		}

		if strings.Contains(items[i].Email, "@") &&
			strings.Contains(items[i].Email, ".") {
			items[i].Status = ACCEPTED
			items[i].Reason = "basic syntax ok"
		} else {
			items[i].Status = REJECTED
			items[i].Reason = "invalid syntax"
		}
	}
	return items
}
