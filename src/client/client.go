package client

import (
	"bufio"
	"io"
	"net"
	"time"

	"github.com/nathan-hello/nat-sync/src/client/players"
	"github.com/nathan-hello/nat-sync/src/messages"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type ClientParams struct {
	ServerAddress string
	JoinedRooms   []int64
	Player        players.Player
	InputReader   io.Reader
}

func CreateClient(p *ClientParams) {
	var conn net.Conn
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		conn, err = net.Dial("tcp", p.ServerAddress)
		if err == nil {
			break
		}
		if i == maxRetries-1 {
			utils.ErrorLogger.Println("Error connecting to server:", err)
			return
		}
		utils.DebugLogger.Println("Couldn't connect to server, trying again:  " + p.ServerAddress)
		time.Sleep(250 * time.Millisecond)
	}
	utils.DebugLogger.Println("Started client connected to " + p.ServerAddress)

	go receive(conn, p)
	go transmit(conn, p)

}

func receive(conn net.Conn, p *ClientParams) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {

		utils.DebugLogger.Printf("waitreader\n\n\n")
		msgs, err := messages.WaitReader(reader)
		utils.DebugLogger.Printf("got msg %#v\n\n", msgs)
		if err == io.EOF {
			utils.DebugLogger.Printf("connection closed EOF\n")
			return
		}
		if err != nil {
			utils.ErrorLogger.Printf("client got a bad message. error: %s\n", err)
		}

		for _, v := range msgs {
			switch msg := v.Sub.(type) {
			case messages.PlayerCommand:
				utils.DebugLogger.Printf("appending cmd to playerqueue. cmd: %#v\n", msg)
				p.Player.AppendQueue(msg)
			// case messages.AdminCommand:
			// 	switch admin := msg.(type) {
			// 	case *impl.Join:
			// 		p.CurrentRoom = admin.RoomId
			// 		if !slices.Contains(p.JoinedRooms, admin.RoomId) {
			// 			p.JoinedRooms = append(p.JoinedRooms, admin.RoomId)
			// 		}
			// 	}
			default:
				utils.ErrorLogger.Printf("client was given a command that was not a player command! %#v\n", msg)
			}
		}
	}
}

func transmit(conn net.Conn, p *ClientParams) {
	scanner := bufio.NewScanner(p.InputReader)
	for scanner.Scan() { // this blocks the terminal
		text := scanner.Text()
		if IsLocalCommand(text) {
			playerCmd, err := NewLocalCmd(text, p.Player)
			if err != nil {
				utils.ErrorLogger.Println(err)
				continue
			}
			newPlayer, err := playerCmd.Sub.ExecuteClient()
			if err != nil {
				utils.ErrorLogger.Println(err)
				continue
			}
			p.Player = newPlayer
			continue
		}

		macro := messages.IsMacro(text)
		if macro != nil {
			sendMsgs(conn, macro)
			continue
		}

		msgs, err := messages.New(text, nil)
		if err != nil {
			utils.ErrorLogger.Println(err)
			continue
		}
		sendMsgs(conn, msgs)
	}

}

func sendMsgs(conn net.Conn, msgs []messages.Message) {
	for _, m := range msgs {
		bits, err := m.ToBits()
		if err != nil {
			utils.ErrorLogger.Println("cmd.ToBits() in client transmit. err: ", err)
			continue
		}
		_, err = conn.Write(bits)
		if err != nil {
			utils.ErrorLogger.Printf("client writing bits failed. bits: %b\n", bits)
		}
	}
}
