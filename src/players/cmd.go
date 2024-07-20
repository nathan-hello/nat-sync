package players

import (
	"errors"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type playHead uint8

type Local struct {
	Length  uint16
	Type    utils.MsgType
	Head    playHead
	Version uint16
	UserId  uint16
	Sub     SubLocal
}

type SubLocal interface {
	ExecuteClient() (Player, error)
	NewFromString([]string) error
}

var Head = struct {
	LaunchPlayer playHead
	QuitPlayer   playHead
}{
	LaunchPlayer: 1,
	QuitPlayer:   2,
}

func NewPlayerCmd(i string, player Player) (*Local, error) {
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
		Type:    utils.MsgCommand,
		Version: utils.CurrentVersion,
		Sub:     sub,
	}, nil
}

// TODO: stop assumption that mpv is the only player
type LaunchPlayer struct {
	CurrentPlayer Player
	Target        utils.PlayerTargets
}

func (l *LaunchPlayer) NewFromString(s []string) error {
	l.Target = utils.TargetMpv
	return nil
}

func (l *LaunchPlayer) ExecuteClient() (Player, error) {
	if l.CurrentPlayer != nil {
		l.CurrentPlayer.Quit()
	}
	asdf := New(l.Target)
	asdf.Launch()
	return asdf, nil
}

type QuitPlayer struct {
	Player Player
}

func (l *QuitPlayer) NewFromString(s []string) error {
	return nil
}

func (l *QuitPlayer) ExecuteClient() (Player, error) {
	if l.Player != nil {
		l.Player.Quit()
		return nil, nil
	}
	return nil, errors.New("player is already dead")
}

// Register new commands here
func newSub(head playHead, player Player) (SubLocal, error) {
	switch head {
	case Head.LaunchPlayer:
		return &LaunchPlayer{CurrentPlayer: player}, nil
	case Head.QuitPlayer:
		return &QuitPlayer{Player: player}, nil
	}
	return nil, utils.ErrNoCmdHeadFound(uint8(head))
}

// Register new strings here
func getHeadFromString(s string) (playHead, error) {
	switch strings.ToLower(s) {
	case "launch":
		return Head.LaunchPlayer, nil
	case "quit":
		return Head.QuitPlayer, nil
	default:
		return 0, utils.ErrBadArgs([]string{s})
	}

}

func IsPlayerCommand(s string) bool {
	return strings.HasPrefix(s, "$")
}
