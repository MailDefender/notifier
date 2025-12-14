package templates

import (
	_ "embed"
	"html/template"
	"strings"
)

var (
	//go:embed mail-frame.eml
	mailFrameContent string

	mailFrameTemplate *template.Template = template.Must(template.New("").Funcs(template.FuncMap{
		"serializeEmails": func(val []string) string {
			return strings.Join(val, ",")
		},
		"formatReplyTo": func(val string) string {
			if !strings.HasPrefix(val, "<") {
				val = "<" + val
			}
			if !strings.HasSuffix(val, ">") {
				val = val + ">"
			}
			return val
		},
	}).Parse(mailFrameContent))
)

func MailFrameTemplate() *template.Template {
	return mailFrameTemplate
}
