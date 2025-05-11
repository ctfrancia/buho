package ports

import (
// "github.com/ctfrancia/buho/internal/core/domain"
)

// EmailSender is the interface for sending emails and is the Outbound interface
type EmailSender interface {
	SendEmail(email string, subject string, body string) error
}
