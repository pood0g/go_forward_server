package main

import (
	"os"
	"log"

	"github.com/akamensky/argparse"
)

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