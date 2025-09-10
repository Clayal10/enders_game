package server

import (
	"net"

	"github.com/Clayal10/enders_game/lib/lurk"
)

// These functions may need thread protection before getting called.

func (g *game) sendStart(conn net.Conn) error {
	version := &lurk.Version{
		Type:  lurk.TypeVersion,
		Major: 2,
		Minor: 3,
	}

	g.game = &lurk.Game{
		Type:          lurk.TypeGame,
		InitialPoints: initialPoints,
		StatLimit:     initialPoints,
		GameDesc:      gameDescription,
	}

	ba, err := lurk.Marshal(version)
	if err != nil {
		return err
	}

	if _, err = conn.Write(ba); err != nil {
		return err
	}

	if ba, err = lurk.Marshal(g.game); err != nil {
		return err
	}

	_, err = conn.Write(ba)
	return err
}

func (g *game) sendRoom(room *room, player string, conn net.Conn) error {
	ba, err := lurk.Marshal(room.r)
	if err != nil {
		return err
	}
	if _, err = conn.Write(ba); err != nil {
		return err
	}

	if err = g.sendCharacters(room, player, conn); err != nil {
		return err
	}

	return g.sendConnections(room, conn)
}

func (g *game) sendCharacters(room *room, player string, conn net.Conn) (err error) {
	// all characters and monsters in that room
	var ba []byte
	for k, user := range g.users {
		// should we include current user?
		if k == player || user.c.RoomNum != room.r.RoomNumber {
			continue
		}
		if ba, err = lurk.Marshal(user.c); err != nil {
			return
		}
		if _, err = conn.Write(ba); err != nil {
			return
		}
	}

	for _, npc := range g.monsters {
		if npc.RoomNum != room.r.RoomNumber {
			continue
		}
		if ba, err = lurk.Marshal(npc); err != nil {
			return
		}
		if _, err = conn.Write(ba); err != nil {
			return
		}
	}
	return
}

func (g *game) sendConnections(room *room, conn net.Conn) (err error) {
	var ba []byte
	for _, connection := range room.connections {
		if ba, err = lurk.Marshal(connection); err != nil {
			return err
		}
		if _, err = conn.Write(ba); err != nil {
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

	ba, err := lurk.Marshal(accept)
	if err != nil {
		return err
	}
	_, err = conn.Write(ba)
	return err
}
