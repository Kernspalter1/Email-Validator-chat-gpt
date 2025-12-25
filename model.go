package main

// Status-Werte (so wie du wolltest)
const (
	StatusAccepted  = "ACCEPTED"
	StatusUncertain = "UNCERTAIN"
	StatusRejected  = "REJECTED"
)

// Ein einzelner Eintrag aus dem Parser (GUI darf auch Duplikate sehen)
type ParsedEmail struct {
	Email     string `json:"email"`
	Duplicate bool   `json:"duplicate"`
	Status    string `json:"status"`
	Reason    string `json:"reason,omitempty"`
}

// Request vom Frontend
type ParseRequest struct {
	Text           string `json:"text"`
	KeepDuplicates bool   `json:"keep_duplicates"`
}

// Response ans Frontend
type ParseResponse struct {
	Items            []ParsedEmail `json:"items"`              // Liste für GUI (kann Duplikate enthalten)
	UniqueSorted      []string      `json:"unique_sorted"`      // für TXT Export: accepted oder accepted+uncertain (sortiert)
	AcceptedSorted    []string      `json:"accepted_sorted"`    // nur accepted, sortiert
	AcceptedUncertain []string      `json:"accepted_uncertain"` // accepted + uncertain, sortiert
}
