package impl

import "github.com/nathan-hello/nat-sync/src/client/players"

type Pause struct {
}

func (c *Pause) ExecuteClient(_ players.Player) ([]byte, error) {
	return []byte(`{"command":["set_property","pause",true]}`), nil
}
func (c *Pause) ExecuteServer() ([]byte, error) { return nil, nil }
func (c *Pause) IsEchoed() bool                 { return true }
func (c *Pause) NewFromBits(bits []byte) error {
	return nil
}

func (c *Pause) NewFromString(s []string) error {
	return nil
}

func (c *Pause) ToBits() ([]byte, error) {
	return nil, nil
}
