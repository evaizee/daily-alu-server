package smtp

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"time"

	"github.com/spf13/viper"
)

type Smtp struct {
	Config *Config
}

// Config holds the configuration for SMTP connections
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SendEmailData struct {
	From,
	To,
	Subject,
	Content string
}

func InitSmtp() *Smtp {
	return &Smtp{
		Config: NewConfig(),
	}
}

// NewConfig creates a new SMTP configuration from viper settings
func NewConfig() *Config {
	return &Config{
		Host:     viper.GetString("smtp.host"),
		Port:     viper.GetInt("smtp.port"),
		Username: viper.GetString("smtp.username"),
		Password: viper.GetString("smtp.password"),
		From:     viper.GetString("smtp.from"),
	}
}

// GetSMTPAuth returns the SMTP authentication object
func (s *Smtp) GetSMTPAuth() smtp.Auth {
	return smtp.PlainAuth("", s.Config.Username, s.Config.Password, s.Config.Host)
}

// GetSMTPAddress returns the SMTP server address with port
func (s *Smtp) GetSMTPAddress() string {
	return fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)
}

func (s *Smtp) Send(data SendEmailData) (output string, err error) {
	// If from is empty, use the default from address
	from := data.From
	if from == "" {
		from = s.Config.From
	}

	// Create email message with proper headers
	body := &bytes.Buffer{}

	// Setup multipart
	writer := multipart.NewWriter(body)
	boundary := writer.Boundary()

	// Set headers
	headers := fmt.Sprintf(
		"From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/alternative; boundary=%s\r\n\r\n",
		from, data.To, data.Subject, boundary)

	// Add text part
	textPart, err := writer.CreatePart(map[string][]string{
		"Content-Type":              {"text/plain; charset=UTF-8"},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create text part: %w", err)
	}

	qpWriter := quotedprintable.NewWriter(textPart)
	// _, err = qpWriter.Write([]byte(data.TextMessage))
	// if err != nil {
	// 	return "", fmt.Errorf("failed to write text message: %w", err)
	// }
	// qpWriter.Close()

	// Add HTML part if provided
	if data.Content != "" {
		htmlPart, err := writer.CreatePart(map[string][]string{
			"Content-Type":              {"text/html; charset=UTF-8"},
			"Content-Transfer-Encoding": {"quoted-printable"},
		})
		if err != nil {
			return "", fmt.Errorf("failed to create HTML part: %w", err)
		}

		qpWriter = quotedprintable.NewWriter(htmlPart)
		_, err = qpWriter.Write([]byte(data.Content))
		if err != nil {
			return "", fmt.Errorf("failed to write HTML message: %w", err)
		}
		qpWriter.Close()
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Connect to the SMTP server and send the email
	auth := s.GetSMTPAuth()
	smtpAddr := s.GetSMTPAddress()

	// Send the email
	err = smtp.SendMail(
		smtpAddr,
		auth,
		from,
		[]string{data.To},
		[]byte(headers+body.String()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return fmt.Sprintf("Email sent to %s at %s", data.To, time.Now().Format(time.RFC3339)), nil
}