package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"maildefender/notifier/internal/api"
	"maildefender/notifier/internal/client"
	"maildefender/notifier/internal/configuration"
	"maildefender/notifier/internal/connectors"
	"maildefender/notifier/internal/formatters"
	"maildefender/notifier/internal/utils"
)

func main() {
	logrus.Info("Starting notifier...")

	logrus.Info("Retrieving configuration")

	config := smtpConfiguration()
	if err := config.Check(); err != nil {
		logrus.WithError(err).Fatal("invalid SMTP configuration")
	}

	smtpConnector := connectors.NewSmtpConnector()
	if err := smtpConnector.Connect(config); err != nil {
		logrus.WithError(err).Fatal("cannot connect to given smtp server")
	}
	client.SetSmtpClient(client.NewClient(formatters.NewRfc822Formatter(), &smtpConnector))

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Run API in a goroutine
	if err := api.Run(); err != nil {
		logrus.WithError(err).Error("API stopped unexpectedly")
	}
}

func smtpConfiguration() configuration.MailServerConfiguration {
	return configuration.MailServerConfiguration{
		Server: configuration.ServerConfiguration{
			Host: utils.GetEnvString("SMTP_HOST", ""),
			Port: utils.GetEnvInt("SMTP_PORT", 0),
		},
		Authentication: configuration.AuthenticationConfiguration{
			Username: utils.GetEnvString("SMTP_USERNAME", ""),
			Password: utils.GetEnvString("SMTP_PASSWORD", ""),
		},
	}
}
