package user

import (
	"context"

	userport "github.com/rulzi/hexa-go/internal/domain/user/port"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// EmailSenderImpl implements NotificationService (external service adapter)
type EmailSenderImpl struct {
	logger logger.Logger
}

// NewEmailSenderImpl creates a new EmailSenderImpl
func NewEmailSenderImpl(appLogger logger.Logger) *EmailSenderImpl {
	return &EmailSenderImpl{logger: appLogger}
}

// SendWelcomeEmail implements NotificationService interface
func (e *EmailSenderImpl) SendWelcomeEmail(ctx context.Context, email, name string) error {
	// In a real implementation, this would send an actual email
	e.logger.InfoWithFields("sending welcome email", map[string]interface{}{
		"email": email,
		"name":  name,
	})
	// Simulate email sending - in production integrate with SMTP, SendGrid, Mailgun, AWS SES, etc.
	return nil
}

// Ensure EmailSenderImpl implements userport.NotificationService
var _ userport.NotificationService = (*EmailSenderImpl)(nil)

