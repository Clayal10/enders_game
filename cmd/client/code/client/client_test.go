package client

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Clayal10/enders_game/pkg/assert"
	"github.com/Clayal10/enders_game/pkg/lurk"
)

func TestQueueReading(t *testing.T) {
	a := assert.New(t)

	c := newClient(nil, 0)

	c.q.Enqueue(&lurk.Fight{}, &lurk.Leave{}, &lurk.Loot{TargetName: "test"})

	messages := c.dequeueAll()
	a.True(len(messages) == 3)
	_, ok := messages[0].(*lurk.Fight)
	a.True(ok)
	_, ok = messages[1].(*lurk.Leave)
	a.True(ok)
	loot, ok := messages[2].(*lurk.Loot)
	a.True(ok)
	a.True(loot.TargetName == "test")
}

func TestDefaultCharacter(t *testing.T) {
	a := assert.New(t)

	t.Run("TestNilBody", func(_ *testing.T) {
		c := newClient(nil, 1)
		c.Game = &lurk.Game{
			InitialPoints: 66,
		}
		character, err := c.getOrMakeCharacter(nil)
		a.NoError(err)
		a.True(strings.Contains(character.Name, "Character 1"))
		a.True(character.Attack == 22)
		a.True(character.Defense == 22)
		a.True(character.Regen == 22)
	})
	t.Run("TestJavascriptResponse", func(_ *testing.T) {
		c := newClient(nil, 1)
		c.Game = &lurk.Game{
			InitialPoints: 66,
		}

		filename := filepath.Join(os.TempDir(), "test.txt")
		fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.ModePerm)
		a.NoError(err)
		defer fd.Close()

		n, err := fd.Write([]byte(`{
			"name": "Tester!",
			"attack": "nil"
		}`))
		a.NoError(err)
		a.True(n > 0)

		file, err := os.Open(filename)
		a.NoError(err)

		character, err := c.getOrMakeCharacter(file)
		a.NoError(err)
		a.True(strings.Contains(character.Name, "Tester!"))
		a.True(character.Attack == 22)
		a.True(character.Defense == 22)
		a.True(character.Regen == 22)
	})

	t.Run("TestNilGame", func(_ *testing.T) {
		c := newClient(nil, 1)
		_, err := c.getOrMakeCharacter(nil)
		a.Error(err)
	})

}
