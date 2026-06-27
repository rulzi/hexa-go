package user

import (
	"context"
	"testing"

	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmailSenderImpl(t *testing.T) {
	sender := NewEmailSenderImpl(logger.NewSimpleLogger())
	require.NotNil(t, sender)
}

func TestEmailSenderImpl_SendWelcomeEmail(t *testing.T) {
	sender := NewEmailSenderImpl(logger.NewSimpleLogger())

	err := sender.SendWelcomeEmail(context.Background(), "user@example.com", "Test User")

	assert.NoError(t, err)
}
