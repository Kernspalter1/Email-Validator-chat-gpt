package main

import (
	"regexp"
)

// einfacher Syntax-Check zusätzlich zum Parser
var simpleEmailSyntax = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

// ValidateParsed nimmt Parser-Ergebnis und setzt Status.
// Aktuell:
// - gültige Syntax => UNCERTAIN (weil noch keine MX/SMTP/Mailbox-Checks)
// - ungültig => REJECTED
//
// Sobald du MX/SMTP/Mailbox drin hast, wird hier daraus ACCEPTED / UNCERTAIN / REJECTED.
func ValidateParsed(items []ParsedEmail) []ParsedEmail {
	out := make([]ParsedEmail, 0, len(items))

	for _, it := range items {
		if !simpleEmailSyntax.MatchString(it.Email) {
			it.Status = StatusRejected
			it.Reason = "invalid_syntax"
			out = append(out, it)
			continue
		}

		// Platzhalter-Logik:
		// Syntax OK, aber ohne Netzprüfung => UNCERTAIN
		it.Status = StatusUncertain
		it.Reason = "syntax_ok_no_network_checks_yet"

		out = append(out, it)
	}

	return out
}
