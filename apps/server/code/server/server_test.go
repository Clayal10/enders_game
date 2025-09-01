package server

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Clayal10/enders_game/lib/assert"
	"github.com/Clayal10/enders_game/lib/cross"
	"github.com/Clayal10/enders_game/lib/lurk"
)

func TestServerFunctionality(t *testing.T) {
	a := assert.New(t)
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

	// This one doesn't leave.
	t.Run("TestGameStart", func(_ *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf(":%v", port))
		a.NoError(err)

		buffer := make([]byte, 1024)

		n, err := conn.Read(buffer)
		a.NoError(err)

		msg, err := lurk.Unmarshal(buffer[:n])
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeGame)

		char := &lurk.Character{
			Type:       lurk.TypeCharacter,
			Name:       "Clay",
			Attack:     100,
			Defense:    90,
			Regen:      80,
			RoomNum:    2,
			PlayerDesc: "A guy who is just programming a game server",
		}

		ba, err := lurk.Marshal(char)
		a.NoError(err)

		_, err = conn.Write(ba) // send character
		a.NoError(err)

		acceptBuffer := buffer[:2]       // keep reusing the buffer
		n, err = conn.Read(acceptBuffer) // read accept
		a.NoError(err)

		msg, err = lurk.Unmarshal(buffer[:n])
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeAccept)
		n, err = conn.Read(buffer)
		a.NoError(err)

		msg, err = lurk.Unmarshal(buffer[:n])
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeCharacter)
	})
	t.Run("TestInvalidCharacterStats", func(_ *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf(":%v", port))
		a.NoError(err)

		buffer := make([]byte, 1024)

		n, err := conn.Read(buffer)
		a.NoError(err)

		msg, err := lurk.Unmarshal(buffer[:n])
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeGame)

		char := &lurk.Character{
			Type:       lurk.TypeCharacter,
			Name:       "Clay",
			Attack:     100,
			Health:     100,
			Defense:    90,
			Regen:      80,
			RoomNum:    2,
			PlayerDesc: "A guy who is just programming a game server with stats that I'm not supposed to put in the character",
		}

		ba, err := lurk.Marshal(char)
		a.NoError(err)

		_, err = conn.Write(ba)
		a.NoError(err)

		n, err = conn.Read(buffer)
		a.NoError(err)

		msg, err = lurk.Unmarshal(buffer[:n])
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeError)
		e, ok := msg.(*lurk.Error)
		a.True(ok)
		a.True(e.ErrCode == cross.StatError)
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

		buffer, _, err := readAll(conn) // read the 'room'
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

		//buffer, _, err = readAll(conn)
		//a.NoError(err)
		//a.True(buffer[0] == byte(lurk.TypeAccept))
	})
}

func startClientConnection(a *assert.Assert, cfg *ServerConfig, char *lurk.Character) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf(":%v", cfg.Port))
	a.NoError(err)

	buffer, n, err := readAll(conn)
	a.NoError(err)

	msg, err := lurk.Unmarshal(buffer[:n])
	a.NoError(err)

	a.True(msg.GetType() == lurk.TypeGame)

	ba, err := lurk.Marshal(char)
	a.NoError(err)

	_, err = conn.Write(ba) // send character
	a.NoError(err)

	buffer, n, err = readAll(conn) // read accept
	a.NoError(err)

	msg, err = lurk.Unmarshal(buffer[:n])
	a.NoError(err)

	a.True(msg.GetType() == lurk.TypeAccept)
	buffer, n, err = readAll(conn)
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

	buffer, _, err = readAll(conn)
	a.NoError(err)

	a.True(buffer[0] == byte(lurk.TypeAccept))

	return conn
}
