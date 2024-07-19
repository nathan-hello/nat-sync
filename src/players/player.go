package players

import (
	"errors"
	"net"

	"github.com/nathan-hello/nat-sync/src/messages"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type Player interface {
	launch() error
	connect()
	transmit(conn net.Conn)
	receive(conn net.Conn)

	AppendQueue(messages.Message)
}

func New(p utils.PlayerTargets) (Player, error) {
	switch p {
	case utils.TargetMpv:
		mpv := newMpv()

		err := mpv.launch()
		if err != nil {
			return nil, err
		}

		go mpv.connect()

		return mpv, nil
	case utils.TargetVlc:
		return nil, utils.ErrNotImplemented("vlc")
	case utils.TargetTest:
		return nil, nil
	}

	return nil, errors.New("newplayer was given incorrect .Player arg")
}
