package src

import (
	"fmt"
	"time"
)

const CurrentVersion = "0.0.1"

func MarshalErrToJSON(raw any, err error) string {
	return fmt.Sprintf(`{error: "%s", raw_command: "%s" }`, err.Error(), raw)

}

type BaseCommand struct {
	Command string `json:"command"`
	Version string `json:"version"`
}

type RenderedCommand struct {
	BaseCommand
	Argument any `json:"arg"`
}

type Command interface {
	Render() error
}

type Seek struct {
	BaseCommand
	NewTime string `json:"new_time"`
}

func (s *Seek) Render() error {
	s.Command = "seek"
	s.Version = CurrentVersion
	t, err := time.ParseDuration(s.NewTime)
	if err != nil {
		return err
	}
	// s.NewTime = fmt.Sprintf("%fh:%fm:%fs", t.Hours(), t.Minutes(), t.Seconds())
	s.NewTime = fmt.Sprintf("%.0f", t.Seconds())
	return nil
}
