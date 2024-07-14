package impl

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

// uint16 to prevent binary reader
// from interpreting 1010 (decimal 10) as EOF
type Seek struct {
	Hours uint16
	Mins  uint16
	Secs  uint16
}

func (c *Seek) ToBits() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, c.Hours); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, c.Mins); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, c.Secs); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

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
			i, err := strconv.ParseUint(flag, 10, 16)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.Hours = uint16(i)
			init = true
		case strings.HasPrefix(v, "--mins="):
			flag, _ := strings.CutPrefix(v, "--mins=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 16)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.Mins = uint16(i)
			init = true
		case strings.HasPrefix(v, "--secs="):
			flag, _ := strings.CutPrefix(v, "--secs=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 16)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.Secs = uint16(i)
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
