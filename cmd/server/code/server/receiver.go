package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Clayal10/enders_game/pkg/cross"
)

type receiver struct {
	listener *net.TCPListener
	// queue for messages
	shouldRun bool
	*game
}

var terminationTimeout = 2 * time.Second

func newReceiver(cfg *Config, game *game) (*receiver, error) {
	address := fmt.Sprintf("0.0.0.0:%v", cfg.Port)

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
	defer rec.cleanup(player)
	if err != nil {
		log.Printf("%v: error registering player", err.Error())
		return
	}

	if err := rec.startGameplay(player, conn); err != nil && !errors.Is(err, errDisconnect) {
		log.Printf("%v: error during gameplay", err.Error())
		return
	}
	rec.cleanup(player)
	log.Printf("%v left.", player)
	time.Sleep(terminationTimeout)
}

func (rec *receiver) stop() {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	rec.shouldRun = false
}

func (rec *receiver) cleanup(player string) {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	delete(rec.users, player)
}
