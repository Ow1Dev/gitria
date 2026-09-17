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
	"github.com/Ow1Dev/gitria-git/internal/db"
	"github.com/Ow1Dev/gitria-git/internal/git"
	"github.com/Ow1Dev/gitria-git/internal/httpserver"
	"github.com/Ow1Dev/gitria-git/internal/sqliteStore"
	server "github.com/Ow1Dev/gitria-git/internal/ssh"
	"github.com/Ow1Dev/gitria-git/migrations"
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

	err = createDataFolder(cfg.Data, logger)
	if err != nil {
		return fmt.Errorf("Could not create data path: %w", err)
	}

	logger.Info().Msgf("Connecting to DB") 
	conn, err := db.Open(cfg.GetDbFilePath())
	if err != nil {
		return fmt.Errorf("Could not connect to DB: %w", err)
	}

	migrations.RunMigration(conn, logger)
	queries := sqliteStore.New(conn);

	gitsrv, err := git.New(cfg.GetRepoFolderPath())
	if err != nil {
		return fmt.Errorf("Could new create git repo dir: %w", err)
	}

	sshsrv, err := server.New(cfg.Ssh, gitsrv, &logger)
	if err != nil {
		return fmt.Errorf("server: %w", err)
	}

	handler := httpserver.NewHttpServer(conn, *queries, logger)
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
	logger.Info().Msg("shutting http")
	if err := httpsrv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
		if err := httpsrv.Close(); err != nil {
			logger.Error().Err(err).Msg("force close error")
		}
	}

	logger.Info().Msg("shutting ssh")
	if err := sshsrv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("graceful shutdown failed")
	}

	wg.Wait()
	return nil
}

func createDataFolder(dataPath string, logger zerolog.Logger) error {
	info, err := os.Stat(dataPath)

	switch {
	case err == nil && info.IsDir():
		logger.Info().Msgf("Data directory exists: %s", dataPath)

	case os.IsNotExist(err):
		logger.Info().Msgf("Data directory does not exist, creating: %s", dataPath)

	case err != nil:
		return fmt.Errorf("could not check data path: %w", err)

	default:
		return fmt.Errorf("data path exists but is not a directory: %s", dataPath)
	}

	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return fmt.Errorf("could not create data path: %w", err)
	}

	return nil
}
