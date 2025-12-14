package models

import (
	"html/template"
	"time"
)

type MailStructure struct {
	ReplyTo     string
	To          []string
	From        string
	Subject     string
	ThreadTopic string
	Body        template.HTML
	Date        time.Time
}

func (e MailStructure) FormatDate() string {
	return e.Date.Format(time.RFC822)
}
