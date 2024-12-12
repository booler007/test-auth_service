package service

import (
	"net/smtp"
	"os"

	"authentication_medods/cmd/api/config"
)

func NewEmailService(eml *config.Email) *EmailService {
	return &EmailService{eml.Account, eml.Password, eml.Host, eml.Port}
}

type EmailService struct {
	Account  string
	Password string
	Host     string
	Port     string
}

func (e *EmailService) NotificationNewIP(userEmail string) error {
	subject := os.Getenv("Subject: Attention! Detected new IP\n")
	body := "You have requested a new access token from a new IP"
	message := []byte(subject + body)
	auth := smtp.PlainAuth("", e.Account, e.Password, e.Host)
	address := e.Host + ":" + e.Port

	return smtp.SendMail(address, auth, e.Account, []string{userEmail}, message)
}
