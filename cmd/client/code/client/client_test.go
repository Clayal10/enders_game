package client

import (
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

	c := newClient(nil, 1)
	character, err := c.getOrMakeCharacter(nil)
	a.NoError(err)
	a.True(strings.Contains(character.Name, "Character 1"))
}
