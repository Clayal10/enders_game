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

	t.Run("TestBasicMessageSend", func(_ *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf(":%v", port))
		a.NoError(err)

		char := &lurk.Character{
			Type:       lurk.TypeCharacter,
			Name:       "Clay",
			Attack:     100,
			Defense:    90,
			Regen:      80,
			Health:     -70,
			Gold:       60,
			RoomNum:    2,
			PlayerDesc: "A guy who is just programming a game server",
		}

		ba, err := lurk.Marshal(char)
		a.NoError(err)
		conn.SetDeadline(time.Now().Add(time.Second))
		n, err := conn.Write(ba)
		a.NoError(err)

		buffer := make([]byte, n)
		n, err = conn.Read(buffer) //  currently expecting the same value to be sent back.
		a.NoError(err)

		msg, err := lurk.Unmarshal(buffer[:n])
		a.NoError(err)

		character, ok := msg.(*lurk.Character)
		a.True(ok)
		a.True(character.Type == char.Type)
		a.True(character.Name == char.Name)
		a.True(character.Attack == char.Attack)
		a.True(character.Defense == char.Defense)
		a.True(character.Health == char.Health)
		a.True(character.PlayerDesc == char.PlayerDesc)
	})
}
