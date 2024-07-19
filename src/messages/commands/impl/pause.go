package impl

type Pause struct {
}

func (c *Pause) ExecuteClient() ([]byte, error) { return nil, nil }
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

func (c *Pause) ToMpv() (string, error) {
	return `{"command":["set_property","pause",true]}`, nil
}
