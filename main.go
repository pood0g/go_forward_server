package main

import (
	"fmt"
	b64 "encoding/base64"
	"io"
	"log"

	"github.com/gliderlabs/ssh"
)

func main() {

	banner, _ := b64.StdEncoding.DecodeString(BANNER)
	fmt.Printf("%s\n", banner)
	
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