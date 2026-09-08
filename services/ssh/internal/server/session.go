package server

import (
	"sync"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh"
)

func handleSession( channel ssh.Channel, requests <-chan *ssh.Request, logger *zerolog.Logger) {
		defer channel.Close()

		var requestWG sync.WaitGroup

		requestWG.Go(func() {
				for req := range requests {
					switch req.Type {
						case "pty-req": _ = req.Reply(true, nil)
						case "exec":        handleExec(channel, req, logger)
						default:
						_ = req.Reply(false, nil)
					}
				}
		})

		requestWG.Wait()
	
}

func handleExec(_ ssh.Channel, req *ssh.Request, logger *zerolog.Logger) {
	var payload struct {
		Command string
	}

	if err := ssh.Unmarshal(req.Payload, &payload); err != nil {
		_ = req.Reply(false, nil)
		return
	}

	logger.Info().
		Str("command", payload.Command).
		Msg("exec request")

	_ = req.Reply(true, nil)
}
