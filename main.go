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
)

func main() {
	args := flag.NewFlagSet("natsync", flag.ExitOnError)
	_ = args.String("client", "4000", "Create a client to join natsync servers")
	_ = args.String("server", "4000", "Become a server for natsync clients")

	utils.InitLogger()

	err := server.CreateServer(&server.ServerParams{ServerAddress: serverAddr})
	if err != nil {
		utils.ErrorLogger.Fatalf("server could not be started. err: %s", err)
	}

	player, err := players.New(utils.TargetMpv)

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
	fmt.Printf("\nReceived signal %s, exiting\n", <-signalListener)
}
