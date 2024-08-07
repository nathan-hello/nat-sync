package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/nathan-hello/nat-sync/src/messages/impl"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type Message struct {
	Length  uint16
	RoomId  int64
	Head    uint16
	Version uint16
	Content []byte
	Sub     Command
}

type Command interface {
	New(any) error
	ToBits() ([]byte, error)
	GetHead() string
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
	Execute(executor interface{}) ([]byte, error)
}

type RegisteredHead struct {
	Code uint16
	Name string
	Impl Command
}

var RegisteredHeads = []RegisteredHead{
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

func New(i any, roomId *int64) ([]Message, error) {

	msgs := []Message{}
	switch t := i.(type) {
	case []byte:
		m, err := newMsgFromBits(t)
		if err != nil {
			return nil, err
		}
		return append(msgs, *m), nil

	case string:
		utils.DebugLogger.Printf("string got: %s\n", t)
		delmited := strings.Split(t, ";")
		rId, err := getRoomIdFromString(&[]string{delmited[0]})
		if err != nil {
			return nil, err
		}

		for _, s := range delmited[1:] {

			m, err := newMsgFromString(s, &rId)
			if err != nil {
				return nil, utils.ErrBadString(t, err)
			}

			// if there is a final ; but no commands afterwards,
			// msg will be nil because of the if len(parts) == 0 {return nil, nil}
			if m == nil {
				continue
			}
			msgs = append(msgs, *m)
		}
		return msgs, nil

	case Command:
		// Placing a Command here assumes that there is a roomId passed
		// to this function. There is no error returned because when I
		// assume that it's done correctly.
		h, _ := getHeadFromString(t.GetHead())
		subBits, _ := t.ToBits()
		m := Message{
			Head:    h,
			RoomId:  *roomId,
			Version: utils.CurrentVersion,
			Sub:     t,
			Content: subBits,
		}

		return append(msgs, m), nil

	}
	return nil, utils.ErrImpossible
}

// Returns a *Command without UserId field
func newMsgFromString(s string, roomId *int64) (*Message, error) {
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
		RoomId:  *roomId,
		Version: utils.CurrentVersion,
		Sub:     sub,
		Content: content,
	}, nil

}

func getSubFromHead(head uint16) (Command, error) {
	for _, v := range RegisteredHeads {
		if v.Code == head {
			return v.Impl, nil
		}
	}
	return nil, utils.ErrNoCmdHeadFound(uint8(head))
}

func getHeadFromString(s string) (uint16, error) {
	for _, v := range RegisteredHeads {
		if v.Name == s {
			return v.Code, nil
		}
	}
	return 0, utils.ErrBadString(s, nil)
}

func getRoomIdFromString(s *[]string) (int64, error) {
	for i, v := range *s {
		v = strings.ToLower(v)
		v = strings.TrimPrefix(v, "-")
		v = strings.TrimPrefix(v, "-")
		if strings.HasPrefix(v, "roomid=") {
			flag := strings.TrimPrefix(v, "roomid=")
			num, err := strconv.ParseInt(flag, 10, 64)
			if err != nil {
				return -1, err
			}

			*s = append((*s)[:i], (*s)[i+1:]...)
			return num, nil
		}
	}

	return -1, utils.ErrNoRoomClient
}

// WaitReader blocks until a new message of the appropiate format is received.
// If a byte stream is read but is not of the format Message.ToBits() provides,
// this will hang indefinitely.
func WaitReader(reader io.Reader) ([]Message, error) {

	lengthBytes := make([]byte, 2)

	// Block until a new message is read, then advance reader two bytes.
	// We will add those bytes back in just a second.
	_, err := io.ReadFull(reader, lengthBytes)
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("reader in waitreader failed. io.ReadFull err: %w", err)
	}

	length := binary.BigEndian.Uint16(lengthBytes)

	message := make([]byte, length)

	n, err := io.ReadFull(reader, message)
	if err != nil || uint16(n) != length {
		return nil, fmt.Errorf("bytes read: %d, expected: %d", n, length)
	}

	// concat the two slices, putting length in the beginning
	m := append(append([]byte{}, lengthBytes...), message...)
	asdf, err := New(m, nil)
	if err != nil {
		return nil, err
	}

	return asdf, nil
}

func (cmd *Message) ToBits() ([]byte, error) {
	bits := new(bytes.Buffer)

	if cmd.Sub != nil {
		cmd.Sub = nil
	}
	if cmd.Version == 0 {
		cmd.Version = utils.CurrentVersion
	}

	if err := binary.Write(bits, binary.BigEndian, cmd.RoomId); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, cmd.Head); err != nil {
		return nil, err
	}
	if err := binary.Write(bits, binary.BigEndian, cmd.Version); err != nil {
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
	if err := binary.Read(buf, binary.BigEndian, &msg.RoomId); err != nil {
		utils.DebugLogger.Println("binary.Read failed (RoomId):", err)
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
