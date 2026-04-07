package configuration

import "errors"

type AuthenticationConfiguration struct {
	Username string
	Password string
}

type ServerConfiguration struct {
	Host string
	Port int
}

type MailServerConfiguration struct {
	Server         ServerConfiguration
	Authentication AuthenticationConfiguration
}

func (s MailServerConfiguration) Check() error {
	if err := s.Server.Check(); err != nil {
		return err
	}

	if err := s.Authentication.Check(); err != nil {
		return err
	}

	return nil
}

func (s ServerConfiguration) Check() error {
	if s.Host == "" {
		return errors.New("SMTP host not configured")
	}

	if s.Port <= 0 || s.Port > 65535 {
		return errors.New("SMTP port must be between 1 and 65535")
	}
	return nil
}

func (s AuthenticationConfiguration) Check() error {
	if s.Username == "" {
		return errors.New("SMTP username not configured")
	}

	if s.Password == "" {
		return errors.New("SMTP password not configured")
	}
	return nil
}
