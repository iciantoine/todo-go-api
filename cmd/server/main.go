package main

import (
	"context"
	"os"

	"github.com/iciantoine/todo-go-api/cmd"
	"github.com/iciantoine/todo-go-api/server"
)

func main() {
	if err := cmd.Run(run); err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	return server.Listen(ctx,
		server.WithApplicationAddress(
			cmd.Env("APPLICATION_ADDR", "127.0.0.1:8080"),
		),
		server.WithDatabase(
			cmd.Env("DB_USERNAME", "todo"),
			cmd.Env("DB_PASSWORD", "todo"),
			cmd.Env("DB_ADDR", "localhost"),
			cmd.Env("DB_PORT", "5432"),
			cmd.Env("DB_NAME", "todo"),
			cmd.Env("DB_SSL", "disable"),
		),
		server.WithLogLevel(
			cmd.Env("LOGLEVEL", "debug"),
		),
	)
}
