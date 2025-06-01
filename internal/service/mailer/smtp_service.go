package mailer

import (
	"bytes"
	"context"
	"dailyalu-server/internal/service/mailer/domain"
	"dailyalu-server/pkg/mailer/smtp"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
)

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
	content, err := m.GetEmailHTML(emailVerificationData, "verification.html")
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

func (m *SmtpMailerService) GetEmailHTML(data any, templateName string) (string, error) {
	// Get the template file path
	_, filename, _, _ := runtime.Caller(0)
	templateDir := filepath.Dir(filename)
	templateFile := filepath.Join(templateDir, "templates/html", templateName)

	// Parse the template
	tmpl, err := template.ParseFiles(templateFile)
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
