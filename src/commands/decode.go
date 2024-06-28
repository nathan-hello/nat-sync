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
	fixed := FixedLengthCommand{}

	binary.Read(buf, binary.BigEndian, fixed)

	fmt.Printf("fixed after binary read: %#v\n", fixed)

	var user string
	for _, v := range fixed.Creator {
		if v > 127 {
			user = "anon"
			break
		}
		if v != ' ' {
			user = user + string(v)
		}
	}

	fmt.Printf("user: %s\n", user)

	var sub SubCommand
	switch fixed.Type {
	case SeekHead:
		fmt.Printf("found seek")
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
	sub.FromBits(fixed.Content)
	fmt.Printf("con: %#v\n", sub)

	cmd := Command{
		Type:    fixed.Type,
		Version: fixed.Version,
		Creator: user,
		Content: sub,
	}
	fmt.Printf("command full: %#v\n", cmd)

	return &cmd, nil

}
