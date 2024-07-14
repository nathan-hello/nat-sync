package commands

import (
	"fmt"
	"strings"

	"github.com/nathan-hello/nat-sync/src/commands/impl"
	"github.com/nathan-hello/nat-sync/src/utils"
)

// Returns a *Command without UserId field
func CmdFromString(s string) (*Command, error) {
	parts := strings.Fields(s)
	var head CmdHead
	var sub SubCommand

	fmt.Println("text given to cmdfromstring: ", s)
	fmt.Println("parts: ", parts)
	fmt.Println("parts[0]: ", parts[0])

	switch strings.ToLower(parts[0]) {
	case "change":
		head = ChangeHead
		sub = &impl.Change{}
	case "kick":
		head = KickHead
		sub = &impl.Kick{}
	case "join":
		head = JoinHead
		sub = &impl.Join{}
	case "pause":
		head = PauseHead
		sub = &impl.Pause{}
	case "play":
		head = PlayHead
		sub = &impl.Play{}
	case "seek":
		head = SeekHead
		sub = &impl.Seek{}
	default:
		return nil, utils.ErrBadArgs(parts)
	}

	if len(parts) > 0 {
		parts = parts[1:]
	}

	fmt.Println("parts sending to fromstring(): ", parts)
	err := sub.FromString(parts)

	if err != nil {
		return nil, err
	}

	content, err := sub.ToBits()
	if err != nil {
		return nil, err
	}

	return &Command{
		Head:    head,
		Version: CurrentVersion,
		Sub:     sub,
		Content: content,
	}, nil

}
