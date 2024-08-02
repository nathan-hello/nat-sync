package impl

import "github.com/nathan-hello/nat-sync/src/client/players"

type Quitter interface {
	Quit()
}

type Quit struct {
	QuitFunc Quitter
}

func (c *Quit) ExecuteClient(p players.Player) ([]byte, error) {
	return []byte(`{"command":["quit", 0]}`), nil

}
func (c *Quit) ExecuteServer() ([]byte, error) { return nil, nil }
func (c *Quit) IsEchoed() bool                 { return true }
func (c *Quit) NewFromBits(bits []byte) error {
	return nil
}

func (c *Quit) NewFromString(s []string) error {
	return nil
}

func (c *Quit) ToBits() ([]byte, error) {
	return nil, nil
}
