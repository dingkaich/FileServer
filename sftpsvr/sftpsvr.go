package sftpsvr

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var username = "testuser"
var userauth = "testuser"

func authSftp(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	// Should use constant-time compare (or better, salt+hash) in
	// a production setting.
	log.Printf("Login: %s\n", c.User())

	err := AuthUsrPwd(c.User(), string(pass))
	if err != nil {
		log.Println("sftp auth error=", err)
		return nil, fmt.Errorf("password rejected for %q", c.User())
	}
	return nil, nil

}

func SftpServer(addr string) {

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		PasswordCallback: authSftp,
	}

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key", err)
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	//	// accepted.
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to listen for connection", err)
	}
	log.Printf("Listening on %v\n", listener.Addr())

	for {

		nConn, err := listener.Accept()
		if err != nil {
			log.Fatal("failed to accept incoming connection", err)
		}

		go func() {
			// Before use, a handshake must be performed on the incoming
			// net.Conn.
			_, chans, reqs, err := ssh.NewServerConn(nConn, config)
			if err != nil {
				log.Fatal("failed to handshake", err)
			}
			log.Println("SSH server established\n")

			// The incoming Request channel must be serviced.
			go ssh.DiscardRequests(reqs)

			// Service the incoming Channel channel.
			for newChannel := range chans {
				// Channels have a type, depending on the application level
				// protocol intended. In the case of an SFTP session, this is "subsystem"
				// with a payload string of "<length=4>sftp"
				log.Printf("Incoming channel: %s\n", newChannel.ChannelType())

				if newChannel.ChannelType() != "session" {
					newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
					log.Printf("Unknown channel type: %s\n", newChannel.ChannelType())
					continue
				}
				channel, requests, err := newChannel.Accept()
				if err != nil {
					log.Println("could not accept channel.", err)
					break
				}
				log.Println("Channel accepted\n")

				// Sessions have out-of-band requests such as "shell",
				// "pty-req" and "env".  Here we handle only the
				// "subsystem" request.
				go func(in <-chan *ssh.Request) {
					for req := range in {
						log.Printf("Request: %v\n", req.Type)
						ok := false
						switch req.Type {
						case "subsystem":
							log.Printf("Subsystem: %s\n", req.Payload[4:])
							if string(req.Payload[4:]) == "sftp" {
								ok = true
							}
						}
						log.Printf(" - accepted: %v\n", ok)
						req.Reply(ok, nil)
					}
				}(requests)

				server, err := sftp.NewServer(
					channel,
				)
				if err != nil {
					log.Fatal(err)
				}
				if err := server.Serve(); err == io.EOF {
					server.Close()
					log.Print("sftp client exited session.")
				} else if err != nil {
					log.Fatal("sftp server completed with error:", err)
				}
			}

		}()
	}

}
