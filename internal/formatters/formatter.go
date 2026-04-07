package formatters

// Formatter defines the interface for converting mail content into a sendable format.
type Formatter interface {
	// Format converts the input content into email recipients and formatted email content.
	// It returns the list of recipients, the formatted email content as a string, and any error that occurred.
	Format(content any) ([]string, string, error)
}
