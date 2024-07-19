package server

import (
	"net"
	"sync"
)

type user struct {
	ip   net.Addr
	name string
	id   uint16
}

type manager struct {
	room    uint8
	clients []user
	lock    sync.Mutex
}
