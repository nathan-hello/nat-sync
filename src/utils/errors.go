package utils

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidArg    = errors.New("command was given invalid arguments")
	ErrInvalidCmd    = errors.New("command given to server is not valid")
	ErrNoContent     = errors.New("command given has no content")
	ErrUsernameBlank = errors.New("username string in command was blank")
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
