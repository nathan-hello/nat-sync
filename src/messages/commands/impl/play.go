package impl

type Play struct {
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
	return `{"command":["set_property","pause",false]}`, nil
}
