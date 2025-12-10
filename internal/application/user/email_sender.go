package user

import "context"

// EmailSender is an interface for sending emails (external service port)
type EmailSender interface {
	SendWelcomeEmail(ctx context.Context, email, name string) error
}

