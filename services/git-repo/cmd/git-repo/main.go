package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/Ow1Dev/gitria.git/gen/git/v0"
	"github.com/Ow1Dev/gitria.git/services/git-repo/internal/config"
	"github.com/Ow1Dev/gitria.git/services/git-repo/internal/server"
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

	grpcLis, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterGitServiceServer(grpcSrv, server.NewServer())
	wg.Go(func() {
		logger.Info().Str("address", cfg.GRPCAddress).Msg("grpc listening")
		if err := grpcSrv.Serve(grpcLis); err != nil {
			logger.Error().Err(err).Msg("grpc server error")
		}
	})

	<-ctx.Done()
	logger.Info().Msg("shutting down")
	grpcSrv.GracefulStop()

	wg.Wait()
	return nil
}
