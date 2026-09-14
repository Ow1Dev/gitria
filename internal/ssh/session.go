package server

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh"
)

type GitService interface {
	UploadPack(repo string, channel io.ReadWriter) error
	ReceivePack(repo string, channel io.ReadWriter) error
	UploadArchive(repo string, channel io.ReadWriter) error
}

func handleSession( channel ssh.Channel, requests <-chan *ssh.Request, git GitService, logger *zerolog.Logger) {
		defer channel.Close()

		var requestWG sync.WaitGroup

		requestWG.Go(func() {
				for req := range requests {
					switch req.Type {
						case "pty-req": _ = req.Reply(true, nil)
						case "exec":        handleExec(channel, req, git, logger)
						default:
						_ = req.Reply(false, nil)
					}
				}
		})

		requestWG.Wait()
	
}

func handleExec(channel ssh.Channel, req *ssh.Request, git GitService, logger *zerolog.Logger) {
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

  cmd, repo, err := parseGitCommand(payload.Command)
	if err != nil {
  	logger.Warn().
			Err(err).
			Str("command", payload.Command).
			Msg("invalid git command")

			_ = req.Reply(false, nil)
			return
	}

  switch cmd {
    case "git-upload-pack":
        // Handle clone/fetch
				err = git.UploadPack(repo, channel)

		default: 
    		logger.Warn().
            Str("command", cmd).
            Msg("unsupported git command")
 				_ = req.Reply(false, nil)
        return
	}

	status := uint32(0)
	if err != nil {
			logger.Error().
					Err(err).
					Str("repo", repo).
					Msg("git command failed")

			return
	}

	_, _ = channel.SendRequest(
		"exit-status",
		false,
		ssh.Marshal(struct {
			Status uint32
		}{status}),
	)

	_ = channel.Close()
}


func parseGitCommand(command string) (string, string, error) {
	args := strings.Split(command, " ")
	if len(args) != 2 {
		return "", "", fmt.Errorf("invalid git command: expected exactly 2 arguments, got %d", len(args))
	}

	return args[0], args[1], nil
}
