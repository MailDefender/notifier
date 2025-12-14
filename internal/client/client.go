package client

import (
	"maildefender/notifier/internal/connectors"
	"maildefender/notifier/internal/formatters"
)

type Client interface {
	Send(content any) error
}

type clientImpl struct {
	formatter formatters.Formatter
	connector connectors.Connector
}

func NewClient(formatter formatters.Formatter, connector connectors.Connector) clientImpl {
	return clientImpl{
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
