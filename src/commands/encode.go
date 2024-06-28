package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"

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

	err = binary.Write(bits, binary.BigEndian, fixedCmd.Type)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.BigEndian, fixedCmd.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.BigEndian, fixedCmd.Creator)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.BigEndian, fixedCmd.ContentLength)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.BigEndian, fixedCmd.Content)
	if err != nil {
		return nil, err
	}

	fmt.Printf("decoded bytes: %b ", bits.Bytes())

	return bits.Bytes(), nil
}

func cmdToFixedLength(cmd *Command) (*FixedLengthCommand, error) {
	content := cmd.Content.ToBits()
	contentLength := len(content)
	if contentLength > 65535 {
		return nil, utils.ErrFixedContentLength(cmd, cmd.Content)
	}

	return &FixedLengthCommand{
		Type:          uint8(cmd.Type),
		Version:       cmd.Version,
		Creator:       userFixer(cmd.Creator),
		ContentLength: uint16(contentLength),
		Content:       content,
	}, nil
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
