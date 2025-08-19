package lurk_test

import (
	"testing"

	"github.com/Clayal10/enders_game/lib/assert"
	"github.com/Clayal10/enders_game/lib/lurk"
)

func TestUnmarshal(t *testing.T) {
	a := assert.New(t)
	t.Run("TestMessageUnmarshal", func(_ *testing.T) {
		m, err := lurk.Unmarshal(sampleMessage)
		msg, ok := m.(*lurk.Message)
		a.True(ok)
		a.NoError(err)
		a.True(msg.Type == lurk.TypeMessage)
		a.True(msg.Narration)
		a.True(msg.RName == "Raymond")
		a.True(msg.SName == "Clay")
	})
}

// With narration
var sampleMessage = []byte{0x1, 0x5, 0x00, 0x52, 0x61, 0x79, 0x6d, 0x6f, 0x6e, 0x64, 0x0, 0x43, 0x6c, 0x61, 0x79, 0x0, 0x1, 0x48, 0x65, 0x6c, 0x6c, 0x6f}
