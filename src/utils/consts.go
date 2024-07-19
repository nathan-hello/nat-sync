package utils

const CurrentVersion = 1

type MsgType uint8

const (
	MsgCommand MsgType = iota
	MsgAck     MsgType = iota
)
