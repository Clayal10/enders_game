package server

import (
	"fmt"
	"net"
	"sync"

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

	mu sync.Mutex
}

// default game values.
const (
	initialPoints = 100
)

// when creating a new game, we need to initialize the rooms and all entities.
func newGame() *game {
	g := &game{
		users:    make(map[string]*lurk.Character),
		monsters: make(map[string]*lurk.Character),
		rooms:    make(map[uint16]*lurk.Room),
	}

	g.createRooms()
	g.createMonsters()

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

func (g *game) registerPlayer(conn net.Conn) (string, error) {
	id, err := g.addUser(conn)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, bufferLength)

	for {
		n, err := conn.Read(buffer) // accept START
		if err != nil {
			return "", err
		}

		msg, err := lurk.Unmarshal(buffer[:n])
		if err != nil {
			return "", err
		}
		if msg.GetType() == lurk.TypeStart {
			if err = g.sendAccept(conn, lurk.TypeStart); err != nil { // accepted START
				return "", err
			}
			break
		}
		g.sendError(conn, cross.NotReady, "Please send a [START] message")
	}

	return id, nil
}

func (g *game) addUser(conn net.Conn) (characterID string, err error) {

	// In this loop, we get the character and send it back after checking the validity of it.
	for {
		buffer, n, err := g.readAll(conn) // accept CHARACTER
		if err != nil {
			return "", err
		}

		msg, err := lurk.Unmarshal(buffer[:n])
		if err != nil {
			return "", err
		}
		if msg.GetType() != lurk.TypeCharacter {
			if err := g.sendError(conn, cross.Other, "You must send a [CHARACTER] type."); err != nil {
				return "", err
			}
			continue
		}

		g.mu.Lock()
		character := msg.(*lurk.Character)
		if e := g.validateCharacter(character); e != cross.NoError {
			g.mu.Unlock()
			if err := g.sendError(conn, e, "Your [CHARACTER] has invalid stats."); err != nil {
				return "", err
			}
			continue
		}

		if err = g.sendAccept(conn, lurk.TypeCharacter); err != nil { // accepted CHARACTER
			return "", err
		}

		// Character is good at this point, flip flag and wait for their start.
		character.Flags[lurk.Ready] = true
		character.RoomNum = battleSchool
		g.users[character.Name] = character
		characterID = character.Name
		g.mu.Unlock()

		ba, err := lurk.Marshal(character)
		if err != nil {
			return "", err
		}

		if _, err = conn.Write(ba); err != nil {
			return "", err
		}

		break
	}
	return characterID, err
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

// An error returned from here results in termination of the client.
func (g *game) startGameplay(player string, conn net.Conn) error {
	// First, send the user information on their current room.
	if err := g.sendRoom(g.rooms[battleSchool], player, conn); err != nil {
		return err
	}

	for {
		buffer, n, err := g.readAll(conn) // accept MESSAGE || CHARACTER || LEAVE
		if err != nil {
			if err = g.sendError(conn, cross.Other, "Bad message, try again."); err != nil {
				return err
			}
			return err
		}

		lm, err := lurk.Unmarshal(buffer[:n])
		if err != nil {
			if err = g.sendError(conn, cross.Other, "Bad message, try again."); err != nil {
				return err
			}
			continue
		}

		// Not calling 'continue' since this is the last thing in the logic loop.
		switch lm.GetType() {
		case lurk.TypeMessage:
			msg, ok := lm.(*lurk.Message)
			if !ok {
				err = g.sendError(conn, cross.Other, fmt.Sprintf("Message had type %v, but had invalid fields.", lurk.TypeMessage))
				if err != nil {
					return err
				}
			}
			g.handleMessage(msg, conn)
		case lurk.TypeChangeRoom:
			msg, ok := lm.(*lurk.ChangeRoom)
			if !ok {
				err = g.sendError(conn, cross.Other, fmt.Sprintf("Message had type %v, but had invalid fields.", lurk.TypeChangeRoom))
				if err != nil {
					return err
				}
			}
			g.handleChangeRoom(msg, conn)
		case lurk.TypeFight:
			g.handleFight(conn)
		case lurk.TypePVPFight:
			msg, ok := lm.(*lurk.PVPFight)
			if !ok {
				err = g.sendError(conn, cross.Other, fmt.Sprintf("Message had type %v, but had invalid fields.", lurk.TypePVPFight))
				if err != nil {
					return err
				}
			}
			g.handlePVPFight(msg, conn)
		case lurk.TypeLoot:
			msg, ok := lm.(*lurk.Loot)
			if !ok {
				err = g.sendError(conn, cross.Other, fmt.Sprintf("Message had type %v, but had invalid fields.", lurk.TypeLoot))
				if err != nil {
					return err
				}
			}
			g.handleLoot(msg, conn)
		case lurk.TypeCharacter:
			msg, ok := lm.(*lurk.Character)
			if !ok {
				err = g.sendError(conn, cross.Other, fmt.Sprintf("Message had type %v, but had invalid fields.", lurk.TypeCharacter))
				if err != nil {
					return err
				}
			}
			g.handleCharacter(msg, conn)
		case lurk.TypeLeave:
			g.handleLeave(conn)

		}
	}
}

func (g *game) sendRoom(room *lurk.Room, player string, conn net.Conn) error {
	ba, err := lurk.Marshal(room)
	if err != nil {
		return err
	}
	if _, err = conn.Write(ba); err != nil {
		return err
	}

	// all characters and monsters in that room
	for k, v := range g.users {
		// should we include current user?
		if k == player {
			continue
		}
		if ba, err = lurk.Marshal(v); err != nil {
			return err
		}
		if _, err = conn.Write(ba); err != nil {
			return err
		}
	}

	for _, v := range g.monsters {
		if ba, err = lurk.Marshal(v); err != nil {
			return err
		}
		if _, err = conn.Write(ba); err != nil {
			return err
		}
	}
	return nil
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

// We want to read exactly the length of the message. This function will do up to 3
// calls to 'Read' to read exactly one message.
func (g *game) readAll(conn net.Conn) ([]byte, int, error) {
	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	if err != nil {
		return nil, 0, err
	}

	bytesNeeded, ok := lurk.LengthOffset[lurk.MessageType(buffer[0])]
	if !ok {
		return nil, 0, cross.ErrInvalidMessageType
	}

	if bytesNeeded == 1 {
		return buffer, 1, nil
	}

	b := make([]byte, bytesNeeded-1)
	if _, err = conn.Read(b); err != nil {
		return nil, 0, err
	}
	buffer = append(buffer, b...)

	n, err := lurk.GetVariableLength(buffer)
	if err != nil {
		return nil, 0, err
	}

	if n == -1 {
		return buffer, bytesNeeded, nil
	}

	b = make([]byte, n)
	if _, err = conn.Read(b); err != nil {
		return nil, 0, err
	}
	buffer = append(buffer, b...)
	return buffer, n, nil
}
