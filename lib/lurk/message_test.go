package lurk_test

import (
	"testing"

	"github.com/Clayal10/enders_game/lib/assert"
	"github.com/Clayal10/enders_game/lib/lurk"
)

func TestUnmarshal(t *testing.T) {
	a := assert.New(t)
	t.Run("TestMessageType", func(_ *testing.T) {
		m, err := lurk.Unmarshal(sampleMessage)
		msg, ok := m.(*lurk.Message)
		a.True(ok)
		a.NoError(err)
		a.True(msg.Type == lurk.TypeMessage)
		a.True(msg.Narration)
		a.True(msg.RName == "Raymond")
		a.True(msg.SName == "Clay")

		ba, err := lurk.Marshal(msg)
		a.NoError(err)
		a.EqualSlice(ba, sampleMessage)
	})
	t.Run("TestChangeRoomType", func(_ *testing.T) {
		changeRoom := &lurk.ChangeRoom{
			Type:       lurk.TypeChangeRoom,
			RoomNumber: 2,
		}
		ba, err := lurk.Marshal(changeRoom)
		a.NoError(err)
		a.EqualSlice(ba, []byte{0x2, 0x2, 0x0})
		cr2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		changeRoom2, ok := cr2.(*lurk.ChangeRoom)
		a.True(ok)
		a.True(changeRoom.RoomNumber == changeRoom2.RoomNumber)
	})
}

// With narration
var sampleMessage = []byte{0x1, 0x5, 0x00,
	0x52, 0x61, 0x79, 0x6d, 0x6f, 0x6e, 0x64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0x43, 0x6c, 0x61, 0x79, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01,
	0x48, 0x65, 0x6c, 0x6c, 0x6f}
