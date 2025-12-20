package templates

import (
	_ "embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/google/uuid"
)

var (
	//go:embed mail-frame.eml
	mailFrameContent string

	mailFrameTemplate *template.Template = template.Must(template.New("").Funcs(template.FuncMap{
		"serializeEmails": func(val []string) template.HTML {
			return template.HTML(strings.Join(val, ","))
		},
		"formatReplyTo": func(val string) template.HTML {
			if !strings.HasPrefix(val, "<") {
				val = "<" + val
			}
			if !strings.HasSuffix(val, ">") {
				val = val + ">"
			}
			return template.HTML(val)
		},
		"generateMsgId": func() template.HTML {
			return template.HTML(fmt.Sprintf("<%s@%s>", uuid.New().String(), "notifier-mail-defender"))
		},
	}).Parse(mailFrameContent))
)

func MailFrameTemplate() *template.Template {
	return mailFrameTemplate
}
