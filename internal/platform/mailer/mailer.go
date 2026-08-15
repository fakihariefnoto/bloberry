package mailer

import "context"

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
