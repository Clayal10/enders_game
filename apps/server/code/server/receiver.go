package server

import (
	"net"
	"sync"
)

type receiver struct {
	listener net.TCPListener
	mu       sync.Mutex
	// queue for messages
}

func newReceiver(cfg *ServerConfig) (*receiver, error) {
	l, err := net.ListenTCP()

	return &receiver{}
}
