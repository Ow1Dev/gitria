package main

import (
	"errors"
	"log"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

const (
	listenAddr = ":2222"
	hostKey    = "id_rsa"
)

func main() {
	privateBytes, err := os.ReadFile(hostKey)
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	config.AddHostKey(private)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("failed to listen for connection: ", err)
	}
	defer listener.Close()
	log.Printf("SSH server listening on %s", listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			log.Printf("accept error: %v", err)
			continue
		}

		go serveConnection(conn, config)
	}
}

func serveConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()
	log.Printf("new connection from %s", conn.RemoteAddr())

	sshConn, channels, requests, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("SSH handshake failed from %s: %v",
			conn.RemoteAddr(), err)
		return
	}
	defer sshConn.Close()

	log.Printf("SSH connection established from %s", sshConn.RemoteAddr())
	go ssh.DiscardRequests(requests)

	var wg sync.WaitGroup

	for newChannel := range channels {
		if newChannel.ChannelType() != "session" {
			_ = newChannel.Reject(ssh.UnknownChannelType, "only session channels are supported",)
		  continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("failed to accept channel: %v", err)
			continue
		}

		wg.Go(func() {
			handleSession(channel, requests)
		})
	}

	wg.Wait()
	log.Printf("SSH connection closed from %s", sshConn.RemoteAddr())
}

func handleSession(channel ssh.Channel, requests <-chan *ssh.Request) {
		defer channel.Close()

		var requestWG sync.WaitGroup

		requestWG.Go(func() {
				for req := range requests {
					switch req.Type {
						case "pty-req":
							// Accept PTY allocation.
						_ = req.Reply(true, nil)
						case "exec":
						handleExec(channel, req)
						default:
						_ = req.Reply(false, nil)
					}
				}
		})

		requestWG.Wait()
}


func handleExec(_ ssh.Channel, req *ssh.Request) {

		var payload struct {
			Command string
		}

		if err := ssh.Unmarshal(req.Payload, &payload); err != nil {
			_ = req.Reply(false, nil)
			return
		}
		log.Printf("exec request: %q", payload.Command)

		_ = req.Reply(true, nil)
}



