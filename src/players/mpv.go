package players

import (
	"bufio"
	"net"
	"os/exec"
	"time"

	"github.com/nathan-hello/nat-sync/src/messages"
	"github.com/nathan-hello/nat-sync/src/messages/commands"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type mpv struct {
	LaunchCmd  string
	LaunchArg  []string
	SocketPath string
	Conn       net.Conn
	ToPlayer   chan messages.Message
}

func newMpv() *mpv {
	socketPath := "/tmp/nat-sync-mpv-socket"
	c := make(chan messages.Message)

	return &mpv{
		LaunchCmd:  "mpv",
		LaunchArg:  []string{"--idle", "--force-window", "--input-ipc-server=" + socketPath},
		SocketPath: socketPath,
		ToPlayer:   c,
	}
}

func (v *mpv) launch() error {
	cmd := exec.Command(v.LaunchCmd, v.LaunchArg...)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (v *mpv) connect() {
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
	go v.transmit(conn)
	go v.receive(conn)
}

func (v *mpv) transmit(conn net.Conn) {
	for v := range v.ToPlayer {
		switch msg := v.(type) {
		case *commands.Command:
			mpvStr, err := msg.Sub.ToMpv()
			if err != nil {
				utils.ErrorLogger.Printf("parsing command to player format. cmd: %#v err: %s", msg.Sub, err)
				break
			}
			utils.DebugLogger.Printf("sending cmd to player. cmd: %s", mpvStr)
			_, err = conn.Write([]byte(mpvStr + "\n"))
			if err != nil {
				utils.ErrorLogger.Printf("sending command to player socket. cmd: %#v err: %s", msg.Sub, err)
				break
			}
			// TODO: handle acks?
		default:
			continue
		}
	}
}

func (v *mpv) receive(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// TODO: do something interesting with mpv responses
		response, err := reader.ReadString('\n')
		if err != nil {
			utils.ErrorLogger.Printf("reading from socket: %s", err)
			break
		}
		utils.NoticeLogger.Printf("mpv response: %s\n", response)
	}
}
func (p *mpv) AppendQueue(cmd messages.Message) {
	p.ToPlayer <- cmd
}
