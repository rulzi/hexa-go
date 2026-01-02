package user

import "context"

// NotificationService is a port for sending notifications (e.g., emails)
type NotificationService interface {
	// SendWelcomeEmail sends a welcome email to a new user
	SendWelcomeEmail(ctx context.Context, email, name string) error
}

