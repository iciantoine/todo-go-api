package cmd

import (
	"context"
	"os/signal"
	"syscall"
)

// Run runs the given function with a context that is closed as soon as an OS
// signal is caught.
func Run(f func(context.Context) error) error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	return f(ctx)
}
