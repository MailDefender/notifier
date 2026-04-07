package connectors

// Connector defines the interface for sending emails through different backends (SMTP, etc.).
type Connector interface {
	// Connect establishes a connection to the email service using the provided configuration.
	Connect(configuration any) error
	// Send sends an email to the specified recipients with the given content.
	Send(recipients []string, content string) error
}
