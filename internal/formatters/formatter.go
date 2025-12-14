package formatters

type Formatter interface {
	Format(content any) ([]string, string, error)
}
