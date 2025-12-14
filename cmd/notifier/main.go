package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"maildefender/notifier/internal/api"
	"maildefender/notifier/internal/client"
	"maildefender/notifier/internal/connectors"
	"maildefender/notifier/internal/formatters"
	"maildefender/notifier/internal/utils"
)

func main() {
	logrus.Info("Starting notifier...")

	logrus.Info("Retrieving configuration")

	smtpConnector := connectors.NewSmtpConnector()
	if err := smtpConnector.Connect(smtpConfiguration()); err != nil {
		logrus.WithError(err).Fatal("cannot connect to given smtp server")
		os.Exit(1)
	}
	api.SetSmtpClient(client.NewClient(formatters.NewRfc822Formatter(), &smtpConnector))

	api.RegisterMiddlewares()
	api.RegisterHandlers()
	if err := api.Run(); err != nil {
		logrus.WithError(err).Error("API stopped")
	}
}

func smtpConfiguration() connectors.SmtpConfiguration {
	return connectors.SmtpConfiguration{
		Server: connectors.SmtpServerConfiguration{
			Host: utils.GetEnvString("SMTP_HOST", ""),
			Port: int16(utils.GetEnvInt("SMTP_PORT", 0)),
		},
		Authentication: connectors.SmtpAuthenticationConfiguration{
			Username: utils.GetEnvString("SMTP_USERNAME", ""),
			Password: utils.GetEnvString("SMTP_PASSWORD", ""),
		},
	}
}
