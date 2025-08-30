package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/Clayal10/enders_game/lib/lurk"
)

type receiver struct {
	listener *net.TCPListener
	mu       sync.Mutex
	// queue for messages
	shouldRun bool
}

const (
	bufferLength = 128 // Should cover most messages.
	minBuffer    = 48  // type Character's variable length field is at offset 46.
)

const readTimeout = 30 * time.Second

func newReceiver(cfg *ServerConfig) (*receiver, error) {
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
		go rec.handleConnection(conn)
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
func (rec *receiver) handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, bufferLength)

	conn.SetReadDeadline(time.Now().Add(readTimeout))

	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("%v: error reading from the connection", err.Error())
		return
	}

	if n >= bufferLength {
		n, err = lurk.GetVariableLength(buffer)
		if err != nil {
			log.Printf("%v: error processing message", err.Error())
			return
		}
		b := make([]byte, n-bufferLength)
		conn.Read(b)
		buffer = append(buffer, b...)
	}

	msg, err := lurk.Unmarshal(buffer[:n])
	if err != nil {
		log.Printf("%v: error processing full message", err.Error())
	}
	// TODO make a queue for LurkMessages to be taken care of by another thread.
	// For now, just write it back. This is temporary

	ba, err := lurk.Marshal(msg)
	if err != nil {
		log.Printf("%v: error turning message back into bytes", err.Error())
	}

	conn.Write(ba)
}

func (rec *receiver) stop() {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	rec.shouldRun = false
}
