package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Ow1Dev/gitria-git/internal/config"
	"github.com/Ow1Dev/gitria-git/internal/git"
	"github.com/Ow1Dev/gitria-git/internal/httpserver"
	server "github.com/Ow1Dev/gitria-git/internal/ssh"
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
	var wg sync.WaitGroup

	handler := httpserver.NewHttpServer(logger)
	gitsrv := git.New()

	sshsrv, err := server.New(cfg.Ssh, gitsrv, &logger)
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	httpsrv := &http.Server{
		Addr:         cfg.HTTP.Address,
		Handler:      handler,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	wg.Go(func() {
		logger.Info().Str("address", cfg.HTTP.Address).Msg("http listening")
		if err := httpsrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("http server error")
		}
	})

	wg.Go(func() {
		logger.Info().Str("address", cfg.Ssh.PORT).Msg("ssh listening")
		if err := sshsrv.ListenAndServe(); err != nil {
			logger.Error().Err(err).Msg("http server error")
		}
	})

	<-ctx.Done()
	logger.Info().Msg("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpsrv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
		if err := httpsrv.Close(); err != nil {
			logger.Error().Err(err).Msg("force close error")
		}
	}
	if err := sshsrv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
	}

	wg.Wait()
	return nil
}
