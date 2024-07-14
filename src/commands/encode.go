package commands

import (
	"bytes"
	"encoding/binary"
)

func EncodeCommand(cmd *Command) ([]byte, error) {
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

	// fmt.Printf("decoded bytes: %b ", bits.Bytes())

	return bits.Bytes(), nil
}
