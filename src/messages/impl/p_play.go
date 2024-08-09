package impl

import (
	"fmt"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type Play struct {
	User string
}

func (c *Play) IsPlayer() bool  { return true }
func (c *Play) GetHead() string { return "play" }
func (c *Play) ToPlayer(p utils.LocalTarget) ([]byte, error) {

	pause := `{"command":["set_property","pause",false]}`
	msg := fmt.Sprintf(`{ "command": ["show-text", "%s paused", 1000, 0] }`, c.User)
	return []byte(pause + "\n" + msg), nil

}

func (c *Play) New(t any) error { return nil }
func (c *Play) ToBits() ([]byte, error) {
	return nil, nil
}
