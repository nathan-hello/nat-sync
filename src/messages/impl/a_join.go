package impl

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type Join struct {
	RoomLen  uint8
	RoomName string
	UserLen  uint8
	Username string
}

func (c *Join) IsAdmin() bool { return true }

func (c *Join) GetHead() string { return "join" }

func (c *Join) ToBits() ([]byte, error) {
	if len(c.Username) > 255 {
		return nil, utils.ErrBadArgs([]string{"username too large", c.Username})
	}
	c.UserLen = uint8(len(c.RoomName))

	if len(c.RoomName) > 255 {
		return nil, utils.ErrBadArgs([]string{"roomname too large", c.RoomName})
	}
	c.RoomLen = uint8(len(c.RoomName))

	var bits = new(bytes.Buffer)

	if err := binary.Write(bits, binary.BigEndian, c.RoomLen); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, []byte(c.RoomName)); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, c.UserLen); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, []byte(c.Username)); err != nil {
		return nil, err
	}

	return bits.Bytes(), nil
}

func (c *Join) New(t any) error {
	switch s := t.(type) {
	case []byte:
		return c.newFromBits(s)
	case []string:
		return c.newFromString(s)

	default:
		return utils.ErrBadType
	}
}

func (c *Join) newFromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	if err := binary.Read(buf, binary.BigEndian, &c.RoomLen); err != nil {
		return err
	}
	roomBits := make([]byte, c.RoomLen)
	if err := binary.Read(buf, binary.BigEndian, roomBits); err != nil {
		return err
	}
	c.RoomName = string(roomBits)

	if err := binary.Read(buf, binary.BigEndian, &c.UserLen); err != nil {
		return err
	}

	userBits := make([]byte, c.UserLen)

	if err := binary.Read(buf, binary.BigEndian, userBits); err != nil {
		return err
	}

	c.Username = string(userBits)

	return nil
}

// Example:
// ["RoomId=34129"]
func (c *Join) newFromString(s []string) error {
	for _, v := range s {
		v = strings.TrimPrefix(v, "-")
		v = strings.TrimPrefix(v, "-")
		switch {
		case strings.HasPrefix(v, "roomname="):
			v = strings.ToLower(v)
			flag, _ := strings.CutPrefix(v, "roomname=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			if len(flag) > 255 {
				return utils.ErrBadArgs([]string{"roomname too large. max is 255 characters", strings.Join(s, " ")})
			}
			c.RoomLen = uint8(len(c.RoomName))
			c.RoomName = flag
		case strings.HasPrefix(v, "username="):
			flag, _ := strings.CutPrefix(v, "username=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			if len(flag) > 255 {
				return utils.ErrBadArgs([]string{"username too large. max is 255 characters", strings.Join(s, " ")})
			}
			c.UserLen = uint8(len(flag))
			c.Username = flag
		}
	}

	if c.RoomName == "" || c.Username == "" {
		return utils.ErrRequiredArgs("join required arg roomname=<string> username=<string>")
	}

	return nil
}
