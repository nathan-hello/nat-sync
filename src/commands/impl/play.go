package impl

type Play struct {
}

func (c *Play) IsEchoed() bool { return true }

func (c *Play) ToBits() ([]byte, error) {
	return nil, nil
}

func (c *Play) FromBits(bits []byte) error {
	return nil
}

// Example:
// []
func (c *Play) FromString(s []string) error {
	return nil
}

func (c *Play) ToMpv() (string, error) {
	return `{"command":["set_property","pause",false]}`, nil
}
