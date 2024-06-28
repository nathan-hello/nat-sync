package src

import (
	"bufio"
	"fmt"
	"net"

	"github.com/nathan-hello/nat-sync/src/commands"
)

func CreateServer(address string, serverInit chan bool) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	serverInit <- true
	fmt.Println("Started server at " + address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Connection closed")
			return
		}

		if len(message) == 8 {
			continue
		}

		fmt.Printf("Received from client: %b\nstring %s", message, string(message))
		dec, err := commands.DecodeCommand(message)
		if err != nil {
			fmt.Println("err: ", err)
		}
		fmt.Printf("%#v\n", dec)
	}
}

// send := func(s string) {
// 	for {
// 		_, err := fmt.Fprintf(conn, "Message from server: %s\n", s)
// 		if err != nil {
// 			fmt.Println("Error sending message:", err)
// 			return
// 		}
// 	}
// }
