package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nathan-hello/nat-sync/src/client"
	"github.com/nathan-hello/nat-sync/src/client/players"
	"github.com/nathan-hello/nat-sync/src/db"
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
	_, err := utils.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}

	err = db.DbInit()
	if err != nil {
		log.Fatal(err)
	}

	utils.InitLogger()

	serverRooms := map[int64]utils.ServerRoom{
		1: {Id: 1, Name: "asdf", Password: "", Clients: make(map[int64]utils.Client)},
	}

	err = server.CreateServer(&server.ServerParams{ServerAddress: serverAddr, Rooms: serverRooms})
	if err != nil {
		utils.ErrorLogger.Fatalf("server could not be started. err: %s", err)
	}

	player := players.New(utils.TargetMpv)
	err = player.Launch()
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
