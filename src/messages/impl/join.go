package impl

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type Join struct {
	RoomId uint16
}

func (c *Join) ToBits() ([]byte, error) {
	bits := make([]byte, 2)
	binary.BigEndian.PutUint16(bits, c.RoomId)

	return bits, nil
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
			c.RoomId = uint16(i)
		}
	}

	if c.RoomId == 0 {
		return utils.ErrRequiredArgs("join required arg roomid=<uint16>")
	}

	return nil
}
