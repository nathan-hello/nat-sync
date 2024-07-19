package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-hello/nat-sync/src/client"
	"github.com/nathan-hello/nat-sync/src/players"
	"github.com/nathan-hello/nat-sync/src/server"
	"github.com/nathan-hello/nat-sync/src/utils"
)

const (
	serverAddr = ":1412"
)

var (
	signalListener = make(chan os.Signal, 1) // this only works on unix
	cleanupTime    = make(chan bool, 1)
)

func main() {
	args := flag.NewFlagSet("natsync", flag.ExitOnError)
	_ = args.String("client", "4000", "Create a client to join natsync servers")
	_ = args.String("server", "4000", "Become a server for natsync clients")

	utils.InitLogger()

	server.CreateServer(&server.ServerParams{ServerAddress: serverAddr})

	player, err := players.New(players.Mpv)

	if err != nil {
		log.Fatalf("could not start video player for reason: %s\n", err)
	}

	client.CreateClient(
		&client.ClientParams{
			ServerAddress: serverAddr,
			InputReader:   os.Stdin,
			Player:        player,
		},
	)

	signal.Notify(signalListener, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)
	signal := <-signalListener
	fmt.Printf("\nReceived signal %s, exiting\n", signal)
	cleanupTime <- true
}
