package players

import (
	"bufio"
	"io"
	"net"
	"os/exec"
	"time"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type mpv struct {
	LaunchCmd  string
	LaunchArg  []string
	Exec       *exec.Cmd
	Conn       net.Conn
	SocketPath string
	PlayerType utils.LocalTarget
	ToPlayer   chan PlayerExecutor
	FromPlayer chan []byte
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

	pause := []byte(`{"command": ["observe_property", 1, "pause"]`)
	play := []byte(` {"command": ["observe_property", 1, "play" ]`)
	stop := []byte(` {"command": ["observe_property", 1, "stop" ]`)
	// seek := []byte(`{"command": ["observe_property", 2, "playback-time"], "request_id": 69}`)
	subscribes := [][]byte{pause, play, stop}

	for _, v := range subscribes {
		conn.Write(append(v, byte('\n')))
	}

	go v.transmit(conn)
	go v.receive(conn)
	return nil
}

func (v *mpv) transmit(conn net.Conn) {
	for m := range v.ToPlayer {
		mpvBits, err := m.ToPlayer(utils.TargetMpv)
		if err != nil {
			utils.ErrorLogger.Printf("parsing command to player format. cmd: %#v err: %s", m, err)
			break
		}

		_, err = conn.Write(append(mpvBits, byte('\n')))
		if err != nil {
			utils.ErrorLogger.Printf("sending command to player socket. cmd: %#v err: %s", m, err)
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
