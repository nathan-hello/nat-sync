package utils

const CurrentVersion = 1

type MsgType uint8

const (
	MsgCommand MsgType = iota
)

type LocalTarget string

const (
	TargetTest LocalTarget = ""
	TargetMpv  LocalTarget = "mpv"
	TargetVlc  LocalTarget = "vlc"
)
