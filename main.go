package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/gliderlabs/ssh"
)

type Arguments struct {
	port string
	iface string
}

func get_args() Arguments {
	parser := argparse.NewParser("go_forward_server", "A simple SSH server allowing reverse port forwarding.")

	port := parser.String("p", "port",
		&argparse.Options{
			Required: false,
			Help: "The port to run the SSH server on",
			Default: "22",
		})

	iface := parser.String("i", "Interface to listen on",
		&argparse.Options{
			Required: false,
			Help: "The interface to run the SSH server on",
			Default: "0.0.0.0",
		})

	argErr := parser.Parse(os.Args)

	if argErr != nil {
		log.Fatal(parser.Usage(argErr))
	}
	return Arguments{
		port: *port,
		iface: *iface,
	}
}

func main() {

	arguments := get_args()

	log.Printf("starting SSH server on port %s:%s...", arguments.iface, arguments.port)

	forwardHandler := &ssh.ForwardedTCPHandler{}

	server := ssh.Server{
		Addr: fmt.Sprintf("%s:%s", arguments.iface, arguments.port),
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			if ctx.User() == "test"	 && password == "test1234" {
				log.Printf("user %s connected from %s", ctx.User(), ctx.RemoteAddr())
				return true
			}
			return false
		},
		Handler: ssh.Handler(func(s ssh.Session) {
			io.WriteString(s, "Remote forwarding only...\n")
			select {}
		}),
		ReversePortForwardingCallback: ssh.ReversePortForwardingCallback(func(ctx ssh.Context, host string, port uint32) bool {
			log.Println("attempt to bind", host, port, "granted")
			return true
		}),
		RequestHandlers: map[string]ssh.RequestHandler{
			"tcpip-forward":        forwardHandler.HandleSSHRequest,
			"cancel-tcpip-forward": forwardHandler.HandleSSHRequest,
		},
		Version: "go_forward_server",

	}

	log.Fatal(server.ListenAndServe())
}