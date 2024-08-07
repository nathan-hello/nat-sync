package client

import (
	"strconv"
	"strings"

	"github.com/nathan-hello/nat-sync/src/client/players"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type locHead uint8

type Local struct {
	Head    locHead
	Version uint16
	Sub     SubLocal
}

func (c *Local) Execute(p *ClientParams) {
	c.Sub.Execute(p)
}

type SubLocal interface {
	New(any) error
	Execute(*ClientParams)
}

var Head = struct {
	LocalChangeRoom locHead
}{
	LocalChangeRoom: 1,
}

func NewLocal(i string, player players.Player) (*Local, error) {
	i = strings.TrimPrefix(i, "$")
	parts := strings.Fields(i)

	if len(parts) == 0 {
		return nil, nil
	}

	head, err := getHeadFromString(parts[0])
	if err != nil {
		return nil, err
	}

	sub, err := newSub(head)
	if err != nil {
		return nil, err
	}

	err = sub.New(parts[1:])
	if err != nil {
		return nil, err
	}

	return &Local{
		Head:    head,
		Version: utils.CurrentVersion,
		Sub:     sub,
	}, nil
}

// Register new commands here
func newSub(head locHead) (SubLocal, error) {
	switch head {
	case Head.LocalChangeRoom:
		return &Swap{}, nil
	}
	return nil, utils.ErrNoCmdHeadFound(uint8(head))
}

// Register new strings here
func getHeadFromString(s string) (locHead, error) {
	switch strings.ToLower(s) {
	case "swap":
		return Head.LocalChangeRoom, nil
	default:
		return 0, utils.ErrBadArgs([]string{s})
	}

}

func IsLocalCommand(s string) bool {
	return strings.HasPrefix(s, "/")
}

type Swap struct {
	RoomId int64
}

func (c *Swap) New(t any) error {
	switch s := t.(type) {
	case []string:
		for _, v := range s {
			switch {
			case strings.HasPrefix(v, "roomid="):
				v := strings.ToLower(v)
				flag, _ := strings.CutPrefix(v, "roomid=")
				flag, _ = strings.CutPrefix(flag, "\"")
				flag, _ = strings.CutSuffix(flag, "\"")
				i, err := strconv.ParseUint(flag, 10, 16)
				if err != nil {
					return utils.ErrBadArgs(s)
				}
				c.RoomId = int64(i)

			}
		}
	}
	return nil
}

func (c *Swap) Execute(p *ClientParams) {
	p.CurrentRoom = c.RoomId
}
