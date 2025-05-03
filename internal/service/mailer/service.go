package mailer

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type MailerService struct {
	SesClient *sesv2.Client
}

func NewMailerService(sesClient *sesv2.Client) *MailerService {
	return &MailerService{
		SesClient: sesClient,
	}
}

func (m *MailerService) Send(ctx context.Context, from, to, subject, textMessage, htmlMessage string) (output string, err error) {

	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String("noreply@mail.dailyalu.mom"),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Text: &types.Content{
						Data: aws.String(textMessage),
					},
					// Html: &types.Content{
					// 	Data: &htmlMessage,
					// },
				},
				Subject: &types.Content{
					Data: aws.String(subject),
				},
			},
		},
	}

	_,err = m.SesClient.SendEmail(ctx, input)

	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return "ok", nil
}
