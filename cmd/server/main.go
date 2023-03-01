package main

import (
	"context"
	"os"

	"github.com/iciantoine/todo-go-api/cmd"
)

func main() {
	if err := cmd.Run(run); err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	return nil
}
