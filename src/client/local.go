package client

import (
	"strings"

	"github.com/nathan-hello/nat-sync/src/client/players"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type playHead uint8

type Local struct {
	Head    playHead
	Version uint16
	Sub     SubLocal
}

type SubLocal interface {
	ExecuteClient() (players.Player, error)
	NewFromString([]string) error
}

var Head = struct {
	PlayerOpen playHead
	PlayerQuit playHead
}{
	PlayerOpen: 1,
	PlayerQuit: 2,
}

func NewLocalCmd(i string, player players.Player) (*Local, error) {
	i = strings.TrimPrefix(i, "$")
	parts := strings.Fields(i)

	if len(parts) == 0 {
		return nil, nil
	}

	head, err := getHeadFromString(parts[0])
	if err != nil {
		return nil, err
	}

	sub, err := newSub(head, player)
	if err != nil {
		return nil, err
	}

	err = sub.NewFromString(parts[1:])
	if err != nil {
		return nil, err
	}

	return &Local{
		Head:    head,
		Version: utils.CurrentVersion,
		Sub:     sub,
	}, nil
}

// TODO: stop assumption that mpv is the only player
type PlayerOpen struct {
	CurrentPlayer players.Player
	Target        utils.LocalTarget
}

func (l *PlayerOpen) NewFromString(s []string) error {
	l.Target = utils.TargetMpv
	return nil
}

func (l *PlayerOpen) ExecuteClient() (players.Player, error) {
	if l.CurrentPlayer != nil {
		l.CurrentPlayer.Quit()
	}
	asdf := players.New(l.Target)
	asdf.Launch()
	return asdf, nil
}

type PlayerQuit struct {
	Player players.Player
}

func (l *PlayerQuit) NewFromString(s []string) error {
	return nil
}

func (l *PlayerQuit) ExecuteClient() (players.Player, error) {
	if l.Player != nil {
		l.Player.Quit()
		return nil, nil
	}
	return nil, utils.ErrPlayerAlreadyDead
}

// Register new commands here
func newSub(head playHead, player players.Player) (SubLocal, error) {
	switch head {
	case Head.PlayerOpen:
		return &PlayerOpen{CurrentPlayer: player}, nil
	case Head.PlayerQuit:
		return &PlayerQuit{Player: player}, nil
	}
	return nil, utils.ErrNoCmdHeadFound(uint8(head))
}

// Register new strings here
func getHeadFromString(s string) (playHead, error) {
	switch strings.ToLower(s) {
	case "launch":
		return Head.PlayerOpen, nil
	case "quit":
		return Head.PlayerQuit, nil
	default:
		return 0, utils.ErrBadArgs([]string{s})
	}

}

func IsLocalCommand(s string) bool {
	return strings.HasPrefix(s, "/")
}
