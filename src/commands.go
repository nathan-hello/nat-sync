package src

import (
	"fmt"
	"time"
)

const CurrentVersion = "0.0.1"

func MarshalErrToJSON(raw any, err error) string {
	return fmt.Sprintf(`{error: "%s", raw_command: "%s" }`, err.Error(), raw)

}

type Command interface {
	Render() (RenderedCommand, error)
}

type RenderedCommand struct {
	Command string      `json:"command"`
	Version string      `json:"version"`
	Content interface{} `json:"content"`
}

type Seek struct {
	Creator  string `json:"creator"`
	Location string `json:"difference"`
}

func (c *Seek) Render() (*RenderedCommand, error) {
	t, err := time.ParseDuration(c.Location)
	if err != nil {
		return nil, err
	}
	c.Location = fmt.Sprintf("%.0f", t.Seconds())
	return &RenderedCommand{Command: "seek", Version: CurrentVersion, Content: c}, nil
}

type Pause struct {
	Creator string `json:"creator"`
}

func (c *Pause) Render() (*RenderedCommand, error) {
	return &RenderedCommand{Command: "pause", Version: CurrentVersion, Content: c}, nil
}

type Play struct {
	Creator string `json:"creator"`
}

func (c *Play) Render() (*RenderedCommand, error) {
	return &RenderedCommand{Command: "play", Version: CurrentVersion, Content: c}, nil
}

type NewVideo struct {
	Creator string `json:"creator"`
	Uri     string `json:"uri"`
	Local   bool   `json:"local"`
}
