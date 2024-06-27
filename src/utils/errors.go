package utils

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrInvalidArg = errors.New("command was given invalid arguments")
	ErrInvalidCmd = errors.New("command given to server is not valid")
	ErrNoContent  = errors.New("command given has no content")
)

func ByteEncodingErr(buf bytes.Buffer) error {
	return fmt.Errorf("could not encode the following bytes: %b\nstring representation: %s", buf, buf.String())
}
