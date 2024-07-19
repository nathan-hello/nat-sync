package ack

import (
	"bytes"
	"encoding/binary"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type ackHead uint16
type ackCode struct {
	Ok                   ackHead
	InternalServiceError ackHead
	BadEcho              ackHead
}

var AckCode = ackCode{
	Ok:                   200,
	InternalServiceError: 500,
	BadEcho:              601,
}

type Ack struct {
	Length  uint16
	Type    utils.MsgType
	Head    ackHead
	Version uint16
	Content []byte
}

func New[T []byte | string](i T) (*Ack, error) {
	switch t := any(i).(type) {
	case []byte:
		ack, err := newAckFromBits(t)
		if err != nil {
			return nil, err
		}
		return ack, nil
	case string:
		ack := newAckFromString(t)
		return ack, nil
	default:
		return nil, utils.ErrImpossible
	}
}
func newAck(head ackHead, content []byte) *Ack {
	return &Ack{
		Type:    utils.MsgAck,
		Head:    head,
		Version: utils.CurrentVersion,
		Content: content,
	}

}

func (a *Ack) ToBits() ([]byte, error) {

	bits := new(bytes.Buffer)

	if a.Version == 0 {
		a.Version = utils.CurrentVersion
	}
	if a.Type == 0 {
		a.Type = utils.MsgAck
	}

	err := binary.Write(bits, binary.BigEndian, a.Type)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.BigEndian, a.Head)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.BigEndian, a.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bits, binary.BigEndian, a.Content)
	if err != nil {
		return nil, err
	}

	a.Length = uint16(len(bits.Bytes()))
	finalBits := new(bytes.Buffer)

	err = binary.Write(finalBits, binary.BigEndian, a.Length)
	if err != nil {
		return nil, err
	}

	_, err = finalBits.Write(bits.Bytes())
	if err != nil {
		return nil, err
	}

	// utils.DebugLogger.Printf("decoded bytes: %b ", finalBits.Bytes())
	return finalBits.Bytes(), nil
}

func newAckFromBits(bits []byte) (*Ack, error) {
	buf := bytes.NewReader(bits)

	// Read the fixed-length part of the Command struct
	var ack Ack
	if err := binary.Read(buf, binary.BigEndian, &ack.Length); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Length):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &ack.Type); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Type):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &ack.Head); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Head):", err)
		return nil, err
	}
	if err := binary.Read(buf, binary.BigEndian, &ack.Version); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Version):", err)
		return nil, err
	}

	// Read the remaining bytes into Content
	ack.Content = make([]byte, buf.Len())
	if err := binary.Read(buf, binary.BigEndian, &ack.Content); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Content):", err)
		return nil, err
	}

	return &ack, nil
}

func newAckFromString(s string) *Ack {
	switch s {
	case "200":
		return newAck(AckCode.Ok, nil)
	case "500":
		return newAck(AckCode.InternalServiceError, nil)
	default:
		return newAck(AckCode.InternalServiceError, []byte("malformed ack: "+s))
	}

}

func IsAck(bits []byte) bool {
	var ack Ack
	buf := bytes.NewReader(bits)

	_ = binary.Read(buf, binary.BigEndian, &ack.Length)
	_ = binary.Read(buf, binary.BigEndian, &ack.Type)

	return ack.Type == utils.MsgAck
}
