package formatters

import (
	"bytes"
	"errors"

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
			return nil, "", errors.New("cannot format input, invalid input object")
		}
	}

	var bufOut bytes.Buffer
	if err := templates.MailFrameTemplate().Execute(&bufOut, input); err != nil {
		return nil, "", err
	}

	return input.To, bufOut.String(), nil
}
