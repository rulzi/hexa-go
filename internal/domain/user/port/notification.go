package port

import "context"

// NotificationService is a port for sending notifications (e.g., emails).
type NotificationService interface {
	SendWelcomeEmail(ctx context.Context, email, name string) error
}
