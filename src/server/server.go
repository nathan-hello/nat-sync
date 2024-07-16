package server

import (
	"bufio"
	"net"

	"github.com/nathan-hello/nat-sync/src/commands"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type ServerParams struct {
	ServerAddress string
	Init          chan bool
	ToServer      chan commands.Command
}

func CreateServer(p ServerParams) {
	listener, err := net.Listen("tcp", p.ServerAddress)
	if err != nil {
		utils.ErrorLogger.Println("starting server:", err)
		return
	}
	defer listener.Close() // the for loop below means this will never run

	p.Init <- true
	utils.DebugLogger.Println("Started server at " + p.ServerAddress)
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.ErrorLogger.Println("accepting connection:", err)
			continue
		}
		go receive(conn, p)
		go transmit(conn, p)
	}
}

func receive(conn net.Conn, p ServerParams) {
	defer conn.Close()

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

		p.ToServer <- *dec
		utils.DebugLogger.Printf("server: sent decoded cmd to channel\n")
	}
}

func transmit(conn net.Conn, p ServerParams) {
	for cmd := range p.ToServer {
		var response []byte
		if cmd.Sub.IsEchoed() {
			r, err := cmd.ToBits() // \n is put at end here
			if err != nil {
				utils.ErrorLogger.Printf("encoding command. cmd: %#v\n err:%s", cmd, err)
				continue
			}
			response = r
		} else {
			response = []byte("200\n")
		}

		utils.DebugLogger.Printf("Sending bits: %b\n", response)
		conn.Write(response)
	}
}
