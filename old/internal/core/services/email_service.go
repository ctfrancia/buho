package services

import (
	"github.com/ctfrancia/buho/internal/core/ports"
)

type EmailService struct {
	notify ports.EmailSender
}

func NewEmailService(tp ports.EmailSender) *EmailService {
	return &EmailService{
		notify: tp,
	}
}

func (es *EmailService) SendEmail(email string, subject string, body string) error {
	return es.notify.SendEmail(email, subject, body)
}
