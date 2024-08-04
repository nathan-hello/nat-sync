package impl

import (
	"bytes"
	"context"
	"encoding/binary"
	"strings"

	"github.com/nathan-hello/nat-sync/src/db"
	"github.com/nathan-hello/nat-sync/src/utils"
)

var RoomHead = struct {
	Create  uint16
	Destroy uint16
	KickAll uint16
}{
	Create:  1,
	Destroy: 2,
	KickAll: 3,
}

var RoomExecuteResponses = struct {
	Ok      []byte
	KickAll []byte
}{
	Ok:      []byte("ok"),
	KickAll: []byte("kickall"),
}

type Room struct {
	NameLen  uint16
	Name     string
	PassLen  uint16
	Password string
	Head     uint16
}

func (c *Room) GetHead() string { return "connect" }
func (c *Room) ToBits() ([]byte, error) {
	var bits = new(bytes.Buffer)

	nl := len(c.Name)
	if nl > 65535 {
		return nil, utils.ErrBadArgs([]string{"name too long", c.Name})
	}
	pl := len(c.Password)
	if pl > 65535 {
		return nil, utils.ErrBadArgs([]string{"password too long", c.Password})
	}

	c.NameLen = uint16(len(c.Name))
	if err := binary.Write(bits, binary.BigEndian, c.NameLen); err != nil {
		return nil, err
	}
	if err := binary.Write(bits, binary.BigEndian, []byte(c.Name)); err != nil {
		return nil, err
	}

	c.PassLen = uint16(len(c.Password))
	if err := binary.Write(bits, binary.BigEndian, c.PassLen); err != nil {
		return nil, err
	}
	if err := binary.Write(bits, binary.BigEndian, []byte(c.Password)); err != nil {
		return nil, err
	}

	if err := binary.Write(bits, binary.BigEndian, c.Head); err != nil {
		return nil, err
	}
	return bits.Bytes(), nil
}

func (c *Room) New(t any) error {
	switch s := t.(type) {
	case []byte:
		return c.newFromBits(s)
	case []string:
		return c.newFromString(s)

	default:
		return utils.ErrBadType
	}
}

func (c *Room) newFromBits(bits []byte) error {
	buf := bytes.NewReader(bits)

	// room name
	if err := binary.Read(buf, binary.BigEndian, &c.NameLen); err != nil {
		return err
	}
	nameBits := make([]byte, c.NameLen)
	if _, err := buf.Read(nameBits); err != nil {
		return err
	}
	c.Name = string(nameBits)

	// password
	if err := binary.Read(buf, binary.BigEndian, &c.PassLen); err != nil {
		return err
	}
	passBits := make([]byte, c.PassLen)
	if _, err := buf.Read(passBits); err != nil {
		return err
	}
	c.Password = string(passBits)

	if err := binary.Read(buf, binary.BigEndian, &c.Head); err != nil {
		return err
	}

	return nil
}

func (c *Room) Execute() ([]byte, error) {
	d := db.Db()
	ctx := context.Background()

	switch c.Head {
	case RoomHead.Create:
		if err := d.InsertRoom(ctx, db.InsertRoomParams{Name: c.Name, Password: c.Password}); err != nil {
			return nil, err
		}
		return RoomExecuteResponses.Ok, nil
	case RoomHead.Destroy:
		room, err := d.SelectRoomByNameWithPassword(ctx, c.Name)
		if err != nil {
			return nil, err
		}
		if room.Password != c.Password {
			return nil, utils.ErrBadArgs([]string{"bad password to room", c.Name})
		}
		if err := d.DeleteRoom(ctx, room.ID); err != nil {
			return nil, err
		}
		return RoomExecuteResponses.Ok, nil
	case RoomHead.KickAll:
		return RoomExecuteResponses.KickAll, nil
	}

	return nil, nil
}

// Example:
// ["name=nate"]
func (c *Room) newFromString(s []string) error {
	for _, v := range s {
		v = strings.TrimPrefix(v, "-")
		v = strings.TrimPrefix(v, "-")
		switch {
		case strings.HasPrefix(strings.ToLower(v), "name="):
			parts := strings.Split(v, "=")
			if len(parts) < 1 {
				return utils.ErrBadArgs(s)
			}
			name := parts[1]
			nameLen := len(name)
			c.Name = name
			c.NameLen = uint16(nameLen)
			if c.Name == "" {
				return utils.ErrBadArgs(s)
			}
		case strings.HasPrefix(strings.ToLower(v), "password="):
			parts := strings.Split(v, "=")
			if len(parts) < 1 {
				return utils.ErrBadArgs(s)
			}
			pass := parts[1]
			passLen := len(pass)
			c.Password = pass
			c.PassLen = uint16(passLen)
			if c.Password == "" {
				return utils.ErrBadArgs(s)
			}
		case strings.HasPrefix(strings.ToLower(v), "action="):
			v = strings.ToLower(v)
			flag := strings.TrimPrefix(v, "action=")
			if flag == "create" {
				c.Head = RoomHead.Create
			}
			if flag == "destroy" {
				c.Head = RoomHead.Destroy
			}
			if flag == "kickall" {
				c.Head = RoomHead.KickAll
			}
			if c.Head == 0 {
				return utils.ErrBadArgs(s)
			}
		}
	}
	if c.Name == "" || c.Head == 0 {
		return utils.ErrBadArgs(s)
	}
	return nil
}
