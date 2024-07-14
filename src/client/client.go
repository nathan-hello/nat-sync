package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/nathan-hello/nat-sync/src/commands"
)

func CreateClient(address string, init chan bool) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Started client connected to " + address)

	playerInit := make(chan bool)
	go launchPlayer("mpv", playerInit)
	<-playerInit

	init <- true

	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Connection closed")
				return
			}
			fmt.Print("Received from server: ", message)
		}
	}()

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
		fmt.Printf("cmd before encodecommand(): %#v\n\n", cmd)

		bits, err := commands.EncodeCommand(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tosend := string(bits)
		fmt.Fprintln(conn, tosend)
	}
}
