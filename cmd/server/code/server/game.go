package server

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"net"
	"sync"
	"time"

	"github.com/Clayal10/enders_game/pkg/cross"
	"github.com/Clayal10/enders_game/pkg/lurk"
)

// game holds all information and methods needed to create a game that complies to the
// lurk protocol actions.
type game struct {
	// key is name, should be unique.
	users map[string]*user
	// key is name? monster is a generic name for an npc
	monsters map[string]*lurk.Character
	// key is room number. Need to be careful about multithreading this
	rooms map[uint16]*room

	game    *lurk.Game
	version *lurk.Version

	mu           sync.Mutex
	lastActivity map[string]time.Time
	healTimer    map[string]*time.Timer
}

type user struct {
	c    *lurk.Character
	conn net.Conn
	// Key is room number. For conditional rooms. Users won't be able to see or access these rooms until true.
	allowedRoom map[uint16]bool
	killedQueen bool
	killedFleet bool
	terminated  bool
}

type room struct {
	r           *lurk.Room
	connections []*lurk.Connection
}

const (
	initialHealth = 100
	// gold values to unlock certain areas
	erosGold = 100
)

var errDisconnect = errors.New("disconnect")

// when creating a new game, we need to initialize the rooms and all entities.
func newGame() *game {
	g := &game{
		users:        make(map[string]*user),
		monsters:     make(map[string]*lurk.Character),
		rooms:        make(map[uint16]*room),
		lastActivity: make(map[string]time.Time),
		healTimer:    make(map[string]*time.Timer),
		version: &lurk.Version{
			Type:  lurk.TypeVersion,
			Major: 2,
			Minor: 3,
		},

		game: &lurk.Game{
			Type:          lurk.TypeGame,
			InitialPoints: initialPoints,
			StatLimit:     statLimit,
			GameDesc:      gameDescription,
		},
	}

	g.createRooms()
	g.createMonsters()

	return g
}

