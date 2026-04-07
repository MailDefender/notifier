package formatters

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/aymerick/douceur/inliner"
	"github.com/sirupsen/logrus"

	"maildefender/notifier/internal/models"
	"maildefender/notifier/internal/templates"
)

// NewRfc822Formatter creates a new RFC822 email formatter.
type rfc822 struct{}

// NewRfc822Formatter creates a new RFC822 email formatter.
func NewRfc822Formatter() rfc822 {
	return rfc822{}
}

// Format converts a MailStructure into RFC822 format recipients and email content.
func (f rfc822) Format(content any) ([]string, string, error) {
	input, ok := content.(models.MailStructure)
	if !ok {
		logrus.WithField("got_type", fmt.Sprintf("%T", content)).Error("invalid email content type, expected MailStructure")
		return nil, "", fmt.Errorf("invalid email content type: expected models.MailStructure, got %T", content)
	}

	htmlRaw, err := inliner.Inline(string(input.Body))
	if err != nil {
		logrus.WithError(err).Warn("failed to inline CSS from input body, proceeding with original HTML")
		htmlRaw = string(input.Body)
	}

	input.Body = template.HTML(htmlRaw)

	var bufOut bytes.Buffer
	if err := templates.MailFrameTemplate().Execute(&bufOut, input); err != nil {
		logrus.WithError(err).Error("cannot execute html template")
		return nil, "", err
	}

	return input.To, bufOut.String(), nil
}
