package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func CreateClient(cmdQueue chan Command, address string, init chan bool) {
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
		if text == "seek" {
			asdf := Seek{
				Location: "00h23m50s",
				Creator:  "nate",
			}
			result, err := asdf.Render()
			if err != nil {
				fmt.Fprintln(conn, err)
				return
			}
			a, err := json.Marshal(result)
			if err != nil {
				fmt.Fprintln(conn, err)
			} else {
				fmt.Fprintln(conn, string(a))
			}
		}

		if text == "play" {
			asdf := Play{
				Creator: "nate",
			}
			result, err := asdf.Render()
			if err != nil {
				fmt.Fprintln(conn, err)
				return
			}
			a, err := json.Marshal(result)
			if err != nil {
				fmt.Fprintln(conn, err)
			} else {
				fmt.Fprintln(conn, string(a))
			}
		}

		if text == "pause" {
			asdf := Pause{
				Creator: "nate",
			}
			result, err := asdf.Render()
			if err != nil {
				fmt.Fprintln(conn, err)
				return
			}
			a, err := json.Marshal(result)
			if err != nil {
				fmt.Fprintln(conn, err)
			} else {
				fmt.Fprintln(conn, string(a))
			}
		}

		_, err := fmt.Fprintln(conn, text)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}
