package client

import (
	"maildefender/notifier/internal/connectors"
	"maildefender/notifier/internal/formatters"
)

var smtpClient Client

func SetSmtpClient(client Client) {
	smtpClient = client
}

func GetSmtpClient() Client {
	return smtpClient
}

type Client interface {
	Send(content any) error
}

type clientImpl struct {
	formatter formatters.Formatter
	connector connectors.Connector
}

func NewClient(formatter formatters.Formatter, connector connectors.Connector) Client {
	return &clientImpl{
		formatter: formatter,
		connector: connector,
	}
}

func (c clientImpl) Send(content any) error {
	recipients, formattedContent, err := c.formatter.Format(content)
	if err != nil {
		return err
	}

	return c.connector.Send(recipients, formattedContent)
}
