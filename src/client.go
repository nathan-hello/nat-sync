package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func CreateClient(cmdQueue chan Command) {
	conn, err := net.Dial("tcp", ":4000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

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
				NewTime: "00h23m50s",
			}
			asdf.Render()
			a, err := json.Marshal(asdf)
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
