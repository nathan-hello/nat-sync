package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/nathan-hello/nat-sync/src/utils"
)

func EncodeCommand(cmd Command) ([]byte, error) {
	if cmd.Content == nil {
		return nil, utils.ErrNoContent
	}

	bits := new(bytes.Buffer)

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
	err = binary.Write(bits, binary.BigEndian, cmd.Creator)
	if err != nil {
		return nil, err
	}
	// err = binary.Write(bits, binary.BigEndian, fixedCmd.ContentLength)
	// if err != nil {
	// 	return nil, err
	// }
	// err = binary.Write(bits, binary.BigEndian, fixedCmd.Content)
	// if err != nil {
	// 	return nil, err
	// }

	fmt.Printf("decoded bytes: %b ", bits.Bytes())

	return bits.Bytes(), nil
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
