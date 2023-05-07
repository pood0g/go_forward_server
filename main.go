package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gliderlabs/ssh"
	"github.com/akamensky/argparse"
)

type Arguments struct {
	port string
	iface string
	userName string
	passWord string
}

func get_args() Arguments {
	parser := argparse.NewParser("go_forward_server", "A simple SSH server allowing reverse port forwarding.")

	port := parser.String("p", "port",
		&argparse.Options{
			Required: false,
			Help: "The port to run the SSH server on",
			Default: "22",
		})

	iface := parser.String("i", "interface",
		&argparse.Options{
			Required: false,
			Help: "The interface to run the SSH server on",
			Default: "0.0.0.0",
		})

	userName := parser.String("U", "username",
	&argparse.Options{
		Required: true,
		Help: "The username for authentication to the SSH server",
	})

	passWord := parser.String("P", "password",
	&argparse.Options{
		Required: true,
		Help: "The password for authentication to the SSH server",
	})

	argErr := parser.Parse(os.Args)

	if argErr != nil {
		log.Fatal(parser.Usage(argErr))
	}
	return Arguments{
		port: *port,
		iface: *iface,
		userName: *userName,
		passWord: *passWord,
	}
}

func main() {

	arguments := get_args()

	log.Printf("starting SSH server on port %s:%s...", arguments.iface, arguments.port)

	forwardHandler := &ssh.ForwardedTCPHandler{}

	server := ssh.Server{
		Addr: fmt.Sprintf("%s:%s", arguments.iface, arguments.port),
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			if ctx.User() == arguments.userName	&& password == arguments.passWord {
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