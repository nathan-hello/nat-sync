package impl

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

var RoomHead = struct {
	Create     uint8
	Destroy    uint8
	ChangeName uint8
	KickAll    uint8
}{
	Create:     1,
	Destroy:    2,
	ChangeName: 3,
	KickAll:    4,
}

type Room struct {
	NameLen uint8
	Name    string
	Head    uint8
}

func (c *Room) GetHead() string { return "connect" }
func (c *Room) ToBits() ([]byte, error) {
	var bits = new(bytes.Buffer)

	if err := binary.Write(bits, binary.BigEndian, c.NameLen); err != nil {
		return nil, err
	}
	if err := binary.Write(bits, binary.BigEndian, []byte(c.Name)); err != nil {
		return nil, err
	}

	return bits.Bytes(), nil
}

func (c *Room) New(t any) error {
	switch s := t.(type) {
	case []byte:
		return c.newFromBits(s)
	case []string:
		return c.newFromString(s)

	default:
		return utils.ErrBadType
	}
}

func (c *Room) newFromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	if err := binary.Read(buf, binary.BigEndian, &c.NameLen); err != nil {
		return err
	}

	uriBits := make([]byte, c.NameLen)
	if _, err := buf.Read(uriBits); err != nil {
		return err
	}
	c.Name = string(uriBits)

	return nil
}

func (c *Room) Execute() ([]byte, error) {
	// TODO: THIS !!!!!!!!!!!!!!!!
	return nil, nil
}

// Example:
// ["name=nate"]
func (c *Room) newFromString(s []string) error {
	for _, v := range s {
		v = strings.TrimPrefix(v, "-")
		v = strings.TrimPrefix(v, "-")
		switch {
		case strings.HasPrefix(strings.ToLower(v), "name="):
			parts := strings.Split(v, "=")
			if len(parts) < 1 {
				return utils.ErrBadArgs(s)
			}
			name := parts[1]
			nameLen := len(name)
			c.Name = name
			c.NameLen = uint8(nameLen)
			if c.Name == "" {
				return utils.ErrBadArgs(s)
			}
		case strings.HasPrefix(strings.ToLower(v), "action="):
			v = strings.ToLower(v)
			flag := strings.TrimPrefix(v, "action=")
			if flag == "create" {
				c.Head = RoomHead.Create
			}
			if flag == "destroy" {
				c.Head = RoomHead.Destroy
			}
			if flag == "changename" {
				c.Head = RoomHead.ChangeName
			}
			if flag == "kickall" {
				c.Head = RoomHead.KickAll
			}
			if c.Head == 0 {
				return utils.ErrBadArgs(s)
			}
		}
	}
	if c.Name == "" {
		return utils.ErrBadArgs(s)
	}
	return nil
}
