package impl

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type Join struct {
	UserId uint16
	RoomId int64
}

func (c *Join) GetHead() string { return "join" }

func (c *Join) ToBits() ([]byte, error) {
	var bits = new(bytes.Buffer)

	if err := binary.Write(bits, binary.BigEndian, c.UserId); err != nil {
		return nil, err
	}
	if err := binary.Write(bits, binary.BigEndian, c.RoomId); err != nil {
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

	if err := binary.Read(buf, binary.BigEndian, &c.UserId); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &c.RoomId); err != nil {
		return err
	}

	return nil
}

// Example:
// ["RoomId=34129"]
func (c *Join) newFromString(s []string) error {
	for _, v := range s {
		v = strings.ToLower(v)
		v = strings.TrimPrefix(v, "-")
		v = strings.TrimPrefix(v, "-")
		switch {
		case strings.HasPrefix(v, "roomid="):
			flag, _ := strings.CutPrefix(v, "roomid=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 16)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.RoomId = int64(i)
		}
	}

	if c.RoomId == 0 {
		return utils.ErrRequiredArgs("join required arg roomid=<uint16>")
	}

	return nil
}
