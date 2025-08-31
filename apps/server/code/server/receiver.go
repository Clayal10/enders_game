package server

import (
	"fmt"
	"log"
	"net"
	"sync"
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
	defer rec.listener.Close()
	for {
		if !rec.shouldRun {
			break
		}
		conn, err := rec.listener.Accept()
		if err != nil {
			log.Printf("%v: error accepting connection", err.Error())
			continue
		}
		go rec.registerUser(conn)
	}
}

// This function needs to be able to handle reading and sending back information to the client
// at any moment.
//
// Thoughts:
//   - We can try to keep the connection as persistent as possible, send a warning when
//     we disconnect etc. and save the character data on disk for better persistence.
//   - This function can go in some registry that will 'register' the user, handle communication
//     and write their information to something at least persistent to this server execution.
func (rec *receiver) registerUser(conn net.Conn) {
	defer conn.Close()

	if err := rec.game.sendStart(conn); err != nil {
		log.Printf("%v: error starting the game", err.Error())
	}

	if err := rec.game.registerPlayer(conn); err != nil {
		log.Printf("%v: error registering player", err.Error())

	}

}

func (rec *receiver) stop() {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	rec.shouldRun = false
}
