package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/Clayal10/enders_game/lib/cross"
	"github.com/Clayal10/enders_game/lib/lurk"
)

type game struct {
	// key is name, should be unique.
	users map[string]*user
	// key is name? monster is a generic name for an npc
	monsters map[string]*lurk.Character
	// key is room number. Need to be careful about multithreading this
	rooms map[uint16]*room

	game *lurk.Game

	mu sync.Mutex
}

type user struct {
	c    *lurk.Character
	conn net.Conn
	// Key is room number. For conditional rooms. Users won't be able to see or access these rooms until true.
	allowedRoom map[uint16]bool
}

type room struct {
	r           *lurk.Room
	connections []*lurk.Connection
}

// default game values.
const (
	initialPoints = 100
)

// when creating a new game, we need to initialize the rooms and all entities.
func newGame() *game {
	g := &game{
		users:    make(map[string]*user),
		monsters: make(map[string]*lurk.Character),
		rooms:    make(map[uint16]*room),
	}

	g.createRooms()
	g.createMonsters()

	return g
}

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

func (g *game) registerPlayer(conn net.Conn) (string, error) {
	id, err := g.addUser(conn)
	if err != nil {
		return id, err
	}
	log.Printf("Added user %v", id)

	for {
		buffer, _, err := readSingleMessage(conn) // accept START
		if err != nil {
			return id, err
		}
		msg, err := lurk.Unmarshal(buffer)
		if err != nil {
			return id, err
		}
		if msg.GetType() == lurk.TypeStart {
			if err = g.sendAccept(conn, lurk.TypeStart); err != nil { // accepted START
				return id, err
			}
			break
		}
		if err = g.sendError(conn, cross.NotReady, "Please send a [START] message"); err != nil {
			return id, err
		}
	}

	return id, nil
}

func (g *game) addUser(conn net.Conn) (characterID string, err error) {
	// In this loop, we get the character and send it back after checking the validity of it.
	for {
		buffer, n, err := readSingleMessage(conn) // accept CHARACTER
		if err != nil {
			_ = g.sendError(conn, cross.Other, "Bad message, terminating connection.")
			return characterID, err
		}

		msg, err := lurk.Unmarshal(buffer[:n])
		if err != nil {
			return characterID, err
		}
		if msg.GetType() != lurk.TypeCharacter {
			if err := g.sendError(conn, cross.Other, "You must send a [CHARACTER] type."); err != nil {
				return characterID, err
			}
			fmt.Println("Not type character")
			continue
		}

		g.mu.Lock()
		character := msg.(*lurk.Character)
		if e := g.validateCharacter(character); e != cross.NoError {
			g.mu.Unlock()
			if err := g.sendError(conn, e, "Your [CHARACTER] has invalid stats"); err != nil {
				return characterID, err
			}
			continue
		}

		characterID = g.createUser(character, conn)
		g.mu.Unlock()

		ba, err := lurk.Marshal(character)
		if err != nil {
			return characterID, err
		}

		if _, err = conn.Write(ba); err != nil {
			return characterID, err
		}

		if err = g.sendAccept(conn, lurk.TypeCharacter); err != nil { // accepted CHARACTER
			return characterID, err
		}

		break
	}
	return characterID, err
}

func (g *game) createUser(character *lurk.Character, conn net.Conn) string {
	// Character is good at this point, flip flag and wait for their start.
	character.Flags[lurk.Ready] = true
	character.RoomNum = battleSchool
	g.users[character.Name] = &user{
		c:           character,
		conn:        conn,
		allowedRoom: make(map[uint16]bool),
	}
	return character.Name
}

func (g *game) sendError(conn net.Conn, code cross.ErrCode, msg string) error {
	ba, err := lurk.Marshal(&lurk.Error{
		Type:       lurk.TypeError,
		ErrCode:    code,
		ErrMessage: msg,
	})
	if err == nil {
		_, err = conn.Write(ba)
	}
	return err
}

func (g *game) validateCharacter(c *lurk.Character) cross.ErrCode {
	if c.Attack >= g.game.StatLimit ||
		c.Defense >= g.game.StatLimit ||
		c.Regen >= g.game.StatLimit {
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
		if _, ok := g.users[player]; !ok { // User has been removed.
			return nil
		}

		buffer, n, err := readSingleMessage(conn) // accept MESSAGE || CHARACTER || LEAVE
		if err != nil {
			_ = g.sendError(conn, cross.Other, "Bad message, try again.")
			return err
		}

		lm, err := lurk.Unmarshal(buffer[:n])
		if err != nil {
			_ = g.sendError(conn, cross.Other, "Bad message, try again.")
			return err
		}

		if err, ok := g.messageSelection(lm, player, conn); err != nil {
			return err
		} else if ok {
			continue
		}
		// The message did not have proper fields for the message type.
		if err = g.sendError(conn, cross.Other, fmt.Sprintf("Message contains invalid fields for type %d", lm.GetType())); err != nil {
			return err
		}
	}
}

func (g *game) messageSelection(lm lurk.LurkMessage, player string, conn net.Conn) (err error, _ bool) {
	switch lm.GetType() {
	case lurk.TypeMessage:
		msg, ok := lm.(*lurk.Message)
		if !ok {
			return nil, ok
		}
		g.handleMessage(msg, player)
	case lurk.TypeChangeRoom:
		msg, ok := lm.(*lurk.ChangeRoom)
		if !ok {
			return nil, ok
		}
		err = g.handleChangeRoom(msg, conn, player)
	case lurk.TypeFight:
		msg, ok := lm.(*lurk.Fight)
		if !ok {
			return nil, ok
		}
		g.handleFight(msg, player)
	case lurk.TypePVPFight:
		msg, ok := lm.(*lurk.PVPFight)
		if !ok {
			return nil, ok
		}
		g.handlePVPFight(msg, player)
	case lurk.TypeLoot:
		msg, ok := lm.(*lurk.Loot)
		if !ok {
			return nil, ok
		}
		g.handleLoot(msg, player)
	case lurk.TypeCharacter:
		msg, ok := lm.(*lurk.Character)
		if !ok {
			return nil, ok
		}
		g.handleCharacter(msg, player)
	case lurk.TypeLeave:
		g.handleLeave(player)
	default:
		return nil, false
	}
	return err, true
}

func (g *game) sendRoom(room *room, player string, conn net.Conn) error {
	ba, err := lurk.Marshal(room.r)
	if err != nil {
		return err
	}
	if _, err = conn.Write(ba); err != nil {
		return err
	}
	// all characters and monsters in that room
	for k, user := range g.users {
		// should we include current user?
		if k == player || user.c.RoomNum != room.r.RoomNumber {
			continue
		}
		if ba, err = lurk.Marshal(user.c); err != nil {
			return err
		}
		if _, err = conn.Write(ba); err != nil {
			return err
		}
	}

	for _, npc := range g.monsters {
		if npc.RoomNum != room.r.RoomNumber {
			continue
		}
		if ba, err = lurk.Marshal(npc); err != nil {
			return err
		}
		if _, err = conn.Write(ba); err != nil {
			return err
		}
	}

	for _, connection := range room.connections {
		if ba, err = lurk.Marshal(connection); err != nil {
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
