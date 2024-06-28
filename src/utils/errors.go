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
