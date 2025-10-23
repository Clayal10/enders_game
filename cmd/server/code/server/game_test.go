package server

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Clayal10/enders_game/pkg/assert"
	"github.com/Clayal10/enders_game/pkg/cross"
	"github.com/Clayal10/enders_game/pkg/lurk"
)

const bufferLength = 128

func TestReadAll(t *testing.T) {
	a := assert.New(t)
	t.Run("TestExtendedMessage", func(_ *testing.T) {
		port := cross.GetFreePort()

		l, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
		a.NoError(err)

		c, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", port))
		a.NoError(err)

		ctx, cf := context.WithCancel(context.Background())
		defer cf()
		go func(ctx context.Context, c net.Conn) {
			t := time.NewTicker(time.Millisecond * 50)
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					ba := lurk.Marshal(&lurk.Character{ // should overflow the 128 default buffer
						Type:       lurk.TypeCharacter,
						Name:       "Verryyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy long name",
						PlayerDesc: gameDescription + gameDescription + gameDescription,
					})
					n, err := c.Write(ba)
					a.NoError(err)
					a.True(n > bufferLength)
				}
			}
		}(ctx, c)

		conn, err := l.Accept()
		a.NoError(err)

		buffer, n, err := lurk.ReadSingleMessage(conn)
		a.NoError(err)
		a.True(n > bufferLength)

		msg, err := lurk.Unmarshal(buffer)
		a.NoError(err)

		a.True(msg.GetType() == lurk.TypeCharacter)
		character, ok := msg.(*lurk.Character)
		a.True(ok)

		a.True(len(character.Name) == 32)
		a.True(character.PlayerDesc == gameDescription+gameDescription+gameDescription)

	})
}
