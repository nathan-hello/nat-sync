package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-hello/nat-sync/src/client"
	"github.com/nathan-hello/nat-sync/src/commands"
	"github.com/nathan-hello/nat-sync/src/server"
)

const (
	clientAddr = ":1412"
	serverAddr = ":1412"
)

var (
	signalListener = make(chan os.Signal, 1) // this only works on unix
	cleanupTime    = make(chan bool, 1)
	clientInit     = make(chan bool, 1)
	serverInit     = make(chan bool, 1)
	toClientCmds   = make(chan commands.Command, 5)
	toServerCmds   = make(chan commands.Command, 5)
)

func main() {
	args := flag.NewFlagSet("natsync", flag.ExitOnError)
	_ = args.String("client", "4000", "Create a client to join natsync servers")
	_ = args.String("server", "4000", "Become a server for natsync clients")

	serverParams := server.ServerParams{
		ServerAddress: serverAddr,
		Init:          serverInit,
		ToServer:      toServerCmds,
	}
	go server.CreateServer(serverParams)
	<-serverInit

	clientParams := client.ClientParams{
		ServerAddress: serverAddr,
		Init:          clientInit,
		ToClient:      toClientCmds,
	}
	go client.CreateClient(clientParams)
	<-clientInit

	signal.Notify(signalListener, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)
	signal := <-signalListener
	fmt.Printf("\nReceived signal %s, exiting\n", signal)
	cleanupTime <- true
}
