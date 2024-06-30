package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/nathan-hello/nat-sync/src/utils"
)

func DecodeCommand(bits []byte) (*Command, error) {
	// bType := bits[0]
	// bVersion := bits[1:9]
	// bCreator := bits[10:41]
	// bContentLength := bits[42:57]
	// bContent := bits[58:binary.LittleEndian.Uint16(bContentLength)]

	buf := new(bytes.Buffer)
	buf.Write(bits)
	cmd := Command{}

	err := binary.Read(buf, binary.BigEndian, &cmd)
	if err != nil {
		return nil, err
	}

	var user string
	for _, v := range cmd.Creator {
		if v > 127 {
			user = "anon"
			break
		}
		if v != ' ' {
			user = user + string(v)
		}
	}

	var sub SubCommand
	switch cmd.Head {
	case SeekHead:
		sub = &Seek{}
	case PlayHead:
		sub = &Play{}
	case PauseHead:
		sub = &Pause{}
	case NewVideoHead:
		sub = &NewVideo{}
	}
	if sub == nil {
		return nil, utils.ErrNoCmdHeadFound(bits[0])
	}

	fmt.Printf("con: %#v\n", sub)
	sub.FromBits(cmd.Content)
	fmt.Printf("con: %#v\n", sub)

	fmt.Printf("command full: %#v\n", cmd)

	return &cmd, nil

}
