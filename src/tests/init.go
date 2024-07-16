package tests

import (
	"io"
	"os"

	"github.com/nathan-hello/nat-sync/src/client"
	"github.com/nathan-hello/nat-sync/src/client/players"
	"github.com/nathan-hello/nat-sync/src/commands"
	"github.com/nathan-hello/nat-sync/src/server"
	"github.com/nathan-hello/nat-sync/src/utils"
)

const (
	clientAddr = ":1413"
	serverAddr = ":1413"
)

var (
	signalListener = make(chan os.Signal, 1) // this only works on unix
	cleanupTime    = make(chan bool, 1)
	clientInit     = make(chan bool, 1)
	serverInit     = make(chan bool, 1)
	playerInit     = make(chan bool, 1)
	toClientCmds   = make(chan commands.Command, 5)
	toServerCmds   = make(chan commands.Command, 5)
)

func initEnvironment() {
	utils.InitLogger()
}

func initClient(r io.Reader) {
	clientParams := client.ClientParams{
		ServerAddress: serverAddr,
		Init:          clientInit,
		ToClient:      toClientCmds,
		InputReader:   r,
	}

	lp := players.LaunchParams{
		Player:   players.NoLaunchy,
		Init:     playerInit,
		ToClient: toClientCmds,
	}
	go client.CreateClient(&clientParams, &lp)
}

func initServer() {
	serverParams := server.ServerParams{
		ServerAddress: serverAddr,
		Init:          serverInit,
		ToServer:      toServerCmds,
	}
	go server.CreateServer(serverParams)
}

func initAll(r io.Reader) {
	initEnvironment()
	go initServer()
	<-serverInit
	go initClient(r)
	<-clientInit
}
