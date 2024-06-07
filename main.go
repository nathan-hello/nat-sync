package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nathan-hello/nat-sync/src/commands"
)

var severToClient commands.Connection
var clientToSever commands.Connection
var signalListener = make(chan os.Signal, 1) // this only works on unix
var cleanupTime = make(chan bool, 1)

var writeBuffer = make(chan string)

func main() {
	args := flag.NewFlagSet("natsync", flag.ExitOnError)
	_ = args.String("client", "4000", "Create a client to join natsync servers")
	_ = args.String("server", "4000", "Become a server for natsync clients")

	go createClient()
	go createServer()

	cmd := commands.Seek{
		CurrentTime: time.Date(0, 0, 0, 0, 10, 20, 0, time.UTC),
		NewTime:     time.Date(0, 0, 0, 0, 10, 8, 0, time.UTC),
	}
	_ = cmd.Parse()

	go waitForInterrupt()
	<-cleanupTime
	os.Exit(0)
}

func createServer() {
	handleNewClient := func(conn net.Conn) {
		defer conn.Close()
		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)

		go func() {
			for {
				msg := <-writeBuffer
				if msg != "\n" {
					_, err := writer.WriteString("hi")
					if err != nil {
						log.Fatalf("err writing in WriteString: %s", err)
					}
					writer.Flush()
				}
			}
		}()

		go func() {
			for {
				msg, err := reader.ReadString('\n')
				if err != nil {
					log.Fatalf("err writing in WriteString: %s", err)
				}
				fmt.Println("msg from server: ", msg)
				writeBuffer <- "Echo: " + msg
			}
		}()

	}

	server, err := net.Listen("tcp", "localhost:4000")
	if err != nil {
		log.Fatalf("err connecting: %s", err)
	}

	for {
		client, err := server.Accept()
		if err != nil {
			log.Fatalf("err connecting: %s", err)
		}
		go handleNewClient(client)
	}
}

func createClient() {
	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		log.Fatalf("createClient couldn't connect to server: %s", err)
	}
	defer conn.Close()
	clientReaderStdIn := bufio.NewReader(os.Stdin)
	clientWriter := bufio.NewWriter(conn)
	clientReader := bufio.NewReader(conn)

	var wg sync.WaitGroup
	wg.Add(2)

	// write new msg
	go func() {
		defer wg.Done()
		for {
			fmt.Println("Type new message: ")

			msg, err := clientReaderStdIn.ReadString('\n')
			if err != nil {
				log.Fatalf("clientReaderStdIn err: %s", err)
			}

			_, err = clientWriter.WriteString(msg)
			if err != nil {
				log.Fatalf("clientWriter err: %s", err)
			}

			err = clientWriter.Flush()
			if err != nil {
				log.Fatalf("clientWriter flush err: %s", err)
			}

		}
	}()

	go func() {
		defer wg.Done()
		for {
			msg, _ := clientReader.ReadString('\n')
			fmt.Println("From the server: ", msg)
		}
	}()
	wg.Wait()
	<-cleanupTime
}

func waitForInterrupt() {
	signal.Notify(signalListener, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP)
	signal := <-signalListener
	fmt.Printf("\nReceived signal %s, exiting\n", signal)
	cleanupTime <- true

}
