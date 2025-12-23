package main

type MailboxPlausibility string

const (
	PlausibleLikely    MailboxPlausibility = "wahrscheinlich"
	PlausibleUncertain MailboxPlausibility = "unsicher"
	PlausibleUnknown   MailboxPlausibility = "unbekannt"
)

type Result struct {
	Email   string              `json:"email"`
	MX      string              `json:"mx"`   // "ja" / "nein"
	SMTP    string              `json:"smtp"` // "ja" / "nein"
	Mailbox MailboxPlausibility `json:"mailbox"`
	Detail  string              `json:"detail,omitempty"` // optional
}

type JobStatus struct {
	Running bool   `json:"running"`
	Total   int    `json:"total"`
	Done    int    `json:"done"`
	Percent int    `json:"percent"`
	Message string `json:"message"`
}
