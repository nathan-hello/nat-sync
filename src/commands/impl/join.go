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
	bits := []byte{}
	binary.BigEndian.PutUint16(bits, c.RoomId)

	return bits, nil
}

func (c *Join) FromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	if err := binary.Read(buf, binary.BigEndian, &c.RoomId); err != nil {
		return err
	}

	return nil
}

// Example:
// ["--UserId=834129"]
func (c *Join) FromString(s []string) error {
	for _, v := range s {
		switch {
		case strings.HasPrefix(v, "--UserId="):
			flag, _ := strings.CutPrefix(v, "--UserId=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 16)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.RoomId = uint16(i)
		default:
			return utils.ErrBadArgs(s)
		}
	}

	return nil
}
