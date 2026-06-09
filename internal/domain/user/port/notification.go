package port

//go:generate mockgen -destination=mocks/mock_notification_service.go -package=mocks github.com/rulzi/hexa-go/internal/domain/user/port NotificationService

import "context"

// NotificationService is a port for sending notifications (e.g., emails).
type NotificationService interface {
	SendWelcomeEmail(ctx context.Context, email, name string) error
}
