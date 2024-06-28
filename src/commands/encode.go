package commands

import (
	"bytes"
	"encoding/binary"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type FixedLengthCommand struct {
	Type          uint8
	Version       uint16
	Creator       [32]byte
	Content       []byte
	ContentLength uint16
}

func EncodeCommand(cmd Command) ([]byte, error) {
	if cmd.Content == nil {
		return nil, utils.ErrNoContent
	}

	bits := new(bytes.Buffer)

	fixedCmd, err := cmdToFixedLength(&cmd)
	if err != nil {
		return nil, err
	}

	if cmd.Version == 0 {
		cmd.Version = CurrentVersion
	}

	err = binary.Write(bits, binary.LittleEndian, fixedCmd.Type)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.LittleEndian, fixedCmd.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.LittleEndian, fixedCmd.Creator)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.LittleEndian, fixedCmd.ContentLength)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.LittleEndian, fixedCmd.Content)
	if err != nil {
		return nil, err
	}

	return bits.Bytes(), nil
}

func cmdToFixedLength(cmd *Command) (*FixedLengthCommand, error) {
	fixedType := typeFixer[cmd.Type]

	content := cmd.Content.ToBits()
	contentLength := len(content)
	if contentLength > 65535 {
		return nil, utils.ErrFixedContentLength(cmd, cmd.Content)
	}

	return &FixedLengthCommand{
		Type:          fixedType,
		Version:       cmd.Version,
		Creator:       userFixer(cmd.Creator),
		ContentLength: uint16(contentLength),
		Content:       content,
	}, nil
}

var typeFixer = map[CmdHead]uint8{
	"seek":      0000_0001,
	"pause":     0000_0010,
	"play":      0000_0011,
	"new_video": 0000_0100,
	"join":      0000_0101,
}

func userFixer(s string) [32]byte {
	fixedUser := [32]byte{}

	for _, v := range s {
		if v > 127 {
			s = "anon"
			break
		}
	}

	copy(fixedUser[:], s)

	for i := len(s); i < 32; i++ {
		fixedUser[i] = ' '
	}

	return fixedUser
}
