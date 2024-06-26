package src

import (
	"bufio"
	"encoding/json"
	"errors"
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

func handleNewMessage(msg []byte) (*RenderedCommand, error) {
	// Unmarshal into the generic RenderedCommand struct
	var baseCmd RenderedCommand
	err := json.Unmarshal(msg, &baseCmd)
	if err != nil {
		return nil, err
	}

	// Determine the command type and unmarshal the content accordingly
	switch baseCmd.Command {
	case "seek":
		var content Seek
		err := json.Unmarshal(baseCmd.Content, &content)
		if err != nil {
			return nil, err
		}
		baseCmd.Content = content
	case "pause":
		var content Pause
		if err := json.Unmarshal(baseCmd.Content.([]byte), &content); err != nil {
			return nil, err
		}
		baseCmd.Content = content
	case "play":
		var content Play
		if err := json.Unmarshal(baseCmd.Content.([]byte), &content); err != nil {
			return nil, err
		}
		baseCmd.Content = content
	case "newvideo":
		var content NewVideo
		if err := json.Unmarshal(baseCmd.Content.([]byte), &content); err != nil {
			return nil, err
		}
		baseCmd.Content = content
	default:
		return nil, errors.New("unknown command type")
	}

	// Type switch to handle the specific command type
	switch content := baseCmd.Content.(type) {
	case Seek:
		fmt.Println("seek cmd found:: ", content)
		return &baseCmd, nil
	case Pause:
		fmt.Println("pause cmd found:: ", content)
		return &baseCmd, nil
	case Play:
		fmt.Println("play cmd found:: ", content)
		return &baseCmd, nil
	case NewVideo:
		fmt.Println("newvideo cmd found:: ", content)
		return &baseCmd, nil
	default:
		fmt.Println("default: ", content)
	}
	return nil, errors.New("invalid content")
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
		fmt.Print("Received from client: ", string(message))
		cmd, err := handleNewMessage(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Rendered proper command: ", cmd)
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
