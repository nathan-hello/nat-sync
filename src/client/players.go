package client

import (
	"bufio"
	"net"
	"os/exec"
	"time"

	"github.com/nathan-hello/nat-sync/src/commands"
	"github.com/nathan-hello/nat-sync/src/utils"
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
}

func LaunchPlayer(p *LaunchParams) error {
	var cmd *exec.Cmd
	var player string
	var playerArgs []string
	if p.Player == "mpv" {
		player = "mpv"
		playerArgs = []string{"--idle", "--force-window", "--input-ipc-server=" + p.SocketPath}
	}
	utils.DebugLogger.Printf("starting %s with args: %v\n", player, playerArgs)
	cmd = exec.Command(player, playerArgs...)

	err := cmd.Start()
	if err != nil {
		return err
	}
	utils.DebugLogger.Println(p.Player, "started")
	go handlePlayerConnection(p)

	cmd.Wait()
	return nil

}

func handlePlayerConnection(p *LaunchParams) {
	const maxRetries = 10
	var conn net.Conn
	var err error

	for i := 0; i < maxRetries; i++ {
		time.Sleep(250 * time.Millisecond)
		conn, err = net.Dial("unix", p.SocketPath)
		if err == nil {
			break
		}
	}

	if err != nil {
		utils.ErrorLogger.Printf("error connecting to socket after %d retries: %s", maxRetries, err)
		return
	}
	defer conn.Close()

	go readResponses(conn, p)

	utils.DebugLogger.Printf("player is ready for cmds\n")
	p.Init <- true
	for cmd := range p.ToClient {
		mpvStr, err := cmd.Sub.ToMpv()
		if err != nil {
			utils.ErrorLogger.Printf("error parsing command to player format. cmd: %#v err: %s", cmd.Sub, err)
			break
		}
		utils.DebugLogger.Printf("sending cmd to player. cmd: %s", mpvStr)
		_, err = conn.Write([]byte(mpvStr + "\n"))
		if err != nil {
			utils.ErrorLogger.Printf("error sending command to player socket. cmd: %#v err: %s", cmd.Sub, err)
			break
		}
	}
}

func readResponses(conn net.Conn, p *LaunchParams) {
	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			utils.ErrorLogger.Printf("error reading from socket: %s", err)
			break
		}
		utils.NoticeLogger.Printf("Response from MPV: %s\n", response)
	}
}
