package connectors

import (
	_ "embed"
	"errors"
	"fmt"
	nativeSmtp "net/smtp"
)

type SmtpServerConfiguration struct {
	Host string
	Port int16
}

type SmtpAuthenticationConfiguration struct {
	Username string
	Password string
}

type SmtpConfiguration struct {
	Server         SmtpServerConfiguration
	Authentication SmtpAuthenticationConfiguration
}

type smtp struct {
	configuration SmtpConfiguration
	auth          nativeSmtp.Auth
}

func NewSmtpConnector() smtp {
	return smtp{}
}

func (c *smtp) Connect(configuration any) error {
	var ok bool
	if c.configuration, ok = configuration.(SmtpConfiguration); !ok {
		return errors.New("cannot cast smtp configuration")
	}

	c.auth = nativeSmtp.PlainAuth("", c.configuration.Authentication.Username, c.configuration.Authentication.Password, c.configuration.Server.Host)
	return nil
}

func (c smtp) Send(recipients []string, content string) error {
	addr := fmt.Sprintf("%s:%d", c.configuration.Server.Host, c.configuration.Server.Port)

	return nativeSmtp.SendMail(
		addr,
		c.auth,
		c.configuration.Authentication.Username,
		recipients,
		[]byte(content))
}
