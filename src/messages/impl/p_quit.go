package impl

import (
	"github.com/nathan-hello/nat-sync/src/utils"
)

type Quitter interface {
	Quit()
}

type Quit struct {
}

func (c *Quit) IsPlayer() bool { return true }

func (c *Quit) GetHead() string { return "quit" }
func (c *Quit) ToPlayer(p utils.LocalTarget) ([]byte, error) {
	return []byte(`{"command":["quit", 0]}`), nil

}
func (c *Quit) New(t any) error { return nil }
func (c *Quit) ToBits() ([]byte, error) {
	return nil, nil
}
