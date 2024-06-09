package src

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

var writeBuffer = make(chan string)

func CreateServer(stopper chan bool) {
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
		go handleConnection(conn, stopper)
	}
}

func handleConnection(conn net.Conn, stopper chan bool) {
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

	go func() {
		for {
			time.Sleep(5 * time.Second) // Adjust the interval as needed
			_, err := fmt.Fprintf(conn, "Message from server\n")
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		}
	}()

	<-stopper
}
