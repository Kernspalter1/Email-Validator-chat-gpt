package main

type MailboxPlausibility string

const (
    PlausibleLikely    MailboxPlausibility = "wahrscheinlich"
    PlausibleUncertain MailboxPlausibility = "unsicher"
    PlausibleUnknown   MailboxPlausibility = "unbekannt"
)
