package impl

import (
	"github.com/nathan-hello/nat-sync/src/utils"
)

type Pause struct {
}

func (c *Pause) ToPlayer(p utils.LocalTarget) ([]byte, error) {
	return []byte(`{"command":["set_property","pause",true]}`), nil
}

func (c *Pause) New(t any) error {
	return nil
}

func (c *Pause) ToBits() ([]byte, error) {
	return nil, nil
}
