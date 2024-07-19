package impl

type Stop struct {
}

func (c *Stop) NewFromBits([]byte) error     { return nil }
func (c *Stop) NewFromString([]string) error { return nil }

func (c *Stop) ExecuteClient() ([]byte, error) { return nil, nil }
func (c *Stop) ExecuteServer() ([]byte, error) { return nil, nil }
func (c *Stop) IsEchoed() bool                 { return true }
func (c *Stop) ToBits() ([]byte, error)        { return []byte{}, nil }
func (c *Stop) ToMpv() (string, error) {
	return `{"command":["stop"]}`, nil
}
