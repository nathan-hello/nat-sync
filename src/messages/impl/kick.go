package impl

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type Kick struct {
	UserId  uint16
	IsSelf  bool
	HideMsg bool
}

func (c *Kick) New(t any) error {
	switch s := t.(type) {
	case []byte:
		return c.newFromBits(s)
	case []string:
		return c.newFromString(s)

	default:
		return utils.ErrBadType
	}
}

func (c *Kick) newFromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	if err := binary.Read(buf, binary.BigEndian, &c.UserId); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &c.IsSelf); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &c.HideMsg); err != nil {
		return err
	}

	return nil
}

// Example:
// ["UserId=2182", "IsSelf=true", "--HideMsg=false"]
func (c *Kick) newFromString(s []string) error {
	for _, v := range s {
		v = strings.ToLower(v)
		v = strings.TrimPrefix(v, "-")
		v = strings.TrimPrefix(v, "-")
		switch {
		case strings.HasPrefix(v, "userid="):
			flag, _ := strings.CutPrefix(v, "userid=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			i, err := strconv.ParseUint(flag, 10, 16)
			if err != nil {
				return utils.ErrBadArgs(s)
			}
			c.UserId = uint16(i)

		case strings.HasPrefix(v, "isself="):
			flag, _ := strings.CutPrefix(v, "isself=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			if flag == "true" {
				c.IsSelf = true
				continue
			}
			if flag == "false" {
				c.IsSelf = false
				continue
			}
			return utils.ErrBadArgs(s)

		case strings.HasPrefix(v, "hidemsg="):
			flag, _ := strings.CutPrefix(v, "hidemsg=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			if flag == "true" {
				c.HideMsg = true
				continue
			}
			if flag == "false" {
				c.HideMsg = false
				continue
			}
			return utils.ErrBadArgs(s)
		default:
			return utils.ErrBadArgs(s)
		}
	}

	if c.UserId == 0 {
		utils.ErrRequiredArgs("kick requires a userid")
	}

	return nil
}

func (c *Kick) ToBits() ([]byte, error) {

	var bits = new(bytes.Buffer)

	if err := binary.Write(bits, binary.BigEndian, c.UserId); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, c.IsSelf); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, c.HideMsg); err != nil {
		return nil, err
	}

	return bits.Bytes(), nil
}
