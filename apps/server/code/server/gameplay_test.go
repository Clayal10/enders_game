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
					a.NoError(err)
					err = conn.SetReadDeadline(time.Now().Add(1000000 * time.Millisecond))
					a.NoError(err)

					ba, _, err := readSingleMessage(conn)
					a.NoError(err)
					msg, err := lurk.Unmarshal(ba)
					a.NoError(err)
					a.True(msg.GetType() == lurk.TypeError)
					e := msg.(*lurk.Error)
					a.True(strings.Contains(e.ErrMessage, cross.ErrRoomsNotConnected.Error()))
				}
			}
		}(ctx)

		c, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", port))
		a.NoError(err)

		g := newGame()

		// No player
		a.Error(g.handleChangeRoom(&lurk.ChangeRoom{
			Type:       lurk.TypeChangeRoom,
			RoomNumber: 100,
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
	t.Run("TestPVPFight", func(_ *testing.T) {
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

		conn1 := startClientConnection(a, cfg, &lurk.Character{
			Name: "t1",
			Flags: map[string]bool{
				lurk.Alive: true,
			},
			Attack:     100,
			Defense:    0,
			Regen:      0,
			PlayerDesc: "First Tester",
		})

		a.Eventually(func() bool {
			_ = conn1.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
			ba, _, err := readSingleMessage(conn1)
			a.NoError(err)
			lmsg, err := lurk.Unmarshal(ba)
			a.NoError(err)
			if lmsg.GetType() != lurk.TypeCharacter {
				return false
			}
			character, ok := lmsg.(*lurk.Character)
			a.True(ok)
			return character.Name == "t1"
		}, time.Second, 20*time.Millisecond)

		conn2 := startClientConnection(a, cfg, &lurk.Character{
			Name: "t2",
			Flags: map[string]bool{
				lurk.Alive: true,
			},
			Attack:     50,
			Defense:    0,
			Regen:      0,
			PlayerDesc: "Second Tester",
		})

		a.Eventually(func() bool {
			_ = conn1.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
			ba, _, err := readSingleMessage(conn1)
			a.NoError(err)
			lmsg, err := lurk.Unmarshal(ba)
			a.NoError(err)
			if lmsg.GetType() != lurk.TypeCharacter {
				return false
			}
			character, ok := lmsg.(*lurk.Character)
			a.True(ok)
			return character.Name == "t2"
		}, time.Second, 20*time.Millisecond)

		_, err = conn1.Write(lurk.Marshal(&lurk.PVPFight{
			TargetName: "t2",
		}))
		a.NoError(err)
		a.Eventually(func() bool {
			_ = conn1.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
			ba, _, err := readSingleMessage(conn1)
			a.NoError(err)
			lmsg, err := lurk.Unmarshal(ba)
			a.NoError(err)
			if lmsg.GetType() != lurk.TypeCharacter {
				return false
			}
			character, ok := lmsg.(*lurk.Character)
			a.True(ok)
			return character.Name == "t2" && !character.Flags[lurk.Alive]
		}, time.Second*100, 20*time.Millisecond)

		_, err = conn2.Write(lurk.Marshal(&lurk.Leave{}))
		a.NoError(err)

		a.Eventually(func() bool {
			_ = conn1.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
			ba, _, err := readSingleMessage(conn1)
			a.NoError(err)
			lmsg, err := lurk.Unmarshal(ba)
			a.NoError(err)
			if lmsg.GetType() != lurk.TypeMessage {
				return false
			}
			msg, ok := lmsg.(*lurk.Message)
			a.True(ok)
			return msg.Narration && strings.Contains(msg.Text, "t2 left the server!")
		}, time.Second*100, 20*time.Millisecond)
		_, err = conn1.Write(lurk.Marshal(&lurk.Leave{}))
		a.NoError(err)

	})
	t.Run("TestKillingHiveQueenCocoon", func(_ *testing.T) {
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

		conn := startClientConnection(a, cfg, &lurk.Character{
			Name: "Beans Shumaker",
			Flags: map[string]bool{
				lurk.Alive: true,
			},
			Attack:     100,
			Defense:    0,
			Regen:      0,
			PlayerDesc: "Admin",
		})

		_, err = conn.Write(lurk.Marshal(&lurk.ChangeRoom{RoomNumber: eros}))
		a.NoError(err)
		_, err = conn.Write(lurk.Marshal(&lurk.ChangeRoom{RoomNumber: shakespeare}))
		a.NoError(err)

		a.Eventually(func() bool {
			_ = conn.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
			ba, _, err := readSingleMessage(conn)
			a.NoError(err)
			lmsg, err := lurk.Unmarshal(ba)
			a.NoError(err)
			if lmsg.GetType() != lurk.TypeCharacter {
				return false
			}
			character, ok := lmsg.(*lurk.Character)
			a.True(ok)
			return character.Name == hiveQueenCocoon && !character.Flags[lurk.Monster]
		}, time.Second*100, 20*time.Millisecond)

		_, err = conn.Write(lurk.Marshal(&lurk.PVPFight{
			TargetName: hiveQueenCocoon,
		}))
		a.NoError(err)

		a.Eventually(func() bool {
			_ = conn.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
			ba, _, err := readSingleMessage(conn)
			a.NoError(err)
			lmsg, err := lurk.Unmarshal(ba)
			a.NoError(err)
			if lmsg.GetType() != lurk.TypeMessage {
				return false
			}
			msg, ok := lmsg.(*lurk.Message)
			a.True(ok)
			return strings.Contains(msg.Text, "Xenocide")
		}, time.Second*100, 20*time.Millisecond)
	})
	t.Run("TestUpgradingStats", func(_ *testing.T) {
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

		conn := startClientConnection(a, cfg, &lurk.Character{
			Name: "Test Guy",
			Flags: map[string]bool{
				lurk.Alive: true,
			},
			Attack:     50,
			Defense:    25,
			Regen:      25,
			PlayerDesc: "Test guy who will upgrade stats",
		})
		_, err = conn.Write(lurk.Marshal(&lurk.ChangeRoom{
			RoomNumber: battleSchoolGameRoom,
		}))
		a.NoError(err)
		_, err = conn.Write(lurk.Marshal(&lurk.Fight{}))
		a.NoError(err)
		_, err = conn.Write(lurk.Marshal(&lurk.Fight{}))
		a.NoError(err)
		_, err = conn.Write(lurk.Marshal(&lurk.Fight{}))
		a.NoError(err)
		_, err = conn.Write(lurk.Marshal(&lurk.ChangeRoom{
			RoomNumber: battleSchool,
		}))
		a.NoError(err)
		_, err = conn.Write(lurk.Marshal(&lurk.ChangeRoom{
			RoomNumber: battleSchoolBarracks,
		}))
		a.NoError(err)
		a.Eventually(func() bool {
			_ = conn.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
			ba, _, err := readSingleMessage(conn)
			a.NoError(err)
			msg, err := lurk.Unmarshal(ba)
			a.NoError(err)
			fmt.Println(fmt.Sprint(msg.GetType()))
			if msg.GetType() != lurk.TypeMessage {
				return false
			}
			message := msg.(*lurk.Message)
			return strings.Contains(message.Text, "upgrade your stats")
		}, time.Second, time.Millisecond)
		_ = conn.SetReadDeadline(time.Time{})

	})
}
