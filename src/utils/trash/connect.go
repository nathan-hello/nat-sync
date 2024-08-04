// package impl
//
// import (
//
//	"bytes"
//	"encoding/binary"
//	"strconv"
//	"strings"
//
//	"github.com/nathan-hello/nat-sync/src/utils"
//
// )
//
//	type Connect struct {
//		NameLen uint8
//		Name    string
//		RoomId  int64
//	}
//
// func (c *Connect) Execute(man interface{ AddClient })
//
// func (c *Connect) GetHead() string { return "connect" }
//
//	func (c *Connect) ToBits() ([]byte, error) {
//		var bits = new(bytes.Buffer)
//
//		if err := binary.Write(bits, binary.BigEndian, c.NameLen); err != nil {
//			return nil, err
//		}
//		if err := binary.Write(bits, binary.BigEndian, []byte(c.Name)); err != nil {
//			return nil, err
//		}
//		if err := binary.Write(bits, binary.BigEndian, c.RoomId); err != nil {
//			return nil, err
//		}
//
//		return bits.Bytes(), nil
//	}
//
//	func (c *Connect) New(t any) error {
//		switch s := t.(type) {
//		case []byte:
//			return c.newFromBits(s)
//		case []string:
//			return c.newFromString(s)
//
//		default:
//			return utils.ErrBadType
//		}
//	}
//
//	func (c *Connect) newFromBits(bits []byte) error {
//		buf := bytes.NewReader(bits)
//
//		if err := binary.Read(buf, binary.BigEndian, &c.NameLen); err != nil {
//			return err
//		}
//
//		uriBits := make([]byte, c.NameLen)
//		if _, err := buf.Read(uriBits); err != nil {
//			return err
//		}
//		c.Name = string(uriBits)
//
//		return nil
//	}
//
// // Example:
// // ["name=nate"]
//
//	func (c *Connect) newFromString(s []string) error {
//		for _, v := range s {
//			v = strings.TrimPrefix(v, "-")
//			v = strings.TrimPrefix(v, "-")
//			switch {
//			case strings.HasPrefix(strings.ToLower(v), "name="):
//				parts := strings.Split(v, "=")
//				if len(parts) < 1 {
//					return utils.ErrBadArgs(s)
//				}
//				name := parts[1]
//				nameLen := len(name)
//				c.Name = name
//				c.NameLen = uint8(nameLen)
//			case strings.HasPrefix(v, "roomid="):
//				flag, _ := strings.CutPrefix(v, "roomid=")
//				flag, _ = strings.CutPrefix(flag, "\"")
//				flag, _ = strings.CutSuffix(flag, "\"")
//				i, err := strconv.ParseUint(flag, 10, 16)
//				if err != nil {
//					return utils.ErrBadArgs(s)
//				}
//				c.RoomId = int64(i)
//			}
//		}
//		if c.Name == "" {
//			return utils.ErrBadArgs(s)
//		}
//		return nil
//	}
package trash
