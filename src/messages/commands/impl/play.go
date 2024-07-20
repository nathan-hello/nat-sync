package impl

import (
	"fmt"

	"github.com/nathan-hello/nat-sync/src/players"
)

type Play struct {
	User string
}

func (c *Play) ExecuteClient(p players.Player) ([]byte, error) {

	pause := `{"command":["set_property","pause",false]}`
	msg := fmt.Sprintf(`{ "command": ["show-text", "%s paused", 1000, 0] }`, c.User)
	return []byte(pause + "\n" + msg), nil

}
func (c *Play) ExecuteServer() ([]byte, error) { return nil, nil }
func (c *Play) IsEchoed() bool                 { return true }
func (c *Play) NewFromBits(bits []byte) error {
	return nil
}

func (c *Play) NewFromString(s []string) error {
	return nil
}

func (c *Play) ToBits() ([]byte, error) {
	return nil, nil
}
