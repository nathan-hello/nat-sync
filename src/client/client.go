package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/nathan-hello/nat-sync/src/client/players"
	"github.com/nathan-hello/nat-sync/src/commands"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type ClientParams struct {
	ServerAddress string
	Init          chan bool
	ToClient      chan commands.Command
}

func CreateClient(p *ClientParams, lp *players.LaunchParams) {
	conn, err := net.Dial("tcp", p.ServerAddress)
	if err != nil {
		utils.ErrorLogger.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	utils.DebugLogger.Println("Started client connected to " + p.ServerAddress)

	go launchPlayer(lp)
	<-lp.Init

	p.Init <- true
	go receive(conn, p)

	transmit(conn, os.Stdin)

}

func receive(conn net.Conn, p *ClientParams) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			utils.ErrorLogger.Println("Connection closed")
			return
		}

		dec, err := commands.DecodeCommand(message)
		if err != nil {
			utils.ErrorLogger.Println("err: ", err)
		}

		utils.DebugLogger.Printf("sending cmd to ToClient: %#v\n", dec)
		p.ToClient <- *dec
		utils.DebugLogger.Printf("sent cmd to ToClient: %#v\n", dec)
	}
}

func transmit(conn net.Conn, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() { // this blocks the terminal
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

func launchPlayer(p *players.LaunchParams) error {
	if p.Player == "mpv" {
		mpv := players.NewMpv()
		cmd := mpv.GetLaunchCmd()

		err := cmd.Start()
		if err != nil {
			return err
		}

		go mpv.Connect(p)
		cmd.Wait()
	}

	return nil

}
