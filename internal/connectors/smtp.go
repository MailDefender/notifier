package connectors

import (
	"fmt"
	"net"
	nativeSmtp "net/smtp"
	"strconv"
	"time"

	"maildefender/notifier/internal/configuration"
)

// SMTP is a connector implementation for sending emails via SMTP protocol.
type SMTP struct {
	configuration configuration.MailServerConfiguration
	auth          nativeSmtp.Auth
	timeout       time.Duration
}

const defaultTimeout = 30 * time.Second

// NewSmtpConnector creates a new SMTP connector with default timeout.
func NewSmtpConnector() SMTP {
	return SMTP{
		timeout: defaultTimeout,
	}
}

// Connect establishes a connection to the SMTP server using the provided configuration.
func (c *SMTP) Connect(config any) error {
	var ok bool
	if c.configuration, ok = config.(configuration.MailServerConfiguration); !ok {
		return fmt.Errorf("invalid configuration type: expected MailServerConfiguration, got %T", config)
	}

	c.auth = nativeSmtp.PlainAuth("", c.configuration.Authentication.Username, c.configuration.Authentication.Password, c.configuration.Server.Host)
	return nil
}

// Send sends an email with the given recipients and content, with timeout protection.
func (c SMTP) Send(recipients []string, content string) error {
	addr := net.JoinHostPort(c.configuration.Server.Host, strconv.Itoa(c.configuration.Server.Port))

	// Create a connection with timeout
	conn, err := net.DialTimeout("tcp", addr, c.timeout)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer conn.Close()

	// Create SMTP client from connection
	client, err := nativeSmtp.NewClient(conn, c.configuration.Server.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate
	if err := client.Auth(c.auth); err != nil {
		return fmt.Errorf("SMTP AUTH failed: %w", err)
	}

	// Set sender
	if err := client.Mail(c.configuration.Authentication.Username); err != nil {
		return fmt.Errorf("SMTP MAIL failed: %w", err)
	}

	// Set recipients
	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("SMTP RCPT failed for %s: %w", recipient, err)
		}
	}

	// Send message body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA failed: %w", err)
	}

	if _, err := w.Write([]byte(content)); err != nil {
		w.Close()
		return fmt.Errorf("failed to write email content: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close DATA: %w", err)
	}

	// Send QUIT command
	return client.Quit()
}
