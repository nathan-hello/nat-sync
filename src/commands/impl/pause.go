package impl

type Pause struct {
}

func (c *Pause) ToBits() ([]byte, error) {
	return nil, nil
}

func (c *Pause) FromBits(bits []byte) error {
	return nil
}

// Example:
// []
func (c *Pause) FromString(s []string) error {
	return nil
}
