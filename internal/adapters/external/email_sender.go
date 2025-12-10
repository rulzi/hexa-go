package external

import (
	"context"
	"fmt"
	"log"
)

// EmailSenderImpl is the implementation of the email sender (external service adapter)
type EmailSenderImpl struct {
	// In a real implementation, this would have SMTP config, API keys, etc.
}

// NewEmailSenderImpl creates a new EmailSenderImpl
func NewEmailSenderImpl() *EmailSenderImpl {
	return &EmailSenderImpl{}
}

// SendWelcomeEmail sends a welcome email to a new user
func (e *EmailSenderImpl) SendWelcomeEmail(ctx context.Context, email, name string) error {
	// In a real implementation, this would send an actual email
	// For now, we'll just log it
	log.Printf("Sending welcome email to %s (%s)", name, email)
	
	// Simulate email sending
	// In production, you would integrate with:
	// - SMTP server
	// - SendGrid, Mailgun, AWS SES, etc.
	
	fmt.Printf("[EMAIL] Welcome %s! Your account has been created successfully.\n", name)
	return nil
}

