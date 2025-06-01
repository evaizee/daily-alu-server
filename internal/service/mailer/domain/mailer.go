package domain

import "context"

const (
	// Email templates
	VerificationEmailTemplate = "verification_email.html"

	//constants
	NoReplyEmail = "no-reply@dailyalu.mom"

	//subjects
	WelcomeSubject = "Welcome to DailyAlu!"
)

type EmailVerificationData struct {
	Name            string
	VerificationURL string
	To string
}

type IMailerService interface {
	SendVerificationEmail(ctx context.Context, data *EmailVerificationData) (error)
	GetEmailHTML(data any, templateName string) (string, error)
}