package impl

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type Seek struct {
	Hours uint8
	Mins  uint8
	Secs  uint8
}

func (c *Seek) ToBits() ([]byte, error) {
	return []byte{c.Hours, c.Mins, c.Secs}, nil
}

func (c *Seek) FromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	if err := binary.Read(buf, binary.BigEndian, &c.Hours); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.Mins); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &c.Secs); err != nil {
		return err
	}

	return nil
}

// Example:
// ["--Hours=2", "--Mins=19", "--Seconds=0"]
// ["--Uri=\"file:/home/catlover/kitty.jpeg\"", "--IsLocal=true"]
func (c *Seek) FromString(s []string) error {

	init := false
	for _, v := range s {
		v = strings.ToLower(v)
		switch {
		case strings.HasPrefix(v, "--hours="):
			flag, _ := strings.CutPrefix(v, "--hours=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 8)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.Hours = uint8(i)
			init = true
		case strings.HasPrefix(v, "--mins="):
			flag, _ := strings.CutPrefix(v, "--mins=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 8)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.Mins = uint8(i)
			init = true
		case strings.HasPrefix(v, "--secs="):
			flag, _ := strings.CutPrefix(v, "--secs=")
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
		return utils.ErrNoArgs("seek requires --hours, --mins, or --secs. if you want to go to beginning, use \"seek --secs=0\"")
	}

	return nil
}
