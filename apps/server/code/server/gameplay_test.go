package server

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/Clayal10/enders_game/lib/assert"
	"github.com/Clayal10/enders_game/lib/cross"
	"github.com/Clayal10/enders_game/lib/lurk"
)

func TestGameActions(t *testing.T) {
	a := assert.New(t)
	t.Run("TestSendBadRoom", func(_ *testing.T) {
		port := cross.GetFreePort()
		l, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
		a.NoError(err)

		ctx, cf := context.WithCancel(context.Background())
		defer cf()
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					conn, err := l.Accept()
					_ = conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
					a.NoError(err)
					ba, _, err := readSingleMessage(conn)
					a.NoError(err)
					msg, err := lurk.Unmarshal(ba)
					a.NoError(err)
					a.True(msg.GetType() == lurk.TypeError)
					e := msg.(*lurk.Error)
					a.True(strings.Contains(e.ErrMessage, cross.ErrUserNotInServer.Error()))

					ba, _, err = readSingleMessage(conn)
					a.NoError(err)
					msg, err = lurk.Unmarshal(ba)
					a.NoError(err)
					a.True(msg.GetType() == lurk.TypeError)
					e = msg.(*lurk.Error)
					a.True(strings.Contains(e.ErrMessage, cross.ErrRoomsNotConnected.Error()))
				}
			}
		}(ctx)

		c, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", port))
		a.NoError(err)

		g := newGame()

		// No player
		a.NoError(g.handleChangeRoom(&lurk.ChangeRoom{
			Type:       lurk.TypeChangeRoom,
			RoomNumber: 100, // doesn't exist.
		}, c, "Test"))

		testName := "test name"
		g.users[testName] = &user{
			conn: c,
			c: &lurk.Character{
				Type:    lurk.TypeCharacter,
				Name:    testName,
				RoomNum: 1,
			},
		}
		// bad room number
		a.NoError(g.handleChangeRoom(&lurk.ChangeRoom{
			Type:       lurk.TypeChangeRoom,
			RoomNumber: 100, // doesn't exist.
		}, c, testName))
	})
}
