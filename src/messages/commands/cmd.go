package commands

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/nathan-hello/nat-sync/src/client/players"
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
	ToBits() ([]byte, error)

	ExecuteClient(players.Player) ([]byte, error)
	ExecuteServer() ([]byte, error)

	NewFromBits([]byte) error
	NewFromString([]string) error

	// True  if each client in room runs the command themselves, individually (i.e. Pause).
	// False if admin-related things, such as Kick and Join.
	// If false, ExecuteServer() will return a response in []byte the server will echo that.
	// TODO: If false, ExecuteClient() will return a response in []byte and do nothing.
	IsEchoed() bool
}

var Head = struct {
	Ack    cmdHead
	Change cmdHead
	Kick   cmdHead
	Join   cmdHead
	Pause  cmdHead
	Play   cmdHead
	Seek   cmdHead
	Stop   cmdHead
	Wait   cmdHead
}{
	Ack:    1,
	Change: 2,
	Kick:   3,
	Join:   4,
	Pause:  5,
	Play:   6,
	Seek:   7,
	Stop:   8,
	Wait:   9,
}

func New[T []byte | string](i T) (*Command, error) {
	switch t := any(i).(type) {
	case []byte:
		cmd, err := newCmdFromBits(t)
		if err != nil {
			return nil, err
		}
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

func (c *Command) ExecutePlayer(player players.Player) ([]byte, error) {
	switch player.GetPlayerType() {
	case utils.TargetMpv:
		return c.Sub.ExecuteClient(player)
	}
	return nil, nil
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

	sub, err := getSubFromHead(cmd.Head)
	if err != nil {
		return nil, err
	}

	sub.NewFromBits(cmd.Content)
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

	head, err := getHeadFromString(parts[0])
	if err != nil {
		return nil, err
	}

	sub, err := getSubFromHead(head)
	if err != nil {
		return nil, err
	}

	err = sub.NewFromString(parts[1:])
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

// Register new commands here
func getSubFromHead(head cmdHead) (SubCommand, error) {
	switch head {
	case Head.Ack:
		return &impl.Ack{}, nil
	case Head.Change:
		return &impl.Change{}, nil
	case Head.Kick:
		return &impl.Kick{}, nil
	case Head.Join:
		return &impl.Join{}, nil
	case Head.Pause:
		return &impl.Pause{}, nil
	case Head.Play:
		return &impl.Play{}, nil
	case Head.Seek:
		return &impl.Seek{}, nil
	case Head.Stop:
		return &impl.Stop{}, nil
	case Head.Wait:
		return &impl.Wait{}, nil
	}
	return nil, utils.ErrNoCmdHeadFound(uint8(head))
}

// Register new strings here
func getHeadFromString(s string) (cmdHead, error) {
	switch strings.ToLower(s) {
	case "ack":
		return Head.Ack, nil
	case "change":
		return Head.Change, nil
	case "kick":
		return Head.Kick, nil
	case "join":
		return Head.Join, nil
	case "pause":
		return Head.Pause, nil
	case "play":
		return Head.Play, nil
	case "seek":
		return Head.Seek, nil
	case "stop":
		return Head.Stop, nil
	case "wait":
		return Head.Wait, nil
	default:
		return 0, utils.ErrBadArgs([]string{s})
	}

}

func IsCommand(bits []byte) bool {
	var cmd Command
	buf := bytes.NewReader(bits)

	_ = binary.Read(buf, binary.BigEndian, &cmd.Length)
	_ = binary.Read(buf, binary.BigEndian, &cmd.Type)

	return cmd.Type == utils.MsgCommand
}
