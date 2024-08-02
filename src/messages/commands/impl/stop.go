package impl

import "github.com/nathan-hello/nat-sync/src/client/players"

type Stop struct {
}

func (c *Stop) NewFromBits([]byte) error     { return nil }
func (c *Stop) NewFromString([]string) error { return nil }

func (c *Stop) ExecuteClient(p players.Player) ([]byte, error) {

	return []byte(`{"command":["stop"]}`), nil
}
func (c *Stop) ExecuteServer() ([]byte, error) { return nil, nil }
func (c *Stop) IsEchoed() bool                 { return true }
func (c *Stop) ToBits() ([]byte, error)        { return []byte{}, nil }
