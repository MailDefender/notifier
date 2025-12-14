package connectors

type Connector interface {
	Connect(configuration any) error
	Send(recipients []string, content string) error
}
