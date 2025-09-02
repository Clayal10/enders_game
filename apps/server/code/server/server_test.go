package server

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/Clayal10/enders_game/lib/assert"
	"github.com/Clayal10/enders_game/lib/cross"
	"github.com/Clayal10/enders_game/lib/lurk"
)

// Global byte buffer for log checking.
var buf bytes.Buffer

func TestServerFunctionality(t *testing.T) {
	a := assert.New(t)

	log.SetOutput(&buf)

	port := cross.GetFreePort()
	cfg := &ServerConfig{
		Port: port,
	}

	cfs, err := New(cfg)
	a.NoError(err)
	defer func() {
		for _, cf := range cfs {
			cf()
		}
	}()

	t.Run("TestInvalidCharacterStats", func(_ *testing.T) {
		char := &lurk.Character{
			Type:       lurk.TypeCharacter,
			Name:       "Invalid guy",
			Attack:     100,
			Health:     100,
			Defense:    90,
			Regen:      80,
			RoomNum:    2,
			PlayerDesc: "A guy who is just programming a game server",
		}

		conn, err := net.Dial("tcp", fmt.Sprintf(":%v", cfg.Port))
		a.NoError(err)

		buffer, n, err := readSingleMessage(conn)
		a.NoError(err)

		msg, err := lurk.Unmarshal(buffer[:n])
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeVersion)

		buffer, n, err = readSingleMessage(conn)
		a.NoError(err)

		msg, err = lurk.Unmarshal(buffer[:n])
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeGame)

		ba, err := lurk.Marshal(char)
		a.NoError(err)

		_, err = conn.Write(ba) // send character
		a.NoError(err)

		buffer, _, err = readSingleMessage(conn) // read error
		a.NoError(err)

		msg, err = lurk.Unmarshal(buffer)
		a.NoError(err)

		e, ok := msg.(*lurk.Error)
		a.True(ok)

		a.True(strings.Contains(e.ErrMessage, "has invalid stats"))
	})
	t.Run("TestBasicGameplay", func(_ *testing.T) {
		conn := startClientConnection(a, cfg, &lurk.Character{
			Type:       lurk.TypeCharacter,
			Name:       "Tester",
			Attack:     100,
			Defense:    90,
			Regen:      80,
			RoomNum:    2,
			PlayerDesc: "A guy who is just programming a game server",
		})

		conn2 := startClientConnection(a, cfg, &lurk.Character{
			Type:       lurk.TypeCharacter,
			Name:       "Tester 2",
			Attack:     100,
			Defense:    90,
			Regen:      80,
			RoomNum:    2,
			PlayerDesc: "A guy who is just programming a game server",
		})

		buffer, _, err := readSingleMessage(conn) // read the 'room'
		a.NoError(err)

		msg, err := lurk.Unmarshal(buffer)
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeRoom)
		// Do more assertions and other gameplay here!
		// Need to read ALL messages

		ba, err := lurk.Marshal(&lurk.Leave{
			Type: lurk.TypeLeave,
		})
		a.NoError(err)
		_, err = conn.Write(ba)
		a.NoError(err)

		time.Sleep(50 * time.Millisecond)

		a.True(strings.Contains(buf.String(), "User left."))
		/* Termination of conn*/

		buffer, _, err = readSingleMessage(conn2)
		a.NoError(err)

		msg, err = lurk.Unmarshal(buffer)
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeRoom)

		// Send invalid stuff to server.
		ba, err = lurk.Marshal(&lurk.Accept{
			Type:   lurk.TypeAccept,
			Action: lurk.TypeAccept,
		})
		a.NoError(err)

		_, err = conn2.Write(ba)
		a.NoError(err)

		errMessage := readUntil(a, lurk.TypeError, conn2)

		a.True(errMessage != nil)
		a.True(errMessage.GetType() == lurk.TypeError)
		e, ok := errMessage.(*lurk.Error)
		a.True(ok)
		a.True(strings.Contains(e.ErrMessage, "Message contains invalid fields"))

		ba, err = lurk.Marshal(&lurk.Leave{
			Type: lurk.TypeLeave,
		})
		a.NoError(err)
		_, err = conn2.Write(ba)
		a.NoError(err)
	})

	t.Run("TestInitialMessageErrors", func(_ *testing.T) {
		char := &lurk.Character{
			Type:       lurk.TypeCharacter,
			Name:       "Tester",
			Attack:     100,
			Defense:    90,
			Regen:      80,
			RoomNum:    2,
			PlayerDesc: "A guy who is just programming a game server",
		}

		conn, err := net.Dial("tcp", fmt.Sprintf(":%v", cfg.Port))
		a.NoError(err)

		buffer, n, err := readSingleMessage(conn)
		a.NoError(err)

		msg, err := lurk.Unmarshal(buffer[:n])
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeVersion)

		buffer, n, err = readSingleMessage(conn)
		a.NoError(err)

		msg, err = lurk.Unmarshal(buffer[:n])
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeGame)

		ba, err := lurk.Marshal(char)
		a.NoError(err)

		ba[0] = 20

		_, err = conn.Write(ba) // send character with bad type.
		a.NoError(err)

		errMessage := readUntil(a, lurk.TypeError, conn)
		a.True(errMessage != nil)
		a.True(errMessage.GetType() == lurk.TypeError)
		e, ok := errMessage.(*lurk.Error)
		a.True(ok)
		a.True(strings.Contains(e.ErrMessage, "Bad message"))

	})
}

