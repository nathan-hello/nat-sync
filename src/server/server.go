package server

import (
	"bufio"
	"io"
	"net"

	"github.com/nathan-hello/nat-sync/src/messages"
	"github.com/nathan-hello/nat-sync/src/messages/commands"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type ServerParams struct {
	ServerAddress string
}

func CreateServer(p *ServerParams) error {
	listener, err := net.Listen("tcp", p.ServerAddress)
	if err != nil {
		return err
	}

	go listen(listener)
	utils.DebugLogger.Println("started server at " + p.ServerAddress)

	return nil
}

func listen(listener net.Listener) {
	msgChan := make(chan messages.Message)
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.ErrorLogger.Println("accepting connection:", err)
			continue
		}
		go receive(conn, msgChan)
		go transmit(conn, msgChan)
	}
}

func receive(conn net.Conn, msgChan chan messages.Message) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msgs, err := messages.WaitReader(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			utils.ErrorLogger.Printf("server got a bad message. error: %s\n", err)
		}
		for _, v := range msgs {
			msgChan <- v
		}
	}
}

func transmit(conn net.Conn, msgChan chan messages.Message) {
	for v := range msgChan {
		var response []byte
		var err error
		switch msg := v.(type) {
		case *commands.Command:
			response, err = msg.Sub.ExecuteServer()
			if err != nil {
				utils.ErrorLogger.Printf("running cmd on server failed. cmd: %#v\n err:%s", msg, err)
			}
			if msg.Sub.IsEchoed() {
				response, err = msg.ToBits()
				if err != nil {
					utils.ErrorLogger.Printf("encoding command. cmd: %#v\n err:%s", msg, err)
				}
			}

			if len(response) > 0 {
				utils.DebugLogger.Printf("Sending bits: %b\n", response)
				conn.Write(response)
			}
		default:
			utils.ErrorLogger.Printf("server got a non-command message: %#v\n", msg)

		}
	}
}
