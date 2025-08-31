package server

import (
	"fmt"
	"net"
	"testing"

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
}
