package impl

import "fmt"

type Play struct {
	User string
}

func (c *Play) ExecuteClient() ([]byte, error) { return nil, nil }
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
func (c *Play) ToMpv() (string, error) {
	pause := `{"command":["set_property","pause",false]}`
	msg := fmt.Sprintf(`{ "command": ["show-text", "%s paused", 1000, 0] }`, c.User)
	return pause + "\n" + msg, nil
}
