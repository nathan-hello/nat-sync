package impl

import "github.com/nathan-hello/nat-sync/src/utils"

type MpvJson struct {
	Command []string `json:"command"`
}

type Command interface {
	New(any) error
	ToBits() ([]byte, error)
	GetHead() string
}

type PlayerCommand interface {
	New(any) error
	ToBits() ([]byte, error)
	GetHead() string
	ToPlayer(p utils.LocalTarget) ([]byte, error)
	IsPlayer() bool
}

type AdminCommand interface {
	New(any) error
	ToBits() ([]byte, error)
	GetHead() string
	IsAdmin() bool
}

type ServerCommand interface {
	New(any) error
	ToBits() ([]byte, error)
	GetHead() string
	Execute(executor interface{}) ([]byte, error)
	IsServer() bool
}

type registeredHead struct {
	Code uint16
	Name string
	Impl Command
}

var registeredPlayers = []struct {
	Code uint16
	Name string
	Impl PlayerCommand
}{
	{1, "change", &Change{}},
	{2, "pause", &Pause{}},
	{3, "play", &Play{}},
	{4, "seek", &Seek{}},
	{5, "stop", &Stop{}},
	{6, "quit", &Quit{}},
}

var registeredAdmins = []struct {
	Code uint16
	Name string
	Impl AdminCommand
}{
	{100, "ack", &Ack{}},
	{101, "accept", &Accept{}},
	{102, "join", &Join{}},
}
var registeredServers = []struct {
	Code uint16
	Name string
	Impl ServerCommand
}{
	{200, "wait", &Wait{}},
	{104, "kick", &Kick{}},
}

func RegisteredCmds() []registeredHead {
	var RegisteredHeads = []registeredHead{}

	for _, p := range registeredPlayers {
		RegisteredHeads = append(RegisteredHeads, registeredHead{
			Code: p.Code,
			Name: p.Name,
			Impl: p.Impl,
		})
	}

	for _, a := range registeredAdmins {
		RegisteredHeads = append(RegisteredHeads, registeredHead{
			Code: a.Code,
			Name: a.Name,
			Impl: a.Impl,
		})
	}

	for _, s := range registeredServers {
		RegisteredHeads = append(RegisteredHeads, registeredHead{
			Code: s.Code,
			Name: s.Name,
			Impl: s.Impl,
		})
	}

	return RegisteredHeads
}
