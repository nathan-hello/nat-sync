package players

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
	Player   LaunchTargets
	Init     chan bool
	ToClient chan commands.Command
}

type Mpv struct {
	LaunchCmd  string
	LaunchArg  []string
	SocketPath string
	Init       chan bool
}

func NewMpv() *Mpv {
	socketPath := "/tmp/nat-sync-mpv-socket"
	return &Mpv{
		LaunchCmd:  "mpv",
		LaunchArg:  []string{"--idle", "--force-window", "--input-ipc-server=" + socketPath},
		SocketPath: socketPath,
	}
}

func (p *Mpv) GetLaunchCmd() *exec.Cmd {
	return exec.Command(p.LaunchCmd, p.LaunchArg...)
}

func (v *Mpv) Connect(p *LaunchParams) {
	const maxRetries = 10
	var conn net.Conn
	var err error

	for i := 0; i < maxRetries; i++ {
		time.Sleep(250 * time.Millisecond)
		conn, err = net.Dial("unix", v.SocketPath)
		if err == nil {
			break
		}

		if i == maxRetries-1 {
			utils.ErrorLogger.Printf("connecting to socket after %d retries: %s", maxRetries, err)
			return
		}
	}

	defer conn.Close()

	p.Init <- true
	go mpvTransmit(conn, p)
	mpvReceive(conn)

}
func mpvTransmit(conn net.Conn, p *LaunchParams) {
	for cmd := range p.ToClient {
		mpvStr, err := cmd.Sub.ToMpv()
		if err != nil {
			utils.ErrorLogger.Printf("parsing command to player format. cmd: %#v err: %s", cmd.Sub, err)
			break
		}
		utils.DebugLogger.Printf("sending cmd to player. cmd: %s", mpvStr)
		_, err = conn.Write([]byte(mpvStr + "\n"))
		if err != nil {
			utils.ErrorLogger.Printf("sending command to player socket. cmd: %#v err: %s", cmd.Sub, err)
			break
		}
	}
}

func mpvReceive(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			utils.ErrorLogger.Printf("reading from socket: %s", err)
			break
		}
		utils.NoticeLogger.Printf("mpv response: %s\n", response)
	}
}
