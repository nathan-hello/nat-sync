package players

import (
	"net"

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

type Player interface {
	connect() error
	transmit(conn net.Conn)
	receive(conn net.Conn)

	Quit()
	Launch() error
	AppendQueue(PlayerExecutor)
	GetPlayerType() utils.LocalTarget
}

func New(p utils.LocalTarget) Player {
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
