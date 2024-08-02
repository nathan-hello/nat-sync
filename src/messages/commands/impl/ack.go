package impl

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/nathan-hello/nat-sync/src/client/players"
	"github.com/nathan-hello/nat-sync/src/utils"
)

type ackHead uint16

var AckCode = struct {
	Ok                   ackHead
	InternalServiceError ackHead
	BadEcho              ackHead
}{
	Ok:                   200,
	InternalServiceError: 500,
	BadEcho:              601,
}

type Ack struct {
	Code    ackHead
	MsgLen  uint32
	Message string
}

func (c *Ack) ExecuteClient(_ players.Player) ([]byte, error) {
	utils.DebugLogger.Printf("received ack: %#v\n", c)
	return nil, nil
}
func (c *Ack) ExecuteServer() ([]byte, error) { return nil, nil }
func (c *Ack) IsEchoed() bool                 { return false }

func (a *Ack) ToBits() ([]byte, error) {
	bits := new(bytes.Buffer)
	if err := binary.Write(bits, binary.BigEndian, a.Code); err != nil {
		return nil, err
	}

	return bits.Bytes(), nil

}

func (c *Ack) NewFromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	// Read the fixed-length part of the Command struct
	var ack Ack
	if err := binary.Read(buf, binary.BigEndian, &ack.Code); err != nil {
		utils.DebugLogger.Println("binary.Read failed (Version):", err)
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &c.MsgLen); err != nil {
		return err
	}

	msgBits := make([]byte, c.MsgLen)
	if _, err := buf.Read(msgBits); err != nil {
		return err
	}
	c.Message = string(msgBits)

	return nil
}

func (c *Ack) NewFromString(s []string) error {
	for _, v := range s {
		v = strings.ToLower(v)
		v = strings.TrimPrefix(v, "-")
		v = strings.TrimPrefix(v, "-")
		switch {
		case strings.HasPrefix(v, "message="):
			flag, _ := strings.CutPrefix(v, "message=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")
			c.Message = flag
			c.MsgLen = uint32(len(flag))
		case strings.HasPrefix(v, "code="):
			flag, _ := strings.CutPrefix(v, "code=")
			flag, _ = strings.CutPrefix(flag, "\"")
			flag, _ = strings.CutSuffix(flag, "\"")

			switch flag {
			case "ok":
				c.Code = AckCode.Ok
			case "err":
				c.Code = AckCode.InternalServiceError
			default:
				return utils.ErrBadArgs(s)
			}
		}
	}
	return nil
}
