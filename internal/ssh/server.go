package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/Ow1Dev/gitria-git/internal/config"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	address 	string
	sshConfig *ssh.ServerConfig
	logger 		*zerolog.Logger
	listener  net.Listener 

	git GitService

	mu    sync.Mutex
	conns map[net.Conn]struct{}
}

func New(cfg config.SSHConfig, git GitService, logger *zerolog.Logger) (*Server, error) {
	privateBytes, err := os.ReadFile(cfg.HostKeyPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load private key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse private key")
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	config.AddHostKey(private)
	return &Server{
		address:   fmt.Sprintf(":%s", cfg.PORT),
		sshConfig: config,
		logger: 	logger,
		git: git,

		conns:   make(map[net.Conn]struct{}),
	}, nil
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("failed to listen for connection")
	}
	s.listener = listener

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}

			s.logger.Error().Err(err).Msg("accept error")
			continue
		}

		s.mu.Lock()
		s.conns[conn] = struct{}{}
		s.mu.Unlock()

		wg.Go(func() {
			defer func() {
				s.mu.Lock()
				delete(s.conns, conn)
				s.mu.Unlock()
			}()

			s.serveConnection(conn)
		})
	}

	wg.Wait()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	err := s.listener.Close()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.conns {
		_ = conn.Close()
	}

	return nil
}
