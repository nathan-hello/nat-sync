package commands

const CurrentVersion = 1

type CmdHead uint8

const (
	AppendHead CmdHead = 1
	ChangeHead CmdHead = 2
	KickHead   CmdHead = 3
	JoinHead   CmdHead = 4
	PauseHead  CmdHead = 5
	PlayHead   CmdHead = 8
	SeekHead   CmdHead = 7
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
