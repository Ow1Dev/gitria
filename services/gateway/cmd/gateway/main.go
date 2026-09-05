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

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Ow1Dev/gitria.git/gen/git/v0"
	"github.com/Ow1Dev/gitria.git/services/gateway/internal/config"
	"github.com/Ow1Dev/gitria.git/services/gateway/internal/server"
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

	conn, err := grpc.NewClient(cfg.RepoAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to git repo grpc")
	}
	defer conn.Close()
	c := pb.NewGitServiceClient(conn)

	handler := server.New(c, &logger)
	srv := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	wg.Go(func() {
		logger.Info().Str("address", cfg.HTTPAddress).Msg("gateway listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("http server error")
		}
	})

	<-ctx.Done()
	logger.Info().Msg("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
		if err := srv.Close(); err != nil {
			logger.Error().Err(err).Msg("force close error")
		}
	}

	wg.Wait()
	return nil
}
