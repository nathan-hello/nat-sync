package src

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
		cmd := commands.Command{Version: commands.CurrentVersion, Creator: "a"}
		if text == "seek" {
			asdf := commands.Seek{
				Hours: 0,
				Mins:  23,
				Secs:  29,
			}
			cmd.Type = "seek"
			cmd.Content = asdf
			bits, err := commands.EncodeCommand(cmd)
			if err != nil {
				fmt.Println("error in seek: ", err)
			}
			fmt.Fprintln(conn, bits)
		}

		if text == "play" {
			asdf := commands.Play{}
			cmd.Type = "play"
			cmd.Content = asdf
			bits, err := commands.EncodeCommand(cmd)
			if err != nil {
				fmt.Println("error in play: ", err)
			}
			fmt.Fprintln(conn, bits)
		}

		if text == "pause" {
			asdf := commands.Pause{}
			cmd.Type = "pause"
			cmd.Content = asdf
			bits, err := commands.EncodeCommand(cmd)
			if err != nil {
				fmt.Println("error in pause: ", err)
			}
			fmt.Fprintln(conn, bits)
		}

		_, err := fmt.Fprintln(conn, text)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}
