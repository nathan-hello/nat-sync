package server

import (
	"net"
	"sync"
)

type Client struct {
	Name string
	Id   uint16
	Conn net.Conn
	Room uint8
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

func (m *Manager) RemoveClient(room uint8, userId uint16) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	for _, v := range m.Clients {
		if v.Id == userId && v.Room == room {
			v.Conn.Close()
		}
	}
}
