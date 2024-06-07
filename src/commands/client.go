package commands

import (
	"crypto/tls"
	"fmt"
	"time"
)

type Connection struct {
	Conn     tls.Conn
	CmdQueue chan string
}

func (c *Connection) Send(command ClientCommand) error { return nil }

type ClientCommand interface {
	Raw() string
	Parse() string
}

type Seek struct {
	CurrentTime time.Time
	NewTime     time.Time
}

func (s *Seek) Raw() string {
	return fmt.Sprintf("%#v\n")
}

func (s *Seek) Parse() string {
	diff := s.NewTime.Sub(s.CurrentTime)
	return fmt.Sprintf(
		"%02dh:%02dm%02ds",
		diff.Hours(),
		diff.Minutes(),
		diff.Seconds(),
	)
}
