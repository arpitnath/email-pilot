package services

import (
	"context"
	"encoding/base64"
	"errors"
	"log"

	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type EmailService struct {
	Token       *oauth2.Token
	OauthConfig *oauth2.Config
}

// FetchEmails retrieves a specified number of emails from the Gmail API
func (e *EmailService) FetchEmails(limit int) ([]map[string]string, error) {
	if limit <= 0 {
		return nil, errors.New("invalid email limit")
	}

	// Create Gmail client
	ctx := context.Background()
	client := e.OauthConfig.Client(ctx, e.Token)
	gmailService, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	// Fetch email list
	user := "me"
	msgs, err := gmailService.Users.Messages.List(user).MaxResults(int64(limit)).Do()
	if err != nil {
		return nil, err
	}

	// Process each email
	var emails []map[string]string
	for _, msg := range msgs.Messages {
		message, err := gmailService.Users.Messages.Get(user, msg.Id).Format("full").Do()
		if err != nil {
			log.Printf("Failed to fetch email with ID %s: %v", msg.Id, err)
			continue
		}

		// Decode subject and body
		var subject string
		for _, header := range message.Payload.Headers {
			if header.Name == "Subject" {
				subject = header.Value
				break
			}
		}

		body := ""
		if message.Payload.Body != nil {
			decodedBody, _ := base64.URLEncoding.DecodeString(message.Payload.Body.Data)
			body = string(decodedBody)
		}

		emails = append(emails, map[string]string{
			"subject": subject,
			"body":    body,
		})
	}

	return emails, nil
}
