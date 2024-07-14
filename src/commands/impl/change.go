package impl

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type ChangeActions uint8

const (
	ChgAppend    ChangeActions = 1
	ChgImmediate ChangeActions = 2
)

type Change struct {
	Action    ChangeActions
	Timestamp Seek
	UriLength uint32
	Uri       string
}

func (c *Change) ToBits() ([]byte, error) {

	var bits = new(bytes.Buffer)

	if err := binary.Write(bits, binary.BigEndian, c.Action); err != nil {
		return nil, err
	}

	t, err := c.Timestamp.ToBits()
	if err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, t); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, c.UriLength); err != nil {
		return nil, err
	}
	if err := binary.Write(bits, binary.BigEndian, []byte(c.Uri)); err != nil {
		return nil, err
	}

	return bits.Bytes(), nil
}

func (c *Change) FromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	if err := binary.Read(buf, binary.BigEndian, &c.Action); err != nil {
		return err
	}

	var t Seek
	if err := binary.Read(buf, binary.BigEndian, &t); err != nil {
		return err
	}
	c.Timestamp = t

	if err := binary.Read(buf, binary.BigEndian, &c.UriLength); err != nil {
		return err
	}

	uriBits := make([]byte, c.UriLength)
	if _, err := buf.Read(uriBits); err != nil {
		return err
	}
	c.Uri = string(uriBits)

	return nil
}

// Example:
// ["--Uri=\"asdf.com/cats\"", "--Action=\"immediate\""]
// ["--Uri=\"file:/home/catlover/kitty.jpeg\"", "--Action=\"append\""]
// TODO: add timestamp string parsing
func (c *Change) FromString(s []string) error {
	for _, v := range s {
		v = strings.ToLower(v)
		switch {
		case strings.HasPrefix(v, "--uri="):
			flag, _ := strings.CutPrefix(v, "--uri=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			c.Uri = flag
		case strings.HasPrefix(v, "--action="):
			flag, _ := strings.CutPrefix(v, "--action=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			switch strings.ToLower(flag) {
			case "append":
				c.Action = ChgAppend
			case "immediate":
				c.Action = ChgImmediate
			default:
				return utils.ErrBadArgs(s)
			}
		default:
			return utils.ErrBadArgs(s)
		}
	}

	return nil
}
