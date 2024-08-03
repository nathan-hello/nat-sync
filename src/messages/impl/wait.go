package impl

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
	"time"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type Wait struct {
	Secs uint8
}

func (c *Wait) GetHead() string { return "wait" }

func (c *Wait) Execute() ([]byte, error) {
	time.Sleep(time.Duration(c.Secs) * time.Second)
	return []byte("success"), nil
}

func (c *Wait) ToBits() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, c.Secs); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

}

func (c *Wait) New(t any) error {
	switch s := t.(type) {
	case []byte:
		return c.newFromBits(s)
	case []string:
		return c.newFromString(s)
	default:
		return utils.ErrBadType
	}
}

func (c *Wait) newFromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	if err := binary.Read(buf, binary.BigEndian, &c.Secs); err != nil {
		return err
	}

	return nil
}

// Example:
// ["Secs=15"]
// ["Uri=\"file:/home/catlover/kitty.jpeg\"", "--IsLocal=true"]
func (c *Wait) newFromString(s []string) error {

	init := false
	for _, v := range s {
		v = strings.ToLower(v)
		v = strings.TrimPrefix(v, "-")
		v = strings.TrimPrefix(v, "-")
		switch {
		case strings.HasPrefix(v, "secs="):
			flag, _ := strings.CutPrefix(v, "secs=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 8)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.Secs = uint8(i)
			init = true
		default:
			return utils.ErrBadArgs(s)
		}
	}

	if !init {
		return utils.ErrNoArgs("Wait requires \"secs\" argument")
	}

	return nil
}