func TestServerStartupErrors(t *testing.T) {
	a := assert.New(t)

	log.SetOutput(&buf)
	t.Run("TestBadIPandPort", func(_ *testing.T) {
		port := cross.GetFreePort()
		cfg := &ServerConfig{
			Port: port,
		}

		cfs, err := New(cfg)
		a.NoError(err)

		_, err = New(cfg)
		a.Error(err)
		a.True(strings.Contains(buf.String(), "Could not listen on port"))

		for _, cf := range cfs {
			cf()
		}
	})
}

func readUntil(a *assert.Assert, t lurk.MessageType, conn net.Conn) lurk.LurkMessage {
	_ = conn.SetDeadline(time.Now().Add(500 * time.Millisecond))
	for {
		buffer, _, err := readSingleMessage(conn)
		if err != nil {
			if errors.Is(err, cross.ErrInvalidMessageType) {
				continue
			}
			return nil
		}

		msg, err := lurk.Unmarshal(buffer)
		a.NoError(err)

		mt := msg.GetType()
		if mt == t {
			return msg
		}
	}
}

func startClientConnection(a *assert.Assert, cfg *ServerConfig, char *lurk.Character) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%v", cfg.Port))
	a.NoError(err)

	buffer, n, err := readSingleMessage(conn)
	a.NoError(err)

	msg, err := lurk.Unmarshal(buffer[:n])
	a.NoError(err)

	a.True(msg.GetType() == lurk.TypeVersion)

	buffer, n, err = readSingleMessage(conn)
	a.NoError(err)

	msg, err = lurk.Unmarshal(buffer[:n])
	a.NoError(err)

	a.True(msg.GetType() == lurk.TypeGame)

	ba, err := lurk.Marshal(char)
	a.NoError(err)

	_, err = conn.Write(ba) // send character
	a.NoError(err)

	buffer, n, err = readSingleMessage(conn) // read accept
	a.NoError(err)

	msg, err = lurk.Unmarshal(buffer[:n])
	a.NoError(err)

	a.True(msg.GetType() == lurk.TypeAccept)
	buffer, n, err = readSingleMessage(conn)
	a.NoError(err)

	msg, err = lurk.Unmarshal(buffer[:n])
	a.NoError(err)

	a.True(msg.GetType() == lurk.TypeCharacter)

	ba, err = lurk.Marshal(&lurk.Start{
		Type: lurk.TypeStart,
	})
	a.NoError(err)

	_, err = conn.Write(ba)
	a.NoError(err)

	buffer, _, err = readSingleMessage(conn)
	a.NoError(err)

	a.True(buffer[0] == byte(lurk.TypeAccept))

	return conn
}
