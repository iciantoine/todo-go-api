package server

import "github.com/iciantoine/todo-go-api/option"

type config struct {
	Database    option.Postgres
	Application option.Endpoint
}

// Option is a configurable parameter.
type Option func(*config) error

// WithApplicationAddress configures the application listen address.
func WithApplicationAddress(addr string, port string) Option {
	return func(cfg *config) error {
		cfg.Application.Addr = addr
		cfg.Application.Port = port
		return nil
	}
}

// WithDatabase configures the credentials to connect to local Postgres DB.
func WithDatabase(user, pass, host, port, name, sslmode string) Option {
	return func(cfg *config) error {
		cfg.Database.User = user
		cfg.Database.Pass = pass
		cfg.Database.Host = host
		cfg.Database.Port = port
		cfg.Database.Name = name
		cfg.Database.Mode = sslmode
		return nil
	}
}

// WithLogLevel configures the log level.
func WithLogLevel(lvl string) Option {
	return func(cfg *config) error {
		return option.ConfigureLogging(lvl)
	}
}
