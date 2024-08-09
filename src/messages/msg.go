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
	Sub     impl.Command
}

func NewMulti(bits []byte) ([]Message, error) {
	msgs := []Message{}
	s := string(bits)

	delmited := strings.Split(s, ";")
	utils.DebugLogger.Printf("delimited: %#v\n", delmited)

	for _, v := range delmited {
		if v == "" {
			continue
		}
		msg := Message{}

		err := msg.TextUnmarshaller([]byte(v))
		if err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func NewFromSub(c impl.Command, roomId int64) Message {
	head, _ := getHeadFromString(c.GetHead())
	return Message{
		Head:    head,
		Version: utils.CurrentVersion,
		RoomId:  roomId,
		Sub:     c,
	}
}

func (msg *Message) TextUnmarshaller(text []byte) error {
	s := string(text)
	if s == "" {
		return utils.ErrTextNoContent
	}
	parts := strings.Fields(s)

	roomId, err := strconv.ParseInt(parts[0], 16, 64)
	if err != nil {
		return err
	}

	head, err := getHeadFromString(parts[1])
	if err != nil {
		return err
	}

	sub, err := getSubFromHead(head)
	if err != nil {
		return err
	}

	err = sub.New(parts[2:])
	if err != nil {
		return err
	}

	msg.Head = head
	msg.Version = utils.CurrentVersion
	msg.RoomId = roomId
	msg.Sub = sub

	return nil
}

// WaitReader blocks until a new message of the appropiate format is received.
// If a byte stream is read but is not of the format Message.MarshalBinary() provides,
// this will hang indefinitely.
func WaitReader(reader io.Reader) (Message, error) {
	msg := Message{}

	lengthBytes := make([]byte, 2)

	// Block until a new message is read, then advance reader two bytes.
	// We will add those bytes back in just a second.
	_, err := io.ReadFull(reader, lengthBytes)
	if err != nil {
		if err == io.EOF {
			return msg, err
		}
		return msg, fmt.Errorf("reader in waitreader failed. io.ReadFull err: %w", err)
	}

	length := binary.BigEndian.Uint16(lengthBytes)

	message := make([]byte, length)

	n, err := io.ReadFull(reader, message)
	if err != nil || uint16(n) != length {
		return msg, fmt.Errorf("bytes read: %d, expected: %d", n, length)
	}

	// concat the two slices, putting length in the beginning
	m := append(lengthBytes, message...)
	err = msg.UnmarshalBinary(m)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func (cmd *Message) MarshalBinary() ([]byte, error) {
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

func (msg *Message) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Read the fixed-length part of the Command struct
	if err := binary.Read(buf, binary.BigEndian, &msg.Length); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Length):", err)
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.RoomId); err != nil {
		utils.DebugLogger.Println("binary.Read failed (RoomId):", err)
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.Head); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Head):", err)
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.Version); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Version):", err)
		return err
	}

	// Read the remaining bytes into Content
	msg.Content = make([]byte, buf.Len())
	if err := binary.Read(buf, binary.BigEndian, &msg.Content); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Content):", err)
		return err
	}

	sub, err := getSubFromHead(msg.Head)
	if err != nil {
		return err
	}

	sub.New(msg.Content)
	msg.Sub = sub

	utils.DebugLogger.Printf("decoded struct: %#v\n", msg)

	return nil
}

func getSubFromHead(head uint16) (impl.Command, error) {
	for _, v := range impl.RegisteredCmds() {
		if v.Code == head {
			return v.Impl, nil
		}
	}
	return nil, utils.ErrNoCmdHeadFound(uint8(head))
}

func getHeadFromString(s string) (uint16, error) {
	for _, v := range impl.RegisteredCmds() {
		if v.Name == s {
			return v.Code, nil
		}
	}
	return 0, utils.ErrBadString(s, nil)
}
