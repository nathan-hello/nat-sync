package impl

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type Join struct {
	RoomId   int64
	UserId   int64
	UserLen  uint16
	Username string
}

func (c *Join) GetHead() string { return "join" }

func (c *Join) ToBits() ([]byte, error) {
	if len(c.Username) > 65534 {
		return nil, utils.ErrBadArgs([]string{"username too large", c.Username})
	}
	c.UserLen = uint16(len(c.Username))
	var bits = new(bytes.Buffer)

	if err := binary.Write(bits, binary.BigEndian, c.RoomId); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, c.UserId); err != nil {
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

	if err := binary.Read(buf, binary.BigEndian, &c.RoomId); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &c.UserId); err != nil {
		return err
	}

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
		case strings.HasPrefix(v, "roomid="):
			v = strings.ToLower(v)
			flag, _ := strings.CutPrefix(v, "roomid=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 16)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.RoomId = int64(i)
		case strings.HasPrefix(v, "userid="):
			v = strings.ToLower(v)
			flag, _ := strings.CutPrefix(v, "userid=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 16)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.UserId = int64(i)
		case strings.HasPrefix(v, "username="):
			flag, _ := strings.CutPrefix(v, "username=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			if len(flag) > 65534 {
				return utils.ErrBadArgs([]string{"username too large", strings.Join(s, " ")})
			}
			c.UserLen = uint16(len(flag))
			c.Username = flag
		}
	}

	if c.RoomId == 0 {
		return utils.ErrRequiredArgs("join required arg roomid=<uint16>")
	}

	return nil
}
