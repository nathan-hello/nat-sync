package commands

import "encoding/binary"

const CurrentVersion = 0000_0001

const (
	ExitHead     = 1
	SeekHead     = 2
	PauseHead    = 3
	PlayHead     = 4
	NewVideoHead = 5
	JoinRoomHead = 6
)

type Command struct {
	Head    uint8
	Version uint16
	Creator [32]byte
	Content []byte
}

type SubCommand interface {
	FromBits(bits []byte)
	ToBits() []byte
}

type Seek struct {
	Hours uint8
	Mins  uint8
	Secs  uint8
}

type Pause struct{}
type Play struct{}
type NewVideo struct {
	UriLength uint16
	Uri       string
	Local     bool
}

func (c *Seek) ToBits() []byte  { return []byte{c.Hours, c.Mins, c.Secs} }
func (c *Play) ToBits() []byte  { return []byte{} }
func (c *Pause) ToBits() []byte { return []byte{} }
func (c *NewVideo) ToBits() []byte {
	var bits []byte
	var l byte
	if c.Local {
		l = 1
	} else {
		l = 0
	}

	binary.LittleEndian.PutUint16(bits, c.UriLength)
	bits = append(bits, []byte(c.Uri)...)
	bits = append(bits, l)
	return bits

}

func (c *Seek) FromBits(bits []byte)     { c.Hours = bits[0]; c.Mins = bits[1]; c.Secs = bits[2] }
func (c *Play) FromBits(bits []byte)     {}
func (c *Pause) FromBits(bits []byte)    {}
func (c *NewVideo) FromBits(bits []byte) {}
