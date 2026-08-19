package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

// Mailer sends transactional email. The SMTP transport is configured in main;
// this interface is what auth/user depend on.
type Mailer interface {
	Send(ctx context.Context, to, subject, text string) error
}

// Noop is the default for local dev — logs instead of sending.
type Noop struct{ Log func(msg string) }

func (m Noop) Send(_ context.Context, to, subject, text string) error {
	if m.Log != nil {
		m.Log("mail: to=" + to + " subject=" + subject)
	}
	return nil
}

// SMTP is a plain/STARTTLS SMTP sender built on net/smtp. It is used only when
// SMTP_HOST is configured; otherwise main wires Noop.
type SMTP struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

func (s *SMTP) Send(_ context.Context, to, subject, text string) error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	msg := "From: " + s.From + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + text + "\r\n"
	var auth smtp.Auth
	if s.User != "" {
		auth = smtp.PlainAuth("", s.User, s.Password, s.Host)
	}
	// net/smtp sends plaintext on 25; on 465/587 the sendmail helper still
	// works but does not STARTTLS. For TLS you'd use a tls client; the
	// envelope/relay for most providers tolerates this in test setups.
	toList := strings.Split(to, ",")
	return smtp.SendMail(addr, auth, s.From, toList, []byte(msg))
}
