package mailer

import (
	"bytes"
	"context"
	"dailyalu-server/internal/service/mailer/domain"
	"dailyalu-server/pkg/mailer/smtp"
	"fmt"
	"html/template"
	"embed"
)

//go:embed templates/html/*
var emailTemplates embed.FS

// MailerService handles email sending via SMTP
type SmtpMailerService struct {
	Smtp *smtp.Smtp
}

// NewMailerService creates a new SMTP mailer service
func NewSmtpMailerService(smtp *smtp.Smtp) domain.IMailerService {
	return &SmtpMailerService{
		Smtp: smtp,
	}
}

func (m *SmtpMailerService) SendVerificationEmail(ctx context.Context, emailVerificationData *domain.EmailVerificationData) (err error) {
	content, err := m.getEmailHTML(emailVerificationData, "verification.html")
	if err != nil {
		return err
	}
	// Create the email data
	emailData := &smtp.SendEmailData{
		From:    domain.NoReplyEmail,
		To:      emailVerificationData.To,
		Subject: domain.WelcomeSubject,
		Content: content,
	}

	output, err := m.Smtp.Send(*emailData)
	if err != nil {
		return err
	}
	fmt.Println(output)

	return nil
}

func (m *SmtpMailerService) getEmailHTML(data any, templateName string) (string, error) {
	
	// Get the template file path
	tmpl, err := template.ParseFS(emailTemplates, "templates/html/" + templateName)
	if err != nil {
		return "", err
	}

	// Execute the template with data
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
