package main

type ParsedEmail struct {
	Email     string `json:"email"`
	Duplicate bool   `json:"duplicate"`
}
