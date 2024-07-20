package players

import (
	"net"

	"github.com/nathan-hello/nat-sync/src/messages"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type Player interface {
	connect() error
	transmit(conn net.Conn)
	receive(conn net.Conn)

	Quit()
	Launch() error
	AppendQueue(messages.Message)
}

func New(p utils.PlayerTargets) Player {
	switch p {
	case utils.TargetMpv:
		return newMpv()
	case utils.TargetVlc:
		return nil // TODO:
	case utils.TargetTest:
		return nil // TODO:
	}

	return nil
}
