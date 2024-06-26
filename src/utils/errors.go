package utils

import "errors"

var (
	ErrInvalidArg = errors.New("command was given invalid arguments")
	ErrInvalidCmd = errors.New("command given to server is not valid")
)
