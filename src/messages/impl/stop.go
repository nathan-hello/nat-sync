package impl

import (
	"github.com/nathan-hello/nat-sync/src/utils"
)

type Stop struct {
}

func (c *Stop) GetHead() string { return "stop" }
func (c *Stop) New(t any) error {
	switch s := t.(type) {
	case []byte:
		return c.newFromBits(s)
	case []string:
		return c.newFromString(s)
	default:
		return utils.ErrBadType
	}
}

func (c *Stop) newFromBits([]byte) error     { return nil }
func (c *Stop) newFromString([]string) error { return nil }

func (c *Stop) ToPlayer(p utils.LocalTarget) ([]byte, error) {
	return []byte(`{"command":["stop"]}`), nil
}

func (c *Stop) ToBits() ([]byte, error) { return []byte{}, nil }
