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
	Rooms         map[int64]utils.ServerRoom
	Manager       *Manager
}

func CreateServer(p *ServerParams) error {
	listener, err := net.Listen("tcp", p.ServerAddress)
	if err != nil {
		return err
	}

	p.Manager = &Manager{
		Lock:  sync.Mutex{},
		Rooms: p.Rooms,
	}

	go listen(listener, p)
	utils.DebugLogger.Println("started server at " + p.ServerAddress)

	return nil
}

func listen(listener net.Listener, p *ServerParams) {
	msgChan := make(chan messages.Message)
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.ErrorLogger.Println("accepting connection:", err)
			continue
		}
		go receive(conn, msgChan)
		go handle(conn, msgChan, p)
	}
}

func receive(conn net.Conn, msgChan chan messages.Message) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := messages.WaitReader(reader)
		if err == io.EOF {
			close(msgChan)
			return
		}
		if err != nil {
			utils.ErrorLogger.Printf("server got a bad message. error: %s\n", err)
		}
		msgChan <- msg
	}
}

func handle(conn net.Conn, msgChan chan messages.Message, p *ServerParams) {
	for v := range msgChan {
		var response []byte
		var err error
		switch msg := v.Sub.(type) {
		case impl.ServerCommand:
			response, err = msg.Execute(nil)
			if err != nil {
				utils.ErrorLogger.Printf("running cmd on server failed. cmd: %#v\n err:%s", msg, err)
			}
			utils.DebugLogger.Printf("server executed sub %#v\nresponse: %s\n", msg, response)
		case impl.PlayerCommand:
			response, err = v.MarshalBinary()
			if err != nil {
				utils.ErrorLogger.Printf("encoding command. cmd: %#v\n err:%s", msg, err)
			}

			if len(response) > 0 {
				utils.DebugLogger.Printf("Sending bits: %b\tstruct: %#v\n", response, v)
				p.Manager.BroadcastMessage(v.RoomId, response)
			}
		case impl.AdminCommand:
			handleAdminMessage(conn, p, v)
		}
	}
}

func handleAdminMessage(conn net.Conn, p *ServerParams, msg messages.Message) {
	d := db.Db()
	switch admin := msg.Sub.(type) {
	case *impl.Join:

		// this is an INSERT OR IGNORE query
		// if it already exists, that's cool
		// because we're about to select it
		room, err := d.InsertRoom(context.Background(), db.InsertRoomParams{Name: admin.RoomName})
		if err != nil {
			return
		}
		utils.DebugLogger.Printf("got room %#v\n", room)

		video, err := d.SelectCurrentVideoByRoomId(context.Background(), room.ID)
		if err != nil {
			return
		}
		utils.DebugLogger.Printf("got video %#v\n", video)

		p.Manager.AddClient(room.ID, utils.Client{Name: admin.Username, Conn: conn})
		utils.DebugLogger.Printf("adding client name %s", admin.Username)

		accept := messages.NewFromSub(&impl.Accept{Action: impl.AcceptHead.Ok, RoomId: room.ID}, msg.RoomId)
		acceptBits, _ := accept.MarshalBinary()
		conn.Write(acceptBits)
		utils.DebugLogger.Printf("sent response %#v\n", accept)

		change := messages.NewFromSub(&impl.Change{Action: impl.ChgImmediate, Uri: video.Uri}, accept.RoomId)
		changeBits, _ := change.MarshalBinary()
		conn.Write(changeBits)
		utils.DebugLogger.Printf("sent response %#v\n", change)

	default:
		utils.ErrorLogger.Printf("server got a non-command message: %#v\n", msg)

	}
}
