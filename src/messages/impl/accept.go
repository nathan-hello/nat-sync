package impl

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type AcceptAction uint8
type Accept struct {
	Action AcceptAction
	RoomId int64
}

var AcceptHead = struct {
	Ok            AcceptAction
	RoomCreated   AcceptAction
	BadUsername   AcceptAction
	RoomNotExists AcceptAction
	NotAllowed    AcceptAction
}{
	Ok:            1,
	RoomCreated:   2,
	BadUsername:   3,
	RoomNotExists: 4,
	NotAllowed:    5,
}

func (c *Accept) GetHead() string { return "accept" }

func (c *Accept) ToBits() ([]byte, error) {
	var bits = new(bytes.Buffer)

	if err := binary.Write(bits, binary.BigEndian, c.Action); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, c.RoomId); err != nil {
		return nil, err
	}

	return bits.Bytes(), nil
}

func (c *Accept) New(t any) error {
	switch s := t.(type) {
	case []byte:
		return c.newFromBits(s)
	case []string:
		return c.newFromString(s)

	default:
		return utils.ErrBadType
	}
}

func (c *Accept) newFromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	if err := binary.Read(buf, binary.BigEndian, &c.Action); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &c.RoomId); err != nil {
		return err
	}

	return nil
}

// Example:
// ["RoomId=34129"]
func (c *Accept) newFromString(s []string) error {
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
		case strings.HasPrefix(v, "action="):
			v = strings.ToLower(v)
			flag, _ := strings.CutPrefix(v, "action=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 8)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.Action = AcceptAction(i)
		}
	}

	if c.RoomId == 0 || c.Action == 0 {
		return utils.ErrRequiredArgs("accept required args roomid=<uint64> action=<AcceptAction>")
	}

	return nil
}
