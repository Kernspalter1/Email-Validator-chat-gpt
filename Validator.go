package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/mail"
	"strings"
	"time"
)

type Validator struct {
	DNSMXTimeout     time.Duration
	SMTPDialTimeout  time.Duration
	SMTPReadTimeout  time.Duration
	FromAddress      string
}

func NewValidator() *Validator {
	return &Validator{
		DNSMXTimeout:    8 * time.Second,
		SMTPDialTimeout: 15 * time.Second,
		SMTPReadTimeout: 20 * time.Second,
		FromAddress:     "validator@localhost",
	}
}

func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func parseDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func (v *Validator) SyntaxOK(email string) bool {
	email = normalizeEmail(email)
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (v *Validator) LookupMX(domain string) ([]*net.MX, error) {
	ctx, cancel := context.WithTimeout(context.Background(), v.DNSMXTimeout)
	defer cancel()

	type resp struct {
		mx  []*net.MX
		err error
	}
	ch := make(chan resp, 1)

	go func() {
		mx, err := net.LookupMX(domain)
		ch <- resp{mx: mx, err: err}
	}()

	select {
	case r := <-ch:
		return r.mx, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func pickBestMX(mx []*net.MX) string {
	if len(mx) == 0 {
		return ""
	}
	best := mx[0]
	for _, m := range mx[1:] {
		if m.Pref < best.Pref {
			best = m
		}
	}
	return strings.TrimSuffix(best.Host, ".")
}

func (v *Validator) SMTPCheck(mxHost, rcptEmail string) (bool, MailboxPlausibility, string) {
	addr := net.JoinHostPort(mxHost, "25")

	dialer := net.Dialer{Timeout: v.SMTPDialTimeout}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return false, PlausibleUnknown, "dial_failed: " + err.Error()
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(v.SMTPReadTimeout))

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	readLine := func() (string, error) {
		line, err := r.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(line), nil
	}
	writeLine := func(s string) error {
		if _, err := w.WriteString(s + "\r\n"); err != nil {
			return err
		}
		return w.Flush()
	}
	readCode := func(line string) int {
		if len(line) < 3 {
			return 0
		}
		var code int
		_, _ = fmt.Sscanf(line[:3], "%d", &code)
		return code
	}

	banner, err := readLine()
	if err != nil {
		return false, PlausibleUnknown, "banner_read_failed: " + err.Error()
	}
	_ = readCode(banner)

	_ = writeLine("EHLO localhost")
	for {
		line, err := readLine()
		if err != nil {
			return true, PlausibleUnknown, "ehlo_read_failed: " + err.Error()
		}
		c := readCode(line)
		if c >= 400 && c < 500 {
			return true, PlausibleUnknown, "ehlo_tempfail: " + line
		}
		if c >= 500 {
			return true, PlausibleUnknown, "ehlo_reject: " + line
		}
		if strings.HasPrefix(line, "250 ") {
			break
		}
	}

	_ = writeLine(fmt.Sprintf("MAIL FROM:<%s>", v.FromAddress))
	line, err := readLine()
	if err != nil {
		return true, PlausibleUnknown, "mailfrom_read_failed: " + err.Error()
	}
	c := readCode(line)
	if c >= 400 && c < 500 {
		return true, PlausibleUnknown, "mailfrom_tempfail: " + line
	}
	if c >= 500 {
		return true, PlausibleUnknown, "mailfrom_reject: " + line
	}

	_ = writeLine(fmt.Sprintf("RCPT TO:<%s>", rcptEmail))
	line, err = readLine()
	if err != nil {
		return true, PlausibleUnknown, "rcpt_read_failed: " + err.Error()
	}
	c = readCode(line)

	switch {
	case c == 250 || c == 251:
		return true, PlausibleLikely, "rcpt_ok: " + line
	case c >= 400 && c < 500:
		return true, PlausibleUnknown, "rcpt_tempfail: " + line
	case c >= 500:
		return true, PlausibleUncertain, "rcpt_reject: " + line
	default:
		return true, PlausibleUnknown, "rcpt_unknown: " + line
	}
}
