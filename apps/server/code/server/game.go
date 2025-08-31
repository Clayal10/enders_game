package server

import (
	"net"

	"github.com/Clayal10/enders_game/lib/cross"
	"github.com/Clayal10/enders_game/lib/lurk"
)

type game struct {
	// key is name, should be unique.
	users map[string]*lurk.Character
	// key is name? monster is a generic name for an npc
	monsters map[string]*lurk.Character
	// key is room number. Need to be careful about multithreading this
	rooms map[uint16]*lurk.Room
}

// default game values.
const (
	initialPoints = 100
)

func newGame() *game {
	g := &game{
		users:    make(map[string]*lurk.Character),
		monsters: make(map[string]*lurk.Character),
		rooms:    make(map[uint16]*lurk.Room),
	}

	g.rooms = createRooms()

	return g
}

func (g *game) sendStart(conn net.Conn) error {
	gameMessage := lurk.Game{
		Type:          lurk.TypeGame,
		InitialPoints: initialPoints,
		StatLimit:     initialPoints,
		GameDesc:      gameDescription,
	}

	msg, err := lurk.Marshal(&gameMessage)
	if err != nil {
		return err
	}

	_, err = conn.Write(msg)
	return err
}

func (g *game) registerPlayer(conn net.Conn) error {
	if err := g.addToUser(conn); err != nil {
		return err
	}

	buffer := make([]byte, bufferLength)

	var typ lurk.MessageType
	for typ != lurk.TypeStart {
		n, err := conn.Read(buffer)
		if err != nil {
			return err
		}

		msg, err := lurk.Unmarshal(buffer[:n])
		if err != nil {
			return err
		}
		typ = msg.GetType()
	}

	return nil
}

func (g *game) addToUser(conn net.Conn) error {
	buffer := make([]byte, bufferLength)

	// In this loop, we get the character and send it back after checking the validity of it.
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return err
		}

		if n >= bufferLength {
			n, err = lurk.GetVariableLength(buffer)
			if err != nil {
				return err
			}

			b := make([]byte, n-bufferLength)
			if _, err = conn.Read(b); err != nil {
				return err
			}
			buffer = append(buffer, b...)
		}

		msg, err := lurk.Unmarshal(buffer[:n])
		if err != nil {
			return err
		}
		if msg.GetType() != lurk.TypeCharacter {
			if err := g.sendError(conn, cross.Other, "You must send a [CHARACTER] type."); err != nil {
				return err
			}
			continue
		}

		character := msg.(*lurk.Character)
		if e := g.validateCharacter(character); e != cross.NoError {
			if err := g.sendError(conn, e, "Your [CHARACTER] has invalid stats."); err != nil {
				return err
			}
			continue
		}

		// Character is good at this point, flip flag and wait for their start.
		character.Flags[lurk.Ready] = true
		ba, err := lurk.Marshal(character)
		if err != nil {
			return err
		}

		if _, err = conn.Write(ba); err != nil {
			return err
		}

		g.users[character.Name] = character
		break
	}
	return nil
}

func (g *game) sendError(conn net.Conn, code cross.ErrCode, msg string) error {
	ba, err := lurk.Marshal(&lurk.Error{
		Type:       lurk.TypeError,
		ErrCode:    code,
		ErrMessage: msg,
	})
	if err != nil {
		return err
	}
	if _, err = conn.Write(ba); err != nil {
		return err
	}
	return nil
}

func (g *game) validateCharacter(c *lurk.Character) cross.ErrCode {
	if c.Gold != 0 || c.Health != 0 {
		return cross.StatError
	}
	if _, ok := g.users[c.Name]; ok {
		return cross.PlayerAlreadyExists
	}
	return cross.NoError
}
