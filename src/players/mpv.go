package players

import (
	"bufio"
	"io"
	"net"
	"os/exec"
	"time"

	"github.com/nathan-hello/nat-sync/src/utils"
)

// One-method interface because that's all we're
// using here. Messages interface implements this.
//
// This means players package doesn't need to import
// messages package for the Messages interface.
// Otherwise, this interface isn't being used for
// any special composability. Just the lack of import.
type PlayerExecutor interface {
	ExecutePlayer(Player) ([]byte, error)
}

type mpv struct {
	LaunchCmd  string
	LaunchArg  []string
	Exec       *exec.Cmd
	Conn       net.Conn
	SocketPath string
	PlayerType utils.LocalTarget
	ToPlayer   chan PlayerExecutor
}

func newMpv() *mpv {
	socketPath := "/tmp/nat-sync-mpv-socket"
	c := make(chan PlayerExecutor)

	return &mpv{
		LaunchCmd:  "mpv",
		LaunchArg:  []string{"--idle", "--force-window", "--input-ipc-server=" + socketPath},
		SocketPath: socketPath,
		ToPlayer:   c,
		PlayerType: utils.TargetMpv,
	}
}

func (v *mpv) Launch() error {
	cmd := exec.Command(v.LaunchCmd, v.LaunchArg...)
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = v.connect()
	v.Exec = cmd
	return err
}

func (p *mpv) AppendQueue(cmd PlayerExecutor) {
	p.ToPlayer <- cmd
}

func (p *mpv) Quit() {
	if p.Conn != nil {
		p.Conn.Close()
	}
	p.Exec.Process.Kill()
	close(p.ToPlayer)
}

func (v *mpv) connect() error {
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
			return err
		}
	}
	utils.DebugLogger.Printf("new conn at %#v\n", conn)

	go v.transmit(conn)
	go v.receive(conn)
	return nil
}

func (v *mpv) transmit(conn net.Conn) {
	for m := range v.ToPlayer {
		mpvStr, err := m.ExecutePlayer(v)
		if err != nil {
			utils.ErrorLogger.Printf("parsing command to player format. cmd: %#v err: %s", m, err)
			break
		}

		utils.DebugLogger.Printf("sending cmd to player. cmd: %s", mpvStr)

		_, err = conn.Write(append(mpvStr, byte('\n')))
		if err != nil {
			utils.ErrorLogger.Printf("sending command to player socket. cmd: %#v err: %s", m, err)
			break
		}
	}
}

func (v *mpv) receive(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			utils.ErrorLogger.Printf("reading from socket: %s", err)
			break
		}
		utils.NoticeLogger.Printf("mpv response: %s\n", response)
	}
}

func (v *mpv) GetPlayerType() utils.LocalTarget {
	return v.PlayerType
}
