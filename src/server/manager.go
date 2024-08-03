package server

import (
	"net"
	"sync"
)

type Room struct {
	Id       uint8
	Name     string
	Password string
}

type Client struct {
	Id     int64
	Name   string
	Conn   net.Conn
	RoomId int64
}

type Manager struct {
	Lock    sync.Mutex
	Clients []Client
}

func (m *Manager) BroadcastMessage(bits []byte) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	for _, c := range m.Clients {
		c.Conn.Write(bits)
	}
}

func (m *Manager) AddClient(c Client) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Clients = append(m.Clients, c)
}

func (m *Manager) RemoveClient(room int64, userId int64) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	for _, v := range m.Clients {
		if v.Id == userId && v.RoomId == room {
			v.Conn.Close()
		}
	}
}
