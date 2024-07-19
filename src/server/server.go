package server

import (
	"bufio"
	"net"

	"github.com/nathan-hello/nat-sync/src/messages"
	"github.com/nathan-hello/nat-sync/src/messages/ack"
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
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.ErrorLogger.Println("accepting connection:", err)
			continue
		}
		go receive(conn)
	}
}

func receive(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := messages.WaitReader(reader)
		if err != nil {
			utils.ErrorLogger.Printf("server got a bad message. error: %s\n", err)
		}
		go echo(conn, msg)
	}
}

func echo(conn net.Conn, v messages.Message) {
	var response []byte
	switch msg := v.(type) {
	case *commands.Command:
		if msg.Sub.IsEchoed() {
			r, err := msg.ToBits()
			if err != nil {
				utils.ErrorLogger.Printf("encoding command. cmd: %#v\n err:%s", msg, err)
			}
			response = r
		}
	case *ack.Ack:
		ack, err := ack.New("200")
		if err != nil {
			utils.ErrorLogger.Printf("creating new ack message. err: %s", err)
		}
		r, err := ack.ToBits()
		if err != nil {
			utils.ErrorLogger.Printf("encoding ack. err: %s", err)
		}
		response = r
	}
	utils.DebugLogger.Printf("Sending bits: %b\n", response)
	conn.Write(response)
}
