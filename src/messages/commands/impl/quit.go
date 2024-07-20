package impl

type Quitter interface {
	Quit()
}

type Quit struct {
	QuitFunc Quitter
}

func (c *Quit) ExecuteClient() ([]byte, error) { return nil, nil }
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
func (c *Quit) ToMpv() (string, error) {
	return `{"command":["quit", 0]}`, nil
}
