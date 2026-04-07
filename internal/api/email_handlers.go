package api

import (
	"fmt"
	"html/template"
	"net/mail"
	"strings"
	"time"

	"github.com/kessaro/shaker"
	"github.com/sirupsen/logrus"

	"maildefender/notifier/internal/client"
	"maildefender/notifier/internal/configuration"
	"maildefender/notifier/internal/models"
)

const (
	maxSubjectLength = 255
	maxBodyLength    = 5 * 1024 * 1024 // 5MB
)

type sendMailIn struct {
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	ReplyTo     string   `json:"replyTo"`
	ThreadTopic string   `json:"threadTopic"`
}

func sendEmail(c *shaker.Context, in *sendMailIn) error {
	// Validate required fields
	if len(in.To) == 0 {
		return ErrMissingFields
	}
	if in.Subject == "" {
		return ErrMissingFields
	}
	if in.Body == "" {
		return ErrMissingFields
	}

	// Validate field lengths
	if len(in.Subject) > maxSubjectLength {
		return ErrMissingFields // Return mapped error for shaker handling
	}
	if len(in.Body) > maxBodyLength {
		return ErrMissingFields // Return mapped error for shaker handling
	}

	// Validate email addresses
	for i, email := range in.To {
		if _, err := mail.ParseAddress(email); err != nil {
			logrus.WithFields(logrus.Fields{
				"index": i,
				"email": email,
			}).WithError(err).Warn("invalid recipient email address")
			return ErrInvalidEmail // Return mapped error for shaker handling
		}
	}

	if in.ReplyTo != "" {
		if _, err := mail.ParseAddress(in.ReplyTo); err != nil {
			logrus.WithField("replyTo", in.ReplyTo).WithError(err).Warn("invalid reply-to address")
			return ErrInvalidEmail // Return mapped error for shaker handling
		}
	}

	smtpClient := client.GetSmtpClient()

	logrus.WithFields(logrus.Fields{
		"to_count": len(in.To),
		"subject":  in.Subject,
	}).Info("Sending email")

	// Sanitize From field to prevent header injection
	fromAddress := configuration.SmtpConfiguration().Authentication.Username
	fromAddress = strings.TrimSpace(fromAddress)
	fromAddress = strings.ReplaceAll(fromAddress, "\n", "")
	fromAddress = strings.ReplaceAll(fromAddress, "\r", "")

	mailStructure := models.MailStructure{
		To:          in.To,
		Subject:     in.Subject,
		ReplyTo:     in.ReplyTo,
		ThreadTopic: in.ThreadTopic,
		Body:        template.HTML(in.Body),
		From:        fromAddress,
		Date:        time.Now(),
	}

	if err := smtpClient.Send(mailStructure); err != nil {
		logrus.WithFields(logrus.Fields{
			"to_count":  len(in.To),
			"subject":   in.Subject,
			"recipient": in.To[0],
		}).WithError(err).Error("failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	logrus.WithField("subject", in.Subject).Info("email sent successfully")
	return nil
}
