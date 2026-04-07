package api

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"maildefender/notifier/internal/client"
	"maildefender/notifier/internal/configuration"
	"maildefender/notifier/internal/connectors"
	"maildefender/notifier/internal/formatters"
)

func testingSmtpConfiguration() configuration.MailServerConfiguration {
	return configuration.MailServerConfiguration{
		Server: configuration.ServerConfiguration{
			Host: "localhost",
			Port: 1025,
		},
		Authentication: configuration.AuthenticationConfiguration{
			Username: "test",
			Password: "test",
		},
	}
}

func TestMain(m *testing.M) {
	logrus.Info("Starting notifier...")

	logrus.Info("Retrieving configuration")

	smtpConnector := connectors.NewSmtpConnector()
	if err := smtpConnector.Connect(testingSmtpConfiguration()); err != nil {
		logrus.WithError(err).Fatal("cannot connect to given smtp server")
	}
	client.SetSmtpClient(client.NewClient(formatters.NewRfc822Formatter(), &smtpConnector))

	os.Exit(m.Run())
}
