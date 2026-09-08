package server

import (
	"log"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

func (s Server) serveConnection(conn net.Conn) {
	defer conn.Close()

	s.logger.Info().
		Str("remote", conn.RemoteAddr().String()).
		Msg("new connection")

	sshConn, channels, requests, err := ssh.NewServerConn(conn, s.sshConfig)
	if err != nil {
		s.logger.Error().Err(err).Str("remote", conn.RemoteAddr().String()).Msg("ssh handshake failed")
		return
	}
	defer sshConn.Close()

	go ssh.DiscardRequests(requests)

	var wg sync.WaitGroup

	for newChannel := range channels {
		if newChannel.ChannelType() != "session" {
			_ = newChannel.Reject(ssh.UnknownChannelType, "only session channels are supported",)
		  continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to accept channel")
			continue
		}

		wg.Go(func() {
			handleSession(channel, requests, s.logger)
		})
	}

	wg.Wait()
	log.Printf("SSH connection closed from %s", sshConn.RemoteAddr())
}
