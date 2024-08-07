package utils

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidArg    = errors.New("command was given invalid arguments")
	ErrInvalidCmd    = errors.New("command given to server is not valid")
	ErrTextNoContent = errors.New("text given has no content")
	ErrUsernameBlank = errors.New("username string in command was blank")
	ErrUsernameLong  = errors.New("username is too long. must be <255 characters")
	ErrNoRoomClient  = errors.New("room id was not specified. you must join a room first")
	ErrNoRoomServer  = errors.New("room id was not specified. server tried sending a command to room 0")
)

var (
	ErrImpossible        = errors.New("this error is here because go's type system is just ok")
	ErrPlayerAlreadyDead = errors.New("player is already dead")
	ErrBadType           = errors.New("new function got a type that was not []byte or []string")
)

func ErrDecodeType(bits []byte) error {
	return fmt.Errorf("could not decode the type from the following bytes: %b", bits)
}

func ErrFixedContentLength(cmd interface{}, subCmd interface{}) error {
	return fmt.Errorf("content is greater than 65535 bits, full command: %#v, content: %#v", cmd, subCmd)
}

func ErrNoCmdHeadFound(head uint8) error {
	return fmt.Errorf("cmd head not found from bit: %b", head)
}

func ErrBadArgs(s []string) error {
	return fmt.Errorf("err parsing cmd string: %s", s)
}

func ErrRequiredArgs(msg string) error {
	return fmt.Errorf("err cmd is missing a required arg: %s", msg)
}

func ErrNoArgs(msg string) error {
	return fmt.Errorf("err cmd was given no args but requires them: %s", msg)
}

func ErrLongString(s []string) error {
	return fmt.Errorf("err arg is longer than 65535: %s", s)
}

func ErrNotImplemented(s string) error {
	return fmt.Errorf("feature is not implemented yet: %s", s)
}

func ErrBadString(s string, err error) error {
	e := fmt.Errorf("string was not able to be turned into a message. string: %s", s)
	return errors.Join(e, err)
}

func ErrBadMsgType(b []byte) error {
	return fmt.Errorf("type bit does not match anything we got: %b", b)
}

func ErrTimeout(s string) error {
	return fmt.Errorf("connection to %s timedout", s)
}
