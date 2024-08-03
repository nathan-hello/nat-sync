package server

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"

	"github.com/nathan-hello/nat-sync/src/db"
	"github.com/nathan-hello/nat-sync/src/messages"
	"github.com/nathan-hello/nat-sync/src/messages/impl"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type ServerParams struct {
	ServerAddress string
	Rooms         []Room
}

func CreateServer(p *ServerParams) error {
	listener, err := net.Listen("tcp", p.ServerAddress)
	if err != nil {
		return err
	}

	man := &Manager{
		Lock:    sync.Mutex{},
		Clients: []Client{},
	}

	go listen(listener, man)
	utils.DebugLogger.Println("started server at " + p.ServerAddress)

	return nil
}

func listen(listener net.Listener, man *Manager) {
	msgChan := make(chan messages.Message)
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.ErrorLogger.Println("accepting connection:", err)
			continue
		}
		go receive(conn, msgChan)
		go handle(conn, msgChan, man)
		whoMsg, _ := messages.New("who name=server isresponse=false")
		msgChan <- whoMsg[0]

	}
}

func receive(conn net.Conn, msgChan chan messages.Message) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msgs, err := messages.WaitReader(reader)
		if err == io.EOF {
			close(msgChan)
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

func handle(conn net.Conn, msgChan chan messages.Message, man *Manager) {
	d := db.Db()
	for v := range msgChan {
		var response []byte
		var err error
		switch msg := v.Sub.(type) {
		case messages.ServerCommand:
			response, err = msg.Execute()
			if err != nil {
				utils.ErrorLogger.Printf("running cmd on server failed. cmd: %#v\n err:%s", msg, err)
			}
			utils.DebugLogger.Printf("server executed sub %#v\nresponse: %s\n", msg, response)
		case messages.PlayerCommand:
			response, err = v.ToBits()
			if err != nil {
				utils.ErrorLogger.Printf("encoding command. cmd: %#v\n err:%s", msg, err)
			}

			if len(response) > 0 {
				// utils.DebugLogger.Printf("Sending bits: %b\n", response)
				man.BroadcastMessage(response)
			}
		case messages.AdminCommand:
			switch admin := msg.(type) {
			case *impl.Connect:
				user, err := d.InsertUser(context.Background(), db.InsertUserParams{Username: admin.Name})
				if err != nil {
					utils.ErrorLogger.Printf("db insert user: %s", err)
					ack, _ := messages.New(impl.Ack{Code: impl.AckCode.InternalServiceError, Message: "server had a database error"})
					ackBits, _ := ack[0].ToBits()
					conn.Write(ackBits)
					continue
				}

				man.AddClient(Client{
					Id:     user.ID,
					Name:   user.Username,
					Conn:   conn,
					RoomId: admin.RoomId,
				})

				ack, _ := messages.New(impl.Ack{Code: impl.AckCode.Ok, Message: "connection accepted"})
				ackBits, _ := ack[0].ToBits()
				conn.Write(ackBits)
			}
		default:
			utils.ErrorLogger.Printf("server got a non-command message: %#v\n", msg)

		}
	}
}
