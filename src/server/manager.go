package server

import (
	"sync"

	"github.com/nathan-hello/nat-sync/src/utils"
)

type Manager struct {
	Lock  sync.Mutex
	Rooms map[int64]utils.ServerRoom
}

func (m *Manager) BroadcastMessage(room int64, bits []byte) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	currentRoom, ok := m.Rooms[room]
	if !ok {
		utils.ErrorLogger.Printf("tried to broadcast message to nonexistent room. room requested: %d, room map: %#v\n", room, m.Rooms)
		return
	}

	for _, c := range currentRoom.Clients {
		i, err := c.Conn.Write(bits)
		if err != nil {
			utils.ErrorLogger.Printf("conn.Write for client failed. err: %s, client: %#v\n", err, c)
		}
		utils.DebugLogger.Printf("conn.Write wrote %d bytes for client: %#v\n", i, c)

	}
}

func (m *Manager) AddClient(room int64, c utils.Client) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	currentRoom, ok := m.Rooms[room]
	if !ok {
		utils.ErrorLogger.Printf("tried to add client to nonexistent room. room requested: %d, room map: %#v\n", room, m.Rooms)
		return
	}
	_, ok = currentRoom.Clients[c.Id]
	if ok {
		utils.ErrorLogger.Printf("tried to add client but it is already in the room. room.Clients: %#v\n, client struct: %#v\n", currentRoom.Clients, m.Rooms)
		return
	}

	currentRoom.Clients[c.Id] = c
}

func (m *Manager) RemoveClient(room int64, c utils.Client) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	currentRoom, ok := m.Rooms[room]
	if !ok {
		utils.ErrorLogger.Printf("tried to delete client from nonexistent room. room requested: %d, room map: %#v\n", room, m.Rooms)
		return
	}

	_, ok = currentRoom.Clients[c.Id]
	if ok {
		utils.ErrorLogger.Printf("tried to delete client but it is not in the room. room.Clients: %#v\n, client struct: %#v\n", currentRoom.Clients, m.Rooms)
		return
	}

	c.Conn.Close()

	delete(currentRoom.Clients, c.Id)
}
