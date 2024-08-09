package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"slices"
	"time"

	"github.com/nathan-hello/nat-sync/src/client/players"
	"github.com/nathan-hello/nat-sync/src/messages"
	"github.com/nathan-hello/nat-sync/src/messages/impl"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type ClientParams struct {
	ServerAddress string
	CurrentRoom   int64
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
		msg, err := messages.WaitReader(reader)
		utils.DebugLogger.Printf("got msg %#v\n\n", msg)
		if err == io.EOF {
			utils.DebugLogger.Printf("connection closed EOF\n")
			return
		}
		if err != nil {
			utils.ErrorLogger.Printf("client got a bad message. error: %s\n", err)
		}

		handleMessage(conn, p, msg)
	}
}

func handleMessage(_ net.Conn, p *ClientParams, msg messages.Message) {
	switch msg := msg.Sub.(type) {
	case messages.PlayerCommand:
		utils.DebugLogger.Printf("appending cmd to playerqueue. cmd: %#v\n", msg)
		p.Player.AppendQueue(msg)
	case messages.AdminCommand:
		switch admin := msg.(type) {
		case *impl.Accept:
			p.CurrentRoom = admin.RoomId
			if !slices.Contains(p.JoinedRooms, admin.RoomId) {
				p.JoinedRooms = append(p.JoinedRooms, admin.RoomId)
			}
		}
	default:
		utils.ErrorLogger.Printf("client was given a command that was not a player command! %#v\n", msg)
	}

}

func transmit(conn net.Conn, p *ClientParams) {
	scanner := bufio.NewScanner(p.InputReader)
	for scanner.Scan() { // this blocks the terminal
		fmt.Printf("<%d> ", p.CurrentRoom)
		text := scanner.Text()
		if text == "" {
			continue
		}
		if IsLocalCommand(text) {
			cmd, err := NewLocal(text[1:], p.Player)
			if err != nil {
				utils.ErrorLogger.Println(err)
				continue
			}
			cmd.Execute(p)
			continue
		}

		macro := messages.IsMacro(text)
		if macro != nil {
			for _, v := range macro {
				sendMsgs(conn, v)
			}
			continue
		}

		text = fmt.Sprintf("%#x %s", p.CurrentRoom, text)
		msg := messages.Message{}

		err := msg.TextUnmarshaller([]byte(text))
		if err != nil {
			utils.ErrorLogger.Println(err)
			continue
		}
		sendMsgs(conn, msg)
	}

}

func sendMsgs(conn net.Conn, msg messages.Message) {
	bits, err := msg.MarshalBinary()
	if err != nil {
		utils.ErrorLogger.Println("cmd.ToBits() in client transmit. err: ", err)
	}

	_, err = conn.Write(bits)
	if err != nil {
		utils.ErrorLogger.Printf("client writing bits failed. bits: %b\n", bits)
	}
}
