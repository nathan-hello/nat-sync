package impl

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type ChangeActions uint16

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

func (c *Change) IsEchoed() bool { return true }

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
func (c *Change) FromString(s []string) error {
	for _, v := range s {
		v = strings.ToLower(v)
		switch {
		case strings.HasPrefix(v, "--uri="):
			flag, _ := strings.CutPrefix(v, "--uri=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			c.Uri = flag
			c.UriLength = uint32(len(flag))
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
				return utils.ErrBadArgs(append(s, flag))
			}
		case strings.HasPrefix(v, "--hours="):
			flag, _ := strings.CutPrefix(v, "--hours=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 8)
			if err != nil {
				return utils.ErrBadArgs(append(s, flag))
			}
			c.Timestamp.Hours = uint16(i)
		case strings.HasPrefix(v, "--mins="):
			flag, _ := strings.CutPrefix(v, "--mins=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 8)
			if err != nil {
				return utils.ErrBadArgs(append(s, flag))
			}
			c.Timestamp.Mins = uint16(i)
		case strings.HasPrefix(v, "--secs="):
			flag, _ := strings.CutPrefix(v, "--secs=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 8)
			if err != nil {
				return utils.ErrBadArgs(append(s, flag))
			}
			c.Timestamp.Secs = uint16(i)
		default:
			return utils.ErrBadArgs(s)
		}
	}

	return nil
}

func (c *Change) ToMpv() (string, error) {
	asdf := MpvJson{}

	asdf.Command = append(asdf.Command, "loadfile")
	asdf.Command = append(asdf.Command, c.Uri)

	if c.Action == ChgAppend {
		asdf.Command = append(asdf.Command, "append-play")
	}
	// if c.Action == ChgImmediate // Immediately playing is the default behavior

	mpvCmd, err := json.Marshal(asdf)
	if err != nil {
		return "", err
	}
	return string(mpvCmd), nil
}
