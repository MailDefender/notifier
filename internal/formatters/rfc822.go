package formatters

import (
	"bytes"
	"errors"
	"html/template"

	"github.com/aymerick/douceur/inliner"
	"github.com/sirupsen/logrus"

	"maildefender/notifier/internal/models"
	"maildefender/notifier/internal/templates"
)

type rfc822 struct{}

func NewRfc822Formatter() rfc822 {
	return rfc822{}
}

func (f rfc822) Format(content any) ([]string, string, error) {
	var input models.MailStructure

	{
		var ok bool
		if input, ok = content.(models.MailStructure); !ok {
			logrus.Error("cannot format input, invalid input object")
			return nil, "", errors.New("cannot format input, invalid input object")
		}
	}

	htmlRaw, err := inliner.Inline(string(input.Body))
	if err != nil {
		logrus.WithError(err).Error("cannot inline CSS from input body")
	}

	input.Body = template.HTML(htmlRaw)

	var bufOut bytes.Buffer
	if err := templates.MailFrameTemplate().Execute(&bufOut, input); err != nil {
		logrus.WithError(err).Error("cannot execute html template")
		return nil, "", err
	}

	return input.To, bufOut.String(), nil
}
