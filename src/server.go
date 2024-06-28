package src

import (
	"bufio"
	"fmt"
	"net"
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
		fmt.Printf("Received from client: %b\nstring %s", message, string(message))
		// cmd, err := handleNewMessage(message)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// fmt.Println("Rendered proper command: ", cmd)
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
//
// for cmd := range cmdQueue {
// 	result, err := cmd.Render()
// 	if err != nil {
// 		send(MarshalErrToJSON(cmd, err))
// 	}
// 	bits, err := json.Marshal(result)
// 	if err != nil {
// 		send(MarshalErrToJSON(cmd, err))
// 	}
// 	send(string(bits))
// }
//
