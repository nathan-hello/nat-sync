package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/nathan-hello/nat-sync/src/messages/impl"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type Message struct {
	Length  uint16
	Head    uint16
	Version uint16
	UserId  uint16
	Content []byte
	Sub     Command
}

type Command interface {
	New(any) error
	ToBits() ([]byte, error)
}

type PlayerCommand interface {
	New(any) error
	ToBits() ([]byte, error)
	ToPlayer(p utils.LocalTarget) ([]byte, error)
}

type AdminCommand interface {
	New(any) error
	ToBits() ([]byte, error)
}

type ServerCommand interface {
	New(any) error
	ToBits() ([]byte, error)
	Execute() ([]byte, error)
}

type RegisteredHead struct {
	Code uint16
	Name string
	Impl Command
}

var registeredHeads = []RegisteredHead{
	{1, "change", &impl.Change{}},
	{2, "pause", &impl.Pause{}},
	{3, "play", &impl.Play{}},
	{4, "seek", &impl.Seek{}},
	{5, "stop", &impl.Stop{}},
	{6, "quit", &impl.Quit{}},
	{100, "ack", &impl.Ack{}},
	{101, "kick", &impl.Kick{}},
	{102, "join", &impl.Join{}},
	{200, "wait", &impl.Wait{}},
}

func New[T string | []byte](i T) ([]Message, error) {
	msgs := []Message{}
	switch t := any(i).(type) {
	case []byte:
		m, err := newMsgFromBits(t)
		if err != nil {
			return nil, err
		}
		return append(msgs, *m), nil

	case string:
		utils.DebugLogger.Printf("string got: %s\n", t)
		delmited := strings.Split(t, ";")
		for _, s := range delmited {
			msg, err := newMsgFromString(s)
			if err != nil {
				return nil, utils.ErrBadString(t, err)
			}

			// if there is a final ; but no commands afterwards,
			// msg will be nil because of the if len(parts) == 0 {return nil, nil}
			if msg != nil {
				msgs = append(msgs, *msg)
			}
		}
		return msgs, nil
	}
	return nil, utils.ErrImpossible
}

func (cmd *Message) ToBits() ([]byte, error) {
	bits := new(bytes.Buffer)

	if cmd.Sub != nil {
		cmd.Sub = nil
	}
	if cmd.Version == 0 {
		cmd.Version = utils.CurrentVersion
	}

	if err := binary.Write(bits, binary.BigEndian, cmd.Head); err != nil {
		return nil, err
	}
	if err := binary.Write(bits, binary.BigEndian, cmd.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(bits, binary.BigEndian, cmd.UserId); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, cmd.Content); err != nil {
		return nil, err
	}

	cmd.Length = uint16(len(bits.Bytes()))
	finalBits := new(bytes.Buffer)

	if err := binary.Write(finalBits, binary.BigEndian, cmd.Length); err != nil {
		return nil, err
	}

	if _, err := finalBits.Write(bits.Bytes()); err != nil {
		return nil, err
	}

	// utils.DebugLogger.Printf("decoded bytes: %b ", finalBits.Bytes())
	utils.DebugLogger.Printf("encoded struct: %#v\n", cmd)
	return finalBits.Bytes(), nil
}

func newMsgFromBits(bits []byte) (*Message, error) {
	buf := bytes.NewReader(bits)

	// Read the fixed-length part of the Command struct
	var msg Message
	if err := binary.Read(buf, binary.BigEndian, &msg.Length); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Length):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.Head); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Head):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.Version); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Version):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.UserId); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Creator):", err)
		return nil, err
	}

	// Read the remaining bytes into Content
	msg.Content = make([]byte, buf.Len())
	if err := binary.Read(buf, binary.BigEndian, &msg.Content); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Content):", err)
		return nil, err
	}

	sub, err := getSubFromHead(msg.Head)
	if err != nil {
		return nil, err
	}

	sub.New(msg.Content)
	msg.Sub = sub

	// utils.DebugLogger.Printf("decoded struct: %#v\n", msg)

	return &msg, nil
}

// Returns a *Command without UserId field
func newMsgFromString(s string) (*Message, error) {
	parts := strings.Fields(s)

	if len(parts) == 0 {
		return nil, nil
	}

	head, err := getHeadFromString(parts[0])
	if err != nil {
		return nil, err
	}

	sub, err := getSubFromHead(head)
	if err != nil {
		return nil, err
	}

	err = sub.New(parts[1:])
	if err != nil {
		return nil, err
	}

	content, err := sub.ToBits()
	if err != nil {
		return nil, err
	}

	return &Message{
		Head:    head,
		Version: utils.CurrentVersion,
		Sub:     sub,
		Content: content,
	}, nil

}

// Register new commands here
func getSubFromHead(head uint16) (Command, error) {
	for _, v := range registeredHeads {
		if v.Code == head {
			return v.Impl, nil
		}
	}
	return nil, utils.ErrNoCmdHeadFound(uint8(head))
}

// Register new strings here
func getHeadFromString(s string) (uint16, error) {
	for _, v := range registeredHeads {
		if v.Name == s {
			return v.Code, nil
		}
	}
	return 0, utils.ErrBadString(s, nil)

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
