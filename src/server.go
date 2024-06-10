package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func CreateServer(cmdQueue chan Command) {
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on :4000")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, cmdQueue)
	}
}

func handleConnection(conn net.Conn, cmdQueue chan Command) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	go func() {
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Connection closed")
				return
			}
			fmt.Print("Received from client: ", message)
		}
	}()

	send := func(s string) {
		for {
			time.Sleep(5 * time.Second) // Adjust the interval as needed
			_, err := fmt.Fprintf(conn, "Message from server: %s\n", s)
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}

	for cmd := range cmdQueue {
		err := cmd.Render()
		if err != nil {
			send(MarshalErrToJSON(cmd, err))
		}
		bits, err := json.Marshal(cmd)
		if err != nil {
			send(MarshalErrToJSON(cmd, err))
		}
		send(string(bits))
	}
}
