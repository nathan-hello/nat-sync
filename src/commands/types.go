package commands

const CurrentVersion = 1

type CmdHead uint8

const (
	ChangeHead CmdHead = iota
	KickHead   CmdHead = iota
	JoinHead   CmdHead = iota
	PauseHead  CmdHead = iota
	PlayHead   CmdHead = iota
	SeekHead   CmdHead = iota
)

type Command struct {
	Head    CmdHead
	Version uint16
	UserId  uint16
	Content []byte
	Sub     SubCommand
}

type SubCommand interface {
	FromString(s []string) error
	FromBits(bits []byte) error
	ToBits() ([]byte, error)
}
