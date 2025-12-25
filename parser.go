package main

import (
	"regexp"
	"sort"
	"strings"
)

// solide Email-Erkennung (nicht perfekt wie RFC, aber praxistauglich)
var emailRegex = regexp.MustCompile(`(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}`)

func ParseEmails(text string, keepDuplicates bool) []ParsedEmail {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")

	seen := make(map[string]bool)
	out := make([]ParsedEmail, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Falls jemand doch Kommentarzeilen drin hat
		if strings.HasPrefix(line, "#") {
			continue
		}

		matches := emailRegex.FindAllString(line, -1)
		if len(matches) == 0 {
			continue
		}

		for _, raw := range matches {
			email := normalizeEmail(raw)
			if email == "" {
				continue
			}

			dup := seen[email]
			if !dup {
				seen[email] = true
			}

			// Wenn Duplikate NICHT behalten: nur erste nehmen
			if dup && !keepDuplicates {
				continue
			}

			out = append(out, ParsedEmail{
				Email:     email,
				Duplicate: dup,
				Status:    StatusUncertain, // Parser allein bewertet noch nicht "echt"
				Reason:    "",
			})
		}
	}

	return out
}

func normalizeEmail(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	// Ein paar häufige Klammern/Trenner entfernen
	s = strings.Trim(s, "<>;,")
	return s
}

// Hilfsfunktionen für sortierte TXT-Ausgabe nach "B" (MX+SMTP+Mailbox)
// -> wir haben noch keine echten Netzchecks, aber das Sortier-Gerüst steht.
// Wenn später echte Checks da sind, wird das automatisch sinnvoll.
type sortKey struct {
	email  string
	weight int
}

func sortByWeightDesc(items []ParsedEmail) []string {
	keys := make([]sortKey, 0, len(items))
	for _, it := range items {
		keys = append(keys, sortKey{
			email:  it.Email,
			weight: statusWeight(it.Status),
		})
	}

	sort.SliceStable(keys, func(i, j int) bool {
		// höheres Gewicht zuerst
		if keys[i].weight != keys[j].weight {
			return keys[i].weight > keys[j].weight
		}
		// bei Gleichstand alphabetisch (stabil)
		return keys[i].email < keys[j].email
	})

	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, k.email)
	}
	return out
}

func statusWeight(status string) int {
	switch status {
	case StatusAccepted:
		return 300
	case StatusUncertain:
		return 200
	case StatusRejected:
		return 100
	default:
		return 0
	}
}
