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
}

// default game values.
const (
	initialPoints   = 100
	gameDescription = ` 
 ____  __ _  ____  ____  ____  _ ____     ___   __   _  _  ____ 
(  __)(  ( \(    \(  __)(  _ \(// ___)   / __) / _\ ( \/ )(  __)
 ) _) /    / ) D ( ) _)  )   /  \___ \  ( (_ \/    \/ \/ \ ) _) 
(____)\_)__)(____/(____)(__\_)  (____/   \___/\_/\_/\_)(_/(____)

The world has been ravaged by the most feared and despised being known to man, the formic. When it comes down to preventing their second massacre, will you be the one to step up and destroy them?`
)

func newGame() *game {
	return &game{
		users:    make(map[string]*lurk.Character),
		monsters: make(map[string]*lurk.Character),
	}
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
			conn.Read(b)
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
