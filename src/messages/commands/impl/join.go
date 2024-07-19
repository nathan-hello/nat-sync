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

func (c *Join) ExecuteClient() ([]byte, error) { return nil, nil }
func (c *Join) ExecuteServer() ([]byte, error) { return nil, nil }

func (c *Join) IsEchoed() bool { return false }

func (c *Join) NewFromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	if err := binary.Read(buf, binary.BigEndian, &c.RoomId); err != nil {
		return err
	}

	return nil
}

// Example:
// ["RoomId=34129"]
func (c *Join) NewFromString(s []string) error {
	for _, v := range s {
		v = strings.ToLower(v)
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

func (c *Join) ToBits() ([]byte, error) {
	bits := make([]byte, 2)
	binary.BigEndian.PutUint16(bits, c.RoomId)

	return bits, nil
}

func (c *Join) ToMpv() (string, error) {
	return "", nil // not a player command!
}
