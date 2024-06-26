package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-hello/nat-sync/src"
)

const (
	clientAddr = ":1412"
	serverAddr = ":1412"
)

var (
	signalListener = make(chan os.Signal, 1) // this only works on unix
	cleanupTime    = make(chan bool, 1)
	clientCmdQueue chan src.Command
	clientInit     = make(chan bool)
	serverInit     = make(chan bool)
)

func main() {
	args := flag.NewFlagSet("natsync", flag.ExitOnError)
	_ = args.String("client", "4000", "Create a client to join natsync servers")
	_ = args.String("server", "4000", "Become a server for natsync clients")

	go src.CreateServer(serverAddr, serverInit)
	<-serverInit
	go src.CreateClient(clientCmdQueue, clientAddr, clientInit)
	<-clientInit

	signal.Notify(signalListener, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)
	signal := <-signalListener
	fmt.Printf("\nReceived signal %s, exiting\n", signal)
	cleanupTime <- true
}
