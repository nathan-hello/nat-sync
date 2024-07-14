package server

import (
	"bufio"
	"fmt"
	"net"

	"github.com/nathan-hello/nat-sync/src/commands"
)

type ServerParams struct {
	ServerAddress string
	Init          chan bool
	ToServer      chan commands.Command
}

func CreateServer(p ServerParams) {
	listener, err := net.Listen("tcp", p.ServerAddress)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	p.Init <- true
	fmt.Println("Started server at " + p.ServerAddress)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
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
		message, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Connection closed")
			return
		}

		// fmt.Printf("server got msg: %v\n", message)

		dec, err := commands.DecodeCommand(message)
		if err != nil {
			fmt.Println("err: ", err)
		}
		// fmt.Printf("accepted cmd: %v\n", dec)

		p.ToServer <- *dec
		// fmt.Printf("server: sent decoded cmd to channel\n")

	}
}

func transmit(conn net.Conn, p ServerParams) {
	for cmd := range p.ToServer {
		var response []byte
		if cmd.Sub.IsEchoed() {
			r, err := commands.EncodeCommand(&cmd) // \n is put at end here
			if err != nil {
				fmt.Println("Error encoding command:", err)
				continue
			}
			response = r
		} else {
			response = []byte("200\n")
		}

		// fmt.Printf("Sending bits: %b\n", response)
		conn.Write(response)
	}
}
