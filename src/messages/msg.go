package messages

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/nathan-hello/nat-sync/src/messages/commands"
	"github.com/nathan-hello/nat-sync/src/players"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type Message interface {
	ToBits() ([]byte, error)
	ExecutePlayer(players.Player) ([]byte, error)
}

func New[T string | []byte](i T) ([]Message, error) {
	msgs := []Message{}
	switch t := any(i).(type) {
	case []byte:
		if commands.IsCommand(t) {
			m, err := commands.New(t)
			if err != nil {
				return nil, err
			}
			return append(msgs, m), nil
		}
		return nil, utils.ErrBadMsgType(t)
	case string:

		utils.DebugLogger.Printf("string got: %s\n", t)
		if m := getMacro(t); m != nil {
			return m, nil
		}

		delmited := strings.Split(t, ";")
		for _, s := range delmited {
			cmd, err := commands.New(s)
			if err != nil {
				return nil, utils.ErrBadString(t, err)
			}

			msgs = append(msgs, cmd)
		}
		return msgs, nil
	}
	return nil, utils.ErrImpossible
}

func WaitReader(reader io.Reader) ([]Message, error) {

	lengthBytes := make([]byte, 2)

	// this is where receieve blocks until new trasmission is heard
	// because it is a read in the gopher's language, it advances the reader two bytes
	// we will add those bytes back in just a second
	_, err := io.ReadFull(reader, lengthBytes)
	if err == io.EOF {
		return nil, err
	}
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
