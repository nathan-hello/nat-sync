package commands

import (
	"bytes"
	"encoding/binary"
)

const CurrentVersion = 1

type CmdHead uint8

const (
	ChangeHead CmdHead = iota
	KickHead   CmdHead = iota
	JoinHead   CmdHead = iota
	PauseHead  CmdHead = iota
	PlayHead   CmdHead = iota
	SeekHead   CmdHead = iota
)

type Command struct {
	Head    CmdHead
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

func (cmd *Command) ToBits() ([]byte, error) {
	bits := new(bytes.Buffer)

	if cmd.Sub != nil {
		cmd.Sub = nil
	}

	if cmd.Version == 0 {
		cmd.Version = CurrentVersion
	}

	err := binary.Write(bits, binary.BigEndian, cmd.Head)
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

	// utils.DebugLogger.Printf("decoded bytes: %b ", bits.Bytes())

	return bits.Bytes(), nil
}
