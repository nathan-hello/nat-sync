package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"slices"

	"github.com/nathan-hello/nat-sync/src/commands"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type ClientParams struct {
	ServerAddress string
	Init          chan bool
	ToClient      chan commands.Command
}

func CreateClient(p ClientParams) {
	conn, err := net.Dial("tcp", p.ServerAddress)
	if err != nil {
		utils.ErrorLogger.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	utils.DebugLogger.Println("Started client connected to " + p.ServerAddress)

	lp := LaunchParams{
		Player:     LaunchMpv,
		SocketPath: "/tmp/nat-sync-mpv-socket",
		Init:       make(chan bool),
		ToClient:   p.ToClient,
	}
	go LaunchPlayer(&lp)
	<-lp.Init

	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadBytes('\n')
			if err != nil {
				utils.ErrorLogger.Println("Connection closed")
				return
			}
			// utils.DebugLogger.Printf("Received from server: %b\n", message)
			if slices.Equal(message, []byte("200")) {
				utils.DebugLogger.Printf("Received OK server: %s\n", message)
				continue
			}

			dec, err := commands.DecodeCommand(message)
			if err != nil {
				utils.ErrorLogger.Println("err: ", err)
			}
			// utils.DebugLogger.Printf("%#v\n", dec)
			go func() {
				utils.DebugLogger.Printf("sending cmd to ToClient: %#v\n", dec)
				p.ToClient <- *dec
			}()
		}
	}()
	p.Init <- true

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		creator := uint16(10126)

		cmd, err := commands.CmdFromString(text)
		if err != nil {
			utils.ErrorLogger.Println(err)
			continue
		}

		cmd.UserId = creator

		bits, err := cmd.ToBits()
		if err != nil {
			utils.ErrorLogger.Println("err in cmd.ToBits() in client transmit. err: ", err)
			continue
		}
		tosend := string(bits)
		fmt.Fprintln(conn, tosend)
	}
}
