package messages

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/nathan-hello/nat-sync/src/messages/ack"
	"github.com/nathan-hello/nat-sync/src/messages/commands"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type Message interface {
	ToBits() ([]byte, error)
}

func New[T string | []byte](i T) (Message, error) {
	switch t := any(i).(type) {
	case []byte:
		if commands.IsCommand(t) {
			utils.DebugLogger.Printf("cmd hit")
			c, err := commands.New(t)
			utils.DebugLogger.Printf("cmd got: %#v\n", c)
			return c, err
		}
		if ack.IsAck(t) {
			return ack.New(t)
		}
		return nil, utils.ErrBadMsgType(t)
	case string:
		utils.DebugLogger.Printf("str hit")
		cmd, err := commands.New(t)
		if err != nil {
			return nil, utils.ErrBadString(t, err)
		}
		utils.DebugLogger.Printf("cmd got: %#v\n", cmd)
		return cmd, nil
	}
	return nil, utils.ErrImpossible
}

func WaitReader(reader io.Reader) (Message, error) {

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

	// concat the two slices, putting length in the beginning
	m := append(append([]byte{}, lengthBytes...), message...)
	asdf, err := New(m)
	if err != nil {
		return nil, err
	}

	return asdf, nil
}
