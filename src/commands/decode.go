package commands

import (
	"encoding/binary"

	"github.com/nathan-hello/nat-sync/src/utils"
)

func DecodeCommand(bits []byte) (*Command, error) {
	bType := bits[0]
	bVersion := bits[1:9]
	bCreator := bits[10:41]
	bContentLength := bits[42:57]
	contentLength := binary.LittleEndian.Uint16(bContentLength)
	_ = bits[58:contentLength]

	cmd := Command{}
	for k, v := range typeFixer {
		if bType == v {
			cmd.Type = k
			break
		}
	}
	if cmd.Type == "" {
		return nil, utils.ErrDecodeType(bits)
	}

	cmd.Version = binary.LittleEndian.Uint16(bVersion)
	user := ""
	for _, v := range bCreator {
		if v != ' ' {
			user += string(v)
		}
	}

	cmd.Creator = user

	return &cmd, nil

}
