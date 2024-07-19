package client

import (
	"bufio"
	"io"
	"net"
	"time"

	"github.com/nathan-hello/nat-sync/src/messages"
	"github.com/nathan-hello/nat-sync/src/players"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type ClientParams struct {
	ServerAddress string
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
		msgs, err := messages.WaitReader(reader)
		if err != nil {
			utils.ErrorLogger.Printf("client got a bad message. error: %s\n", err)
		}
		for _, v := range msgs {
			p.Player.AppendQueue(v)
		}
	}
}

func transmit(conn net.Conn, p *ClientParams) {
	scanner := bufio.NewScanner(p.InputReader)
	for scanner.Scan() { // this blocks the terminal
		text := scanner.Text()

		// utils.DebugLogger.Printf("new reader text: %s\n", text)

		msgs, err := messages.New(text)
		if err != nil {
			utils.ErrorLogger.Println(err)
			continue
		}

		for _, m := range msgs {
			bits, err := m.ToBits()
			if err != nil {
				utils.ErrorLogger.Println("cmd.ToBits() in client transmit. err: ", err)
				continue
			}
			// utils.DebugLogger.Printf("client sending bits: %b, length: %d\n", bits, len(bits))
			_, err = conn.Write(bits)
			if err != nil {
				utils.ErrorLogger.Printf("client writing bits failed. bits: %b\n", bits)
			}
		}
	}
}
