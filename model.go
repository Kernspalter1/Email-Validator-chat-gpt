package main

type Plausibility string

const (
	ACCEPTED  Plausibility = "ACCEPTED"
	UNCERTAIN Plausibility = "UNCERTAIN"
	REJECTED  Plausibility = "REJECTED"
)

type ParsedEmail struct {
	Email     string      `json:"email"`
	Duplicate bool        `json:"duplicate"`
	Status    Plausibility `json:"status"`
	Reason    string      `json:"reason"`
}
