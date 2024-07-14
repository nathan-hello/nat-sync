package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/nathan-hello/nat-sync/src/commands/impl"
	"github.com/nathan-hello/nat-sync/src/utils"
)

func DecodeCommand(bits []byte) (*Command, error) {
	buf := bytes.NewReader(bits)

	// Read the fixed-length part of the Command struct
	var cmd Command
	if err := binary.Read(buf, binary.BigEndian, &cmd.Head); err != nil {
		fmt.Println("binary.Read failed (Head):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &cmd.Version); err != nil {
		fmt.Println("binary.Read failed (Version):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &cmd.UserId); err != nil {
		fmt.Println("binary.Read failed (Creator):", err)
		return nil, err
	}

	// Read the remaining bytes into Content
	cmd.Content = make([]byte, buf.Len())
	if err := binary.Read(buf, binary.BigEndian, &cmd.Content); err != nil {
		fmt.Println("binary.Read failed (Content):", err)
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

	// fmt.Printf("content before FromBits(): %v\n", cmd.Content)

	sub.FromBits(cmd.Content)

	// fmt.Printf("con: %#v\n", sub)

	// fmt.Printf("command full: %#v\n", cmd)

	cmd.Sub = sub

	return &cmd, nil
}
