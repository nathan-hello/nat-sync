package commands

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/nathan-hello/nat-sync/src/commands/impl"
	"github.com/nathan-hello/nat-sync/src/utils"
)

func DecodeCommand(bits []byte) (*Command, error) {
	buf := bytes.NewReader(bits)

	// Read the fixed-length part of the Command struct
	var cmd Command
	if err := binary.Read(buf, binary.BigEndian, &cmd.Head); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Head):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &cmd.Version); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Version):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &cmd.UserId); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Creator):", err)
		return nil, err
	}

	// Read the remaining bytes into Content
	cmd.Content = make([]byte, buf.Len())
	if err := binary.Read(buf, binary.BigEndian, &cmd.Content); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Content):", err)
		return nil, err
	}

	var sub SubCommand

	switch cmd.Head {
	case ChangeHead:
		sub = &impl.Change{}
	case KickHead:
		sub = &impl.Kick{}
	case JoinHead:
		sub = &impl.Join{}
	case PauseHead:
		sub = &impl.Pause{}
	case PlayHead:
		sub = &impl.Play{}
	case SeekHead:
		sub = &impl.Seek{}
	default:
		return nil, utils.ErrNoCmdHeadFound(bits[0])
	}

	// utils.DebugLogger.Printf("content before FromBits(): %v\n", cmd.Content)

	sub.FromBits(cmd.Content)

	// utils.DebugLogger.Printf("con: %#v\n", sub)

	// utils.DebugLogger.Printf("command full: %#v\n", cmd)

	cmd.Sub = sub

	return &cmd, nil
}

// Returns a *Command without UserId field
func CmdFromString(s string) (*Command, error) {
	parts := strings.Fields(s)

	if len(parts) == 0 {
		return nil, nil
	}

	var head CmdHead
	var sub SubCommand

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
