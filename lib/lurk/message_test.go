package lurk_test

import (
	"encoding/binary"
	"testing"

	"github.com/Clayal10/enders_game/lib/assert"
	"github.com/Clayal10/enders_game/lib/cross"
	"github.com/Clayal10/enders_game/lib/lurk"
)

func TestUnmarshalAndMarshal(t *testing.T) {
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
	t.Run("TestFightType", func(_ *testing.T) {
		f := &lurk.Fight{
			Type: lurk.TypeFight,
		}
		ba, err := lurk.Marshal(f)
		a.NoError(err)
		a.EqualSlice(ba, []byte{0x3})
		f2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		fight2, ok := f2.(*lurk.Fight)
		a.True(ok)
		a.True(fight2.Type == f.Type)
	})
	t.Run("TestPVPType", func(_ *testing.T) {
		pvp := &lurk.PVPFight{
			Type:       lurk.TypePVPFight,
			TargetName: "Clay",
		}
		ba, err := lurk.Marshal(pvp)
		a.NoError(err)
		a.EqualSlice(ba, []byte{0x04, 0x43, 0x6c, 0x61, 0x79, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		pvp2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		pvpType, ok := pvp2.(*lurk.PVPFight)
		a.True(ok)
		a.True(pvp.TargetName == pvpType.TargetName)
	})
	t.Run("TestLootType", func(_ *testing.T) {
		loot := &lurk.Loot{
			Type:       lurk.TypeLoot,
			TargetName: "Clay",
		}
		ba, err := lurk.Marshal(loot)
		a.NoError(err)
		a.EqualSlice(ba, []byte{0x05, 0x43, 0x6c, 0x61, 0x79, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		l2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		loot2, ok := l2.(*lurk.Loot)
		a.True(ok)
		a.True(loot.TargetName == loot2.TargetName)
	})
	t.Run("TestStartType", func(_ *testing.T) {
		s := &lurk.Start{
			Type: lurk.TypeStart,
		}
		ba, err := lurk.Marshal(s)
		a.NoError(err)
		a.EqualSlice(ba, []byte{0x6})
		s2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		start2, ok := s2.(*lurk.Start)
		a.True(ok)
		a.True(start2.Type == s.Type)
	})
	t.Run("TestErrorType", func(_ *testing.T) {
		e := &lurk.Error{
			Type:       lurk.TypeError,
			ErrCode:    cross.Other,
			ErrMessage: "Hello",
		}
		ba, err := lurk.Marshal(e)
		a.NoError(err)
		a.EqualSlice(ba, []byte{0x7, 0x0, 0x5, 0x0, 0x48, 0x65, 0x6c, 0x6c, 0x6f})
		e2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		error2, ok := e2.(*lurk.Error)
		a.True(ok)
		a.True(e.ErrMessage == error2.ErrMessage)
		a.True(e.ErrCode == error2.ErrCode)
	})
	t.Run("TestAcceptType", func(_ *testing.T) {
		accept := &lurk.Accept{
			Type:   lurk.TypeAccept,
			Action: lurk.TypeMessage,
		}
		ba, err := lurk.Marshal(accept)
		a.NoError(err)
		a.EqualSlice(ba, []byte{0x08, 0x1})
		a2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		accept2, ok := a2.(*lurk.Accept)
		a.True(ok)
		a.True(accept.Action == accept2.Action)
	})
	t.Run("TestRoomType", func(_ *testing.T) {
		room := &lurk.Room{
			Type:       lurk.TypeRoom,
			RoomNumber: 1,
			RoomName:   "Test",
			RoomDesc:   "Test Room",
		}
		ba, err := lurk.Marshal(room)
		a.NoError(err)
		r2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		room2, ok := r2.(*lurk.Room)
		a.True(ok)
		a.True(room2.RoomName == "Test")
		a.True(room2.RoomDesc == "Test Room")
	})
	t.Run("TestCharacterType", func(_ *testing.T) {
		flags := make(map[string]bool)
		flags[lurk.Alive] = true
		flags[lurk.JoinBattle] = false
		flags[lurk.Monster] = false
		flags[lurk.Ready] = false
		flags[lurk.Started] = false
		character := &lurk.Character{
			Type:       lurk.TypeCharacter,
			Name:       "Clay",
			Flags:      flags,
			Attack:     1,
			Defense:    2,
			Regen:      3,
			Health:     -4,
			Gold:       5,
			RoomNum:    6,
			PlayerDesc: "This is Clay",
		}
		ba, err := lurk.Marshal(character)
		a.NoError(err)
		char2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		character2, ok := char2.(*lurk.Character)
		a.True(ok)
		a.True(character2.Flags[lurk.Alive])
		a.True(!character2.Flags[lurk.JoinBattle])
		a.True(character2.Name == character.Name)
		a.True(character2.Health == character.Health)
		a.True(character2.PlayerDesc == character.PlayerDesc)
	})
	t.Run("TestGameType", func(_ *testing.T) {
		game := &lurk.Game{
			Type:          lurk.TypeGame,
			InitialPoints: 50,
			StatLimit:     100,
			GameDesc:      "This is a game",
		}
		ba, err := lurk.Marshal(game)
		a.NoError(err)
		g2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		game2, ok := g2.(*lurk.Game)
		a.True(ok)
		a.True(*game2 == *game)
	})
	t.Run("TestConnectionType", func(_ *testing.T) {
		c := &lurk.Connection{
			Type:       lurk.TypeConnection,
			RoomNumber: 1,
			RoomName:   "Test Room",
			RoomDesc:   "Just a test room",
		}
		ba, err := lurk.Marshal(c)
		a.NoError(err)
		c2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		connection2, ok := c2.(*lurk.Connection)
		a.True(ok)
		a.True(*connection2 == *c)
	})
	t.Run("TestExtensionsType", func(_ *testing.T) {
		e := &lurk.Version{
			Type:  lurk.TypeVersion,
			Major: 2,
			Minor: 3,
			Extensions: [][]byte{
				{0xFF, 0xFF},
				{0xFF, 0xFF},
			},
		}
		ba, err := lurk.Marshal(e)
		a.NoError(err)
		a.True(binary.LittleEndian.Uint16(ba[3:]) == uint16(8))
		a.True(binary.LittleEndian.Uint16(ba[5:]) == uint16(2))
		e2, err := lurk.Unmarshal(ba)
		a.NoError(err)
		extension2, ok := e2.(*lurk.Version)
		a.True(ok)
		a.True(extension2.Major == e.Major)
		a.True(extension2.Minor == e.Minor)
		a.EqualSlice(extension2.Extensions[0], e.Extensions[0])
		a.EqualSlice(extension2.Extensions[1], e.Extensions[1])
	})
}

// With narration
var sampleMessage = []byte{0x1, 0x5, 0x00,
	0x52, 0x61, 0x79, 0x6d, 0x6f, 0x6e, 0x64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0x43, 0x6c, 0x61, 0x79, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x01,
	0x48, 0x65, 0x6c, 0x6c, 0x6f}
