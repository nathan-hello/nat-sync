package client

import (
	"bufio"
	"io"
	"net"

	"github.com/nathan-hello/nat-sync/src/client/players"
	"github.com/nathan-hello/nat-sync/src/commands"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type ClientParams struct {
	ServerAddress string
	Init          chan bool
	ToClient      chan commands.Command
	InputReader   io.Reader
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

	go transmit(conn, p)
	p.Init <- true

	receive(conn, p)

}

func receive(conn net.Conn, p *ClientParams) {
	reader := bufio.NewReader(conn)
	for {
		cmdBits, err := utils.ConnRxCommand(reader)
		if err != nil {
			utils.ErrorLogger.Println(err)
		}
		dec, err := commands.DecodeCommand(cmdBits)
		if err != nil {
			utils.ErrorLogger.Println(err)
		}

		p.ToClient <- *dec
		utils.DebugLogger.Printf("server: sent decoded cmd to channel\n")
	}
}

func transmit(conn net.Conn, p *ClientParams) {
	scanner := bufio.NewScanner(p.InputReader)
	for scanner.Scan() { // this blocks the terminal
		text := scanner.Text()
		creator := uint16(10126)

		utils.DebugLogger.Printf("new reader text: %s\n", text)

		cmd, err := commands.CmdFromString(text)
		if err != nil {
			utils.ErrorLogger.Println(err)
			continue
		}

		cmd.UserId = creator

		bits, err := cmd.ToBits()
		if err != nil {
			utils.ErrorLogger.Println("cmd.ToBits() in client transmit. err: ", err)
			continue
		}
		utils.DebugLogger.Printf("client sending bits: %b, length: %d\n", bits, len(bits))
		_, err = conn.Write(bits)
		if err != nil {
			utils.ErrorLogger.Printf("client writing bits failed. bits: %b\n", bits)
		}
	}
}

func launchPlayer(p *players.LaunchParams) error {
	switch p.Player {
	case players.LaunchMpv:
		mpv := players.NewMpv()
		cmd := mpv.GetLaunchCmd()

		err := cmd.Start()
		if err != nil {
			return err
		}

		go mpv.Connect(p)
		cmd.Wait()
	case players.LaunchVlc:
		p.Init <- true // init is usually handled after connection to player is successful
		return utils.ErrNotImplemented("vlc")
	case players.NoLaunchy:
		p.Init <- true // init is usually handled after connection to player is successful
		return nil
	}

	return nil

}
