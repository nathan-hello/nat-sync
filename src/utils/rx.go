package utils

import (
	"encoding/binary"
	"fmt"
	"io"
)

func ConnRxCommand(reader io.Reader) ([]byte, error) {

	lengthBytes := make([]byte, 2)

	// this is where receieve blocks until new trasmission is heard
	// because it is a read in the gopher's language, it advances the reader two bytes
	// we will add those bytes back in just a second
	_, err := io.ReadFull(reader, lengthBytes)
	if err != nil {
		return nil, fmt.Errorf("connection closed or error reading length bytes: %#v", lengthBytes)
	}

	length := binary.BigEndian.Uint16(lengthBytes)

	message := make([]byte, length)

	n, err := io.ReadFull(reader, message)
	if err != nil || uint16(n) != length {
		return nil, fmt.Errorf("bytes read: %d, expected: %d", n, length)
	}

	msgWithLen := []byte{}
	msgWithLen = append(msgWithLen, lengthBytes...)
	msgWithLen = append(msgWithLen, message...)

	return msgWithLen, nil
}
