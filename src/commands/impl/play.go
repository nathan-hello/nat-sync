package impl

import (
	"encoding/json"

	"github.com/nathan-hello/nat-sync/src/commands/impl/players"
)

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
	asdf := players.MpvJson{}

	asdf.Command = append(asdf.Command, "pause")
	asdf.Command = append(asdf.Command, "false")

	mpvCmd, err := json.Marshal(asdf)
	if err != nil {
		return "", err
	}
	return string(mpvCmd), nil
}
