package api

import "errors"

// ErrMissingFields is returned when required email fields are missing.
var (
	ErrMissingFields = errors.New("missing required fields: to, subject, body")
	ErrInvalidEmail  = errors.New("invalid email address")
)
