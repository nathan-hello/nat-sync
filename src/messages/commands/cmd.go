package commands

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/nathan-hello/nat-sync/src/messages/commands/impl"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type cmdHead uint8

type Command struct {
	Length  uint16
	Type    utils.MsgType
	Head    cmdHead
	Version uint16
	UserId  uint16
	Content []byte
	Sub     SubCommand
}

type SubCommand interface {
	FromString(s []string) error
	FromBits(bits []byte) error
	ToBits() ([]byte, error)
	IsEchoed() bool
	ToMpv() (string, error)
}

var Head = struct {
	Change cmdHead
	Kick   cmdHead
	Join   cmdHead
	Pause  cmdHead
	Play   cmdHead
	Seek   cmdHead
}{
	Change: 1,
	Kick:   2,
	Join:   3,
	Pause:  4,
	Play:   5,
	Seek:   6,
}

func New[T []byte | string](i T) (*Command, error) {
	switch t := any(i).(type) {
	case []byte:
		cmd, err := newCmdFromBits(t)
		if err != nil {
			return nil, err
		}
		cmd.UserId = 1000
		return cmd, nil
	case string:
		cmd, err := newCmdFromString(t)
		if err != nil {
			return nil, err
		}
		cmd.UserId = 1000
		return cmd, nil
	default:
		return nil, utils.ErrImpossible
	}
}

func (cmd *Command) ToBits() ([]byte, error) {
	bits := new(bytes.Buffer)

	if cmd.Sub != nil {
		cmd.Sub = nil
	}
	if cmd.Version == 0 {
		cmd.Version = utils.CurrentVersion
	}
	if cmd.Type == 0 {
		cmd.Type = utils.MsgCommand
	}

	err := binary.Write(bits, binary.BigEndian, cmd.Type)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.BigEndian, cmd.Head)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.BigEndian, cmd.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.BigEndian, cmd.UserId)
	if err != nil {
		return nil, err
	}

	err = binary.Write(bits, binary.BigEndian, cmd.Content)
	if err != nil {
		return nil, err
	}

	cmd.Length = uint16(len(bits.Bytes()))
	finalBits := new(bytes.Buffer)

	err = binary.Write(finalBits, binary.BigEndian, cmd.Length)
	if err != nil {
		return nil, err
	}

	_, err = finalBits.Write(bits.Bytes())
	if err != nil {
		return nil, err
	}

	// utils.DebugLogger.Printf("decoded bytes: %b ", finalBits.Bytes())
	return finalBits.Bytes(), nil
}

func newCmdFromBits(bits []byte) (*Command, error) {
	buf := bytes.NewReader(bits)

	// Read the fixed-length part of the Command struct
	var cmd Command
	if err := binary.Read(buf, binary.BigEndian, &cmd.Length); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Length):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &cmd.Type); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Type):", err)
		return nil, err
	}
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
	case Head.Change:
		sub = &impl.Change{}
	case Head.Kick:
		sub = &impl.Kick{}
	case Head.Join:
		sub = &impl.Join{}
	case Head.Pause:
		sub = &impl.Pause{}
	case Head.Play:
		sub = &impl.Play{}
	case Head.Seek:
		sub = &impl.Seek{}
	default:
		return nil, utils.ErrNoCmdHeadFound(uint8(cmd.Head))
	}

	sub.FromBits(cmd.Content)
	cmd.Sub = sub

	// utils.DebugLogger.Printf("decoded cmd: %#v\n", cmd)

	return &cmd, nil
}

// Returns a *Command without UserId field
func newCmdFromString(s string) (*Command, error) {
	parts := strings.Fields(s)

	if len(parts) == 0 {
		return nil, nil
	}

	var head cmdHead
	var sub SubCommand

	switch strings.ToLower(parts[0]) {
	case "change":
		head = Head.Change
		sub = &impl.Change{}
	case "kick":
		head = Head.Kick
		sub = &impl.Kick{}
	case "join":
		head = Head.Join
		sub = &impl.Join{}
	case "pause":
		head = Head.Pause
		sub = &impl.Pause{}
	case "play":
		head = Head.Play
		sub = &impl.Play{}
	case "seek":
		head = Head.Seek
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
		Type:    utils.MsgCommand,
		Version: utils.CurrentVersion,
		Sub:     sub,
		Content: content,
	}, nil

}

func IsCommand(bits []byte) bool {
	var cmd Command
	buf := bytes.NewReader(bits)

	_ = binary.Read(buf, binary.BigEndian, &cmd.Length)
	_ = binary.Read(buf, binary.BigEndian, &cmd.Type)

	return cmd.Type == utils.MsgCommand
}
