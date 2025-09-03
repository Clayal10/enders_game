package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Clayal10/enders_game/lib/cross"
	"github.com/Clayal10/enders_game/lib/lurk"
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
)

var terminationTimeout = 2 * time.Second

func newReceiver(cfg *ServerConfig, game *game) (*receiver, error) {
	address := fmt.Sprintf("localhost:%v", cfg.Port)

	// Won't fail with the preset localhost and "tcp".
	tcpAddr, _ := net.ResolveTCPAddr("tcp", address)

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
		return
	}

	player, err := rec.registerPlayer(conn)
	defer delete(rec.users, player)
	if err != nil {
		log.Printf("%v: error registering player", err.Error())
		return
	}

	if err := rec.startGameplay(player, conn); err != nil {
		// termination of gameplay
		log.Printf("%v: error during gameplay", err.Error())
		return
	}
	log.Printf("User left.")
	delete(rec.users, player)
	time.Sleep(terminationTimeout)
}

func (rec *receiver) stop() {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	rec.shouldRun = false
}

// We want to read exactly the length of the message. This function will do up to 3
// calls to 'Read' to read exactly one message.
func readSingleMessage(conn net.Conn) ([]byte, int, error) {
	// TODO find a way to flush out after an error
	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	if err != nil {
		return nil, 0, err
	}

	bytesNeeded, ok := lurk.LengthOffset[lurk.MessageType(buffer[0])]
	if !ok {
		return nil, 0, cross.ErrInvalidMessageType
	}

	if bytesNeeded == 1 {
		return buffer, 1, nil
	}

	b := make([]byte, bytesNeeded-1)
	_, err = conn.Read(b)
	if err != nil {
		return nil, 0, err
	}
	buffer = append(buffer, b...)

	n, err := lurk.GetVariableLength(buffer)
	if err != nil {
		return nil, 0, err
	}

	if n == -1 {
		return buffer, bytesNeeded, nil
	}

	b = make([]byte, n)
	if _, err = conn.Read(b); err != nil {
		return nil, 0, err
	}
	buffer = append(buffer, b...)
	return buffer, n + bytesNeeded, nil
}
