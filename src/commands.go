package src

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/nathan-hello/nat-sync/src/utils"
)

const CurrentVersion = "0.0.1"

// Command:
//	seek:      0000_0001
//	pause:     0000_0010
//	play:      0000_0011
//	new_video: 0000_0100

type Command struct {
	Command       uint8
	Version       uint8
	Creator       uint8 // map[uint8]string for mapping bits to usernames
	ContentLength uint8
	Content       Renderer
}

type Renderer interface{ Render() (*Command, error) }

func EncodeCommand(cmd Command) ([]byte, error) {
	if cmd.Content == nil {
		return nil, utils.ErrNoContent
	}
	bits := bytes.Buffer{}

	err := binary.Write(&bits, binary.LittleEndian, cmd.Command)
	if err != nil {
		return nil, utils.ByteEncodingErr(bits)
	}
	err = binary.Write(&bits, binary.LittleEndian, 0000_0001)
	if err != nil {
		return nil, utils.ByteEncodingErr(bits)
	}
	err = binary.Write(&bits, binary.LittleEndian, cmd.Creator)
	if err != nil {
		return nil, utils.ByteEncodingErr(bits)
	}
	err = binary.Write(&bits, binary.LittleEndian, len(cmd.Content))
	if err != nil {
		return nil, utils.ByteEncodingErr(bits)
	}
	err = binary.Write(&bits, binary.LittleEndian, cmd.Content)
	if err != nil {
		return nil, utils.ByteEncodingErr(bits)
	}

	return nil, nil
}

func DecodeCommand(bits []byte) (*Command, error) {
	return nil, nil
}

type Seek struct {
	Location string
}

func (c *Seek) Render() (*Command, error) {
	t, err := time.ParseDuration(c.Location)
	if err != nil {
		return nil, err
	}
	c.Location = fmt.Sprintf("%.0f", t.Seconds())
	return &Command{Command: 0000_0001, Content: c}, nil
}

type Pause struct{}

func (c *Pause) Render() (*Command, error) {
	return &Command{Command: 0000_0010, Content: c}, nil
}

type Play struct{}

func (c *Play) Render() (*Command, error) {
	return &Command{Command: 0000_0100, Content: c}, nil
}

type NewVideo struct {
	Creator string
	Uri     string
	Local   bool
}
