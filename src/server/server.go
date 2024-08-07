package server

import (
	"bufio"
	"context"
	"database/sql"
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
	Rooms         map[int64]utils.ServerRoom
}

func CreateServer(p *ServerParams) error {
	listener, err := net.Listen("tcp", p.ServerAddress)
	if err != nil {
		return err
	}

	man := &Manager{
		Lock:  sync.Mutex{},
		Rooms: p.Rooms,
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
			response, err = msg.Execute(nil)
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
				utils.DebugLogger.Printf("Sending bits: %b\tstruct: %#v\n", response, v)
				man.BroadcastMessage(v.RoomId, response)
			}
		case messages.AdminCommand:
			switch admin := msg.(type) {
			case *impl.Join:
				var user db.User
				var err error
				user, err = d.SelectUserByName(context.Background(), admin.Username)

				if err != nil && err != sql.ErrNoRows {
					utils.DebugLogger.Printf("db select user: %s", err)
					continue
				}

				if err == sql.ErrNoRows {
					user, err = d.InsertUser(context.Background(), admin.Username)
					if err != nil {
						utils.ErrorLogger.Printf("db insert user: %s", err)
						continue
					}
				}

				var room db.SelectRoomByNameRow
				room, err = d.SelectRoomByName(context.Background(), admin.RoomName)
				if err != nil {
					if err == sql.ErrNoRows {
						err = d.InsertRoom(context.Background(), db.InsertRoomParams{Name: admin.RoomName})
						if err != nil {
							utils.ErrorLogger.Printf("could not add requested room. request: %#v\n database err: %s", admin, err)
							continue
						}
						room, err = d.SelectRoomByName(context.Background(), admin.RoomName)
						if err != nil {
							utils.ErrorLogger.Printf("could not select requested room that was just made. request: %#v\n database err: %s", admin, err)
							continue
						}
					} else {
						utils.ErrorLogger.Printf("could not select requested room. request: %#v\n database err: %s", admin, err)
						continue

					}

					utils.DebugLogger.Printf("adding client name %s id %d", user.Username, user.ID)
					man.AddClient(room.ID, utils.Client{Id: user.ID, Name: user.Username, Conn: conn})

					accept, err := messages.New(&impl.Accept{Action: impl.AcceptHead.Ok, RoomId: room.ID}, &v.RoomId)
					if err != nil {
						utils.ErrorLogger.Printf("could not respond to room request: %#v\n database err: %s", admin, err)
					}
					respBits, _ := accept[0].ToBits()
					conn.Write(respBits)

					change, err := messages.New(&impl.Change{Action: impl.ChgImmediate}, &accept[0].RoomId)

				}
			default:
				utils.ErrorLogger.Printf("server got a non-command message: %#v\n", msg)

			}
		}
	}
}
