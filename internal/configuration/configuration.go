package configuration

import (
	"maildefender/notifier/internal/utils"
)

var smtpConfiguration MailServerConfiguration = MailServerConfiguration{
	Server: ServerConfiguration{
		Host: utils.GetEnvString("SMTP_HOST", ""),
		Port: utils.GetEnvInt("SMTP_PORT", 0),
	},
	Authentication: AuthenticationConfiguration{
		Username: utils.GetEnvString("SMTP_USERNAME", ""),
		Password: utils.GetEnvString("SMTP_PASSWORD", ""),
	},
}

func SmtpConfiguration() MailServerConfiguration {
	return smtpConfiguration
}
