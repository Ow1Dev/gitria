package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ow1Dev/gitria.git/services/ssh/internal/config"
	"github.com/Ow1Dev/gitria.git/services/ssh/internal/server"
	"github.com/rs/zerolog"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := run(ctx, os.Args, config.LoadFromEnv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	_ []string,
	loadConfig func(func(string) string) (config.Config, error),
	_ io.Reader,
	writer, _ io.Writer,
) error {
  cfg, err := loadConfig(os.Getenv)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	logger := zerolog.New(writer).With().Timestamp().Logger()

	srv, err := server.New(cfg, &logger)
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	return srv.Run(ctx)
}
