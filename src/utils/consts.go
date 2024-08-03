package utils

const CurrentVersion = 1

type LocalTarget string

const (
	TargetTest LocalTarget = ""
	TargetMpv  LocalTarget = "mpv"
	TargetVlc  LocalTarget = "vlc"
)
