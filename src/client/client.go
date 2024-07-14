package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"slices"

	"github.com/nathan-hello/nat-sync/src/commands"
)

type ClientParams struct {
	ServerAddress string
	Init          chan bool
	ToClient      chan commands.Command
	ToError       chan error
}

func CreateClient(p ClientParams) {
	conn, err := net.Dial("tcp", p.ServerAddress)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Started client connected to " + p.ServerAddress)

	lp := LaunchParams{
		Player:     LaunchMpv,
		SocketPath: "/tmp/nat-sync-mpv-socket",
		Init:       make(chan bool),
		ToClient:   p.ToClient,
		ToError:    p.ToError,
	}
	go LaunchPlayer(&lp)
	<-lp.Init

	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Println("Connection closed")
				return
			}
			fmt.Printf("Received from server: %b\n", message)
			if slices.Equal(message, []byte("200")) {
				fmt.Printf("Received OK server: %s\n", message)
				continue
			}

			dec, err := commands.DecodeCommand(message)
			if err != nil {
				fmt.Println("err: ", err)
			}
			fmt.Printf("%#v\n", dec)
			go func() {
				fmt.Printf("sending cmd to ToClient: %#v\n", dec)
				p.ToClient <- *dec
			}()
		}
	}()
	p.Init <- true

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		creator := uint16(10126)

		cmd, err := commands.CmdFromString(text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		cmd.UserId = creator

		bits, err := commands.EncodeCommand(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tosend := string(bits)
		fmt.Fprintln(conn, tosend)
	}
}
