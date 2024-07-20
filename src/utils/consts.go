package utils

const CurrentVersion = 1

type MsgType uint8

const (
	MsgCommand MsgType = iota
	MsgAck     MsgType = iota
	MsgLocal   MsgType = iota
)

type PlayerTargets string

const (
	TargetTest PlayerTargets = ""
	TargetMpv  PlayerTargets = "mpv"
	TargetVlc  PlayerTargets = "vlc"
)
