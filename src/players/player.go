package players

import (
	"errors"
	"net"

	"github.com/nathan-hello/nat-sync/src/messages"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type PlayerTargets string

const (
	Test PlayerTargets = ""
	Mpv  PlayerTargets = "mpv"
	Vlc  PlayerTargets = "vlc"
)

type Player interface {
	launch() error
	connect()
	transmit(conn net.Conn)
	receive(conn net.Conn)

	AppendQueue(messages.Message)
}

func New(p PlayerTargets) (Player, error) {
	switch p {
	case Mpv:
		mpv := newMpv()

		err := mpv.launch()
		if err != nil {
			return nil, err
		}

		go mpv.connect()

		return mpv, nil
	case Vlc:
		return nil, utils.ErrNotImplemented("vlc")
	case Test:
		return nil, nil
	}

	return nil, errors.New("newplayer was given incorrect .Player arg")
}
