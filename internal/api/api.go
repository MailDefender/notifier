package api

import (
	"net/http"

	"github.com/kessaro/shaker"

	"maildefender/notifier/internal/client"
)

var api = shaker.NewShaker(&shaker.MappedErrors{
	ErrMissingFields: http.StatusBadRequest,
	ErrInvalidEmail:  http.StatusBadRequest,
})

var additionalVariables map[string]any = map[string]any{}

func SetSmtpClient(client client.Client) {
	additionalVariables["smtpClient"] = client
}

func init() {
	api.Post("/v1/send/email", sendEmail, http.StatusOK)
}

// Run starts the HTTP server for the API
func Run() error {
	return api.Shake()
}
