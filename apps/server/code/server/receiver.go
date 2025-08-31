package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/Clayal10/enders_game/lib/cross"
)

type receiver struct {
	listener *net.TCPListener
	mu       sync.Mutex
	// queue for messages
	shouldRun bool
	*game
}

const (
	bufferLength = 128 // Should cover most messages.
	minBuffer    = 48  // type Character's variable length field is at offset 46.
)

func newReceiver(cfg *ServerConfig, game *game) (*receiver, error) {
	address := fmt.Sprintf("localhost:%v", cfg.Port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Printf("Could not resolve address '%v'", address)
		return nil, err
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("Could not listen on port %v", cfg.Port)
		return nil, err
	}

	return &receiver{
		listener:  l,
		shouldRun: true,
		game:      game,
	}, nil
}

func (rec *receiver) start() {
	go rec.run()
}

func (rec *receiver) run() {
	defer cross.LogOnErr(rec.listener.Close)
	for rec.shouldRun {
		conn, err := rec.listener.Accept()
		if err != nil {
			log.Printf("%v: error accepting connection", err.Error())
			continue
		}
		go rec.registerUser(conn)
	}
}

// The 'conn' object will simply get passed through to different functions.
func (rec *receiver) registerUser(conn net.Conn) {
	defer cross.LogOnErr(conn.Close)

	if err := rec.sendStart(conn); err != nil {
		log.Printf("%v: error starting the game", err.Error())
	}

	if err := rec.registerPlayer(conn); err != nil {
		log.Printf("%v: error registering player", err.Error())
	}

}

func (rec *receiver) stop() {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	rec.shouldRun = false
}
