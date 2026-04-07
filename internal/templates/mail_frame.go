package templates

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var (
	//go:embed mail-frame.eml
	mailFrameContent string

	mailFrameTemplate *template.Template
	templateOnce      sync.Once
	templateErr       error
)

// initTemplate lazily initializes the email template with error handling
func initTemplate() error {
	var err error
	templateOnce.Do(func() {
		mailFrameTemplate, err = template.New("").Funcs(template.FuncMap{
			"serializeEmails": func(val []string) template.HTML {
				// Sanitize email addresses to prevent header injection
				sanitized := make([]string, len(val))
				for i, email := range val {
					email = strings.TrimSpace(email)
					email = strings.ReplaceAll(email, "\n", "")
					email = strings.ReplaceAll(email, "\r", "")
					sanitized[i] = email
				}
				return template.HTML(strings.Join(sanitized, ","))
			},
			"formatReplyTo": func(val string) template.HTML {
				// Sanitize the reply-to address to prevent header injection
				val = strings.TrimSpace(val)
				val = strings.ReplaceAll(val, "\n", "")
				val = strings.ReplaceAll(val, "\r", "")
				if !strings.HasPrefix(val, "<") {
					val = "<" + val
				}
				if !strings.HasSuffix(val, ">") {
					val = val + ">"
				}
				return template.HTML(val)
			},
			"generateMsgId": func() template.HTML {
				const messageDomain = "notifier-mail-defender"
				return template.HTML(fmt.Sprintf("<%s@%s>", uuid.New().String(), messageDomain))
			},
		}).Parse(mailFrameContent)
	})
	return err
}

// MailFrameTemplate returns the email template, initializing it lazily on first call
func MailFrameTemplate() *template.Template {
	if err := initTemplate(); err != nil {
		log.Fatalf("failed to initialize email template: %v", err)
	}
	return mailFrameTemplate
}