func (g *game) registerPlayer(conn net.Conn) (string, error) {
	id, err := g.addUser(conn)
	if err != nil {
		return id, err
	}
	log.Printf("Added user %v", id)

	for {
		buffer, _, err := lurk.ReadSingleMessage(conn) // accept START
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
		buffer, n, err := lurk.ReadSingleMessage(conn) // accept CHARACTER
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

		if _, err = conn.Write(lurk.Marshal(character)); err != nil {
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
	character.Flags[lurk.Monster] = false
	character.Flags[lurk.Alive] = true
	character.Flags[lurk.Started] = true

	character.RoomNum = battleSchool
	character.Health = initialHealth
	character.Gold = 0
	character.RoomNum = battleSchool
	u := &user{
		c:           character,
		conn:        conn,
		allowedRoom: make(map[uint16]bool),
	}
	for room := battleSchool; room <= rotterdam; room++ {
		u.allowedRoom[room] = true
	}
	for room := eros; room <= formicHomeWorld; room++ {
		if character.Name == "Beans Shumaker" {
			u.allowedRoom[room] = true
		} else {
			u.allowedRoom[room] = false
		}
	}

	g.users[character.Name] = u
	return character.Name
}

func (g *game) validateCharacter(c *lurk.Character) cross.ErrCode {
	if c.Attack+c.Defense+c.Regen > g.game.InitialPoints {
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
	g.mu.Lock()
	if err := g.sendRoom(g.rooms[battleSchool], player, conn); err != nil {
		g.mu.Unlock()
		return err
	}
	if err := g.notifyNewArrival(player); err != nil {
		g.mu.Unlock()
		return err
	}
	g.mu.Unlock()

	for {
		g.mu.Lock()
		user, ok := g.users[player]
		g.mu.Unlock()
		if !ok { // User has been removed / left.
			return nil
		}

		buffer, n, err := lurk.ReadSingleMessage(conn) // accept MESSAGE || CHARACTER || LEAVE
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
			if err := g.checkStatusChange(user, conn); err != nil {
				return err
			}
			continue
		}

		// The message did not have proper fields for the message type.
		if err = g.sendError(conn, cross.Other, fmt.Sprintf("Message contains invalid fields for type %d", lm.GetType())); err != nil {
			return err
		}
	}
}

func (g *game) notifyNewArrival(newbie string) error {
	newUser, ok := g.users[newbie]
	if !ok {
		return cross.ErrUserNotInServer
	}

	for _, otherUser := range g.users {
		if otherUser.c.RoomNum != battleSchool || otherUser.c.Name == newUser.c.Name {
			continue
		}
		if err := g.sendCharacterUpdate(newUser.c, otherUser.conn, otherUser.c.Name,
			fmt.Sprintf("%s joined battle school!", newUser.c.Name)); err != nil {
			log.Printf("Could not send message to %s", otherUser.c.Name)
		}
	}

	return nil
}

const upgradeCost = 50

// A chance to update character stats after each action.
func (g *game) checkStatusChange(user *user, conn net.Conn) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if err := askForUpgrade(user); err != nil {
		return err
	}

	status := map[uint16]bool{}
	maps.Copy(status, user.allowedRoom)
	if user.c.Gold > erosGold {
		user.allowedRoom[eros] = true
	}
	if user.killedQueen {
		user.allowedRoom[earth] = true
		user.allowedRoom[shakespeare] = true
	}
	if user.killedFleet {
		user.allowedRoom[formicHomeWorld] = true
	}

	if user.allowedRoom[eros] == status[eros] && // no change, don't send update
		user.allowedRoom[formicHomeWorld] == status[formicHomeWorld] &&
		user.allowedRoom[earth] == status[earth] &&
		user.allowedRoom[shakespeare] == status[shakespeare] {
		return nil
	}
	return g.sendConnections(g.rooms[user.c.RoomNum], user.c.Name, conn)
}

func askForUpgrade(user *user) (err error) {
	if user.c.Gold >= upgradeCost && user.c.RoomNum == battleSchoolBarracks {
		_, err = user.conn.Write(lurk.Marshal(&lurk.Message{
			Recipient: user.c.Name,
			Sender:    narrator,
			Text: fmt.Sprintf(
				"Looks like some of your hard work is paying off, spend %d gold to upgrade your stats. (Message %s to increase all stats by 5 points)",
				upgradeCost, narrator),
			Narration: true,
		}))
	}
	return
}

func (g *game) messageSelection(lm lurk.LurkMessage, player string, conn net.Conn) (err error, _ bool) {
	switch lm.GetType() {
	case lurk.TypeMessage:
		msg := lm.(*lurk.Message)
		err = g.handleMessage(msg, conn)
	case lurk.TypeChangeRoom:
		msg := lm.(*lurk.ChangeRoom)
		err = g.handleChangeRoom(msg, conn, player)
	case lurk.TypeFight:
		err = g.handleFight(conn, player)
	case lurk.TypePVPFight:
		msg := lm.(*lurk.PVPFight)
		err = g.handlePVPFight(msg, conn, player)
	case lurk.TypeLoot:
		msg := lm.(*lurk.Loot)
		err = g.handleLoot(conn, msg, player)
	case lurk.TypeLeave:
		g.handleLeave(player)
		return errDisconnect, false
	default:
		return nil, false
	}
	return err, true
}

var monsterHealTime = 10 * time.Second

func (g *game) startHealTimer(monster *lurk.Character) {
	if g.healTimer[monster.Name] == nil {
		g.healTimer[monster.Name] = time.AfterFunc(monsterHealTime, func() {
			g.healMonster(monster)
		})
	} else {
		g.healTimer[monster.Name].Reset(monsterHealTime)
	}
}

func (g *game) healMonster(monster *lurk.Character) {
	g.mu.Lock()
	defer g.mu.Unlock()

	idle := time.Since(g.lastActivity[monster.Name])
	if idle < monsterHealTime {
		return
	}
	monster.Health = monsterHealth[monster.Name]
	monster.Flags[lurk.Alive] = true
	for _, user := range g.users {
		if user.c.RoomNum != monster.RoomNum {
			continue
		}
		if err := g.sendCharacterUpdate(monster, user.conn, user.c.Name, ""); err != nil {
			log.Printf("%s: could not update user %v with updated monster health", err.Error(), user.c.Name)
		}
	}
}

func (g *game) upgradeStats(user *user, conn net.Conn) error {
	if user.c.Gold < upgradeCost || user.c.RoomNum != battleSchoolBarracks {
		return g.sendError(conn, cross.StatError, fmt.Sprintf(
			"You must be in the Battle School Barracks with at least %d gold to upgrade your stats", upgradeCost))
	}
	if totalStats := user.c.Attack + user.c.Defense + user.c.Regen; totalStats-15 > statLimit {
		return g.sendError(conn, cross.StatError, fmt.Sprintf(
			"Your stat sum of %d is too high to upgrade any further", totalStats))
	}
	user.c.Attack += 5
	user.c.Defense += 5
	user.c.Regen += 5
	user.c.Gold -= upgradeCost

	for _, u := range g.users {
		_, _ = u.conn.Write(lurk.Marshal(user.c))
	}
	return nil
}
