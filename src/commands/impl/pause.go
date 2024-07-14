package impl

import (
	"encoding/json"

	"github.com/nathan-hello/nat-sync/src/commands/impl/players"
)

type Pause struct {
}

func (c *Pause) IsEchoed() bool { return true }

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

func (c *Pause) ToMpv() (string, error) {
	asdf := players.MpvJson{}

	asdf.Command = append(asdf.Command, "pause")
	asdf.Command = append(asdf.Command, "true")

	mpvCmd, err := json.Marshal(asdf)
	if err != nil {
		return "", err
	}
	return string(mpvCmd), nil
}
