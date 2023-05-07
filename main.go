package main

import (
	"fmt"
	b64 "encoding/base64"
	"io"
	"log"
	"os"

	"github.com/gliderlabs/ssh"
	"github.com/akamensky/argparse"
)

const BANNER = `
ICAgICAgICAgICAgICAgICAgICAgX19fXyAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgX18KICAgX19fXyBfX19fXyAgICAgICAvIF9fL19fXyAgX19fX19fICAgICAgX19fX19fIF9f
X19fX19fX18vIC8KICAvIF9fIGAvIF9fIFwgICAgIC8gL18vIF9fIFwvIF9fXy8gfCAvfCAvIC8g
X18gYC8gX19fLyBfXyAgLyAKIC8gL18vIC8gL18vIC8gICAgLyBfXy8gL18vIC8gLyAgIHwgfC8g
fC8gLyAvXy8gLyAvICAvIC9fLyAvICAKIFxfXywgL1xfX19fL19fX18vXy8gIFxfX19fL18vICAg
IHxfXy98X18vXF9fLF8vXy8gICBcX18sXy8gICAKL19fX18vICAgICAvX19fX18vICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAK`

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
			Help: "The port to run the SSH server on (Optional)",
			Default: "22",
		})

	iface := parser.String("i", "interface",
		&argparse.Options{
			Required: false,
			Help: "The interface to run the SSH server on (Optional)",
			Default: "0.0.0.0",
		})

	userName := parser.String("U", "username",
	&argparse.Options{
		Required: true,
		Help: "The username for authentication to the SSH server (Required)",
	})

	passWord := parser.String("P", "password",
	&argparse.Options{
		Required: true,
		Help: "The password for authentication to the SSH server (Required)",
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