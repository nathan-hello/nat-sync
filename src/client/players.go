package client

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/nathan-hello/nat-sync/src/commands"
)

type LaunchTargets string

const (
	LaunchMpv LaunchTargets = "mpv"
	LaunchVlc LaunchTargets = "vlc"
)

type LaunchParams struct {
	Player     LaunchTargets
	SocketPath string
	Init       chan bool
	ToClient   chan commands.Command
	ToError    chan error
}

func LaunchPlayer(p *LaunchParams) error {
	var cmd *exec.Cmd
	var player string
	var playerArgs []string
	if p.Player == "mpv" {
		player = "mpv"
		playerArgs = []string{"--idle", "--force-window", "--input-ipc-server=" + p.SocketPath}
	}
	fmt.Printf("starting %s with args: %v\n", player, playerArgs)
	cmd = exec.Command(player, playerArgs...)

	err := cmd.Start()
	if err != nil {
		return err
	}
	fmt.Println(p.Player, "started")
	go handlePlayerConnection(p)

	cmd.Wait()
	return nil

}

func handlePlayerConnection(p *LaunchParams) {
	const maxRetries = 10
	var conn net.Conn
	var err error

	for i := 0; i < maxRetries; i++ {
		conn, err = net.Dial("unix", p.SocketPath)
		if err == nil {
			break
		}
		fmt.Printf("Attempt %d: error connecting to socket: %v\n", i+1, err)
		time.Sleep(500 * time.Millisecond) // Wait before retrying
	}

	if err != nil {
		p.ToError <- fmt.Errorf("error connecting to socket after %d retries: %w", maxRetries, err)
		return
	}
	defer conn.Close()

	go readResponses(conn, p)

	fmt.Printf("player is ready for cmds")
	p.Init <- true
	for cmd := range p.ToClient {
		fmt.Printf("player received cmd on ToClient chan: %#v\n", cmd)
		mpvStr, err := cmd.Sub.ToMpv()
		if err != nil {
			p.ToError <- err
			break
		}
		_, err = conn.Write([]byte(mpvStr + "\n"))
		if err != nil {
			p.ToError <- err
			break
		}
	}
}

func readResponses(conn net.Conn, p *LaunchParams) {
	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			p.ToError <- fmt.Errorf("error reading from socket: %w", err)
			break
		}
		fmt.Printf("Response from MPV: %s\n", response)
	}
}
