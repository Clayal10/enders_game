package server

import (
	"net"

	"github.com/Clayal10/enders_game/lib/cross"
	"github.com/Clayal10/enders_game/lib/lurk"
)

// These functions may need thread protection before getting called.

// default game values.
const (
	initialPoints = 100
	statLimit     = 65535
	narrator      = "Narrator"
)

func (g *game) sendStart(conn net.Conn) error {
	version := &lurk.Version{
		Type:  lurk.TypeVersion,
		Major: 2,
		Minor: 3,
	}

	g.game = &lurk.Game{
		Type:          lurk.TypeGame,
		InitialPoints: initialPoints,
		StatLimit:     statLimit,
		GameDesc:      gameDescription,
	}

	if _, err := conn.Write(lurk.Marshal(version)); err != nil {
		return err
	}

	_, err := conn.Write(lurk.Marshal(g.game))
	return err
}

func (g *game) sendRoom(room *room, player string, conn net.Conn) error {
	if _, err := conn.Write(lurk.Marshal(room.r)); err != nil {
		return err
	}

	if err := g.sendCharacters(room, player, conn); err != nil {
		return err
	}

	return g.sendConnections(room, player, conn)
}

func (g *game) sendCharacters(room *room, player string, conn net.Conn) (err error) {
	// all characters and monsters in that room
	for _, user := range g.users {
		// should we include current user?
		if user.c.RoomNum != room.r.RoomNumber {
			continue
		}
		if _, err = conn.Write(lurk.Marshal(user.c)); err != nil {
			return
		}
	}

	for _, npc := range g.monsters {
		if npc.RoomNum != room.r.RoomNumber {
			continue
		}
		if _, err = conn.Write(lurk.Marshal(npc)); err != nil {
			return
		}
	}
	return
}

// Takes a user object and sends it to conn. Used for notifying other users of a user's status.
// The message will be sent if the the recipient isn't allowed to know what room the user is
// moving to.
func (g *game) sendCharacterUpdate(user *lurk.Character, conn net.Conn, recipient string, message string) error {
	ba := lurk.Marshal(user)
	if _, err := conn.Write(ba); err != nil {
		return err
	}

	if message == "" {
		return nil
	}

	_, err := conn.Write(lurk.Marshal(&lurk.Message{
		Type:      lurk.TypeMessage,
		Recipient: recipient,
		Sender:    narrator,
		Text:      message,
		Narration: true,
	}))
	return err
}

func (g *game) sendConnections(room *room, player string, conn net.Conn) (err error) {
	for _, connection := range room.connections {
		if !g.users[player].allowedRoom[connection.RoomNumber] {
			continue
		}
		if _, err = conn.Write(lurk.Marshal(connection)); err != nil {
			return err
		}
	}
	return
}

func (g *game) sendAccept(conn net.Conn, action lurk.MessageType) error {
	accept := &lurk.Accept{
		Type:   lurk.TypeAccept,
		Action: action,
	}
	_, err := conn.Write(lurk.Marshal(accept))
	return err
}

func (g *game) sendError(conn net.Conn, code cross.ErrCode, msg string) error {
	_, err := conn.Write(lurk.Marshal(&lurk.Error{
		Type:       lurk.TypeError,
		ErrCode:    code,
		ErrMessage: msg,
	}))
	return err
}
