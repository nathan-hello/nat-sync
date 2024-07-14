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

func (c *Kick) ToBits() ([]byte, error) {

	bits := []byte{}
	if c.IsSelf {
		bits = append(bits, 1)
	} else {
		bits = append(bits, 0)
	}

	if c.HideMsg {
		bits = append(bits, 1)
	} else {
		bits = append(bits, 0)
	}

	return bits, nil
}

func (c *Kick) FromBits(bits []byte) error {
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
// ["--UserId=2182", "IsSelf=true", "--HideMsg=false"]
func (c *Kick) FromString(s []string) error {
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
			c.UserId = uint16(i)

		case strings.HasPrefix(v, "--IsSelf="):
			flag, _ := strings.CutPrefix(v, "--IsSelf=")
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

		case strings.HasPrefix(v, "--HideMsg="):
			flag, _ := strings.CutPrefix(v, "--HideMsg=")
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
