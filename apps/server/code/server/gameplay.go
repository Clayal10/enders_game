package server

import (
	"net"

	"github.com/Clayal10/enders_game/lib/lurk"
)

// Descriptions / long text.
const (
	gameDescription = ` 
 ____  __ _  ____  ____  ____  _ ____     ___   __   _  _  ____ 
(  __)(  ( \(    \(  __)(  _ \(// ___)   / __) / _\ ( \/ )(  __)
 ) _) /    / ) D ( ) _)  )   /  \___ \  ( (_ \/    \/ \/ \ ) _) 
(____)\_)__)(____/(____)(__\_)  (____/   \___/\_/\_/\_)(_/(____)

The world has been ravaged by the most feared and despised being known to man, the formic. When it comes down to preventing their second massacre, will you be the one to step up and destroy them?`
)

// Entity names.
const (
	colonelGraph = "Colonel Graph"
)

// Room Numbers
const (
	battleSchool = 1
)

func (g *game) createRooms() {
	g.rooms = map[uint16]*lurk.Room{
		battleSchool: {
			Type:       lurk.TypeRoom,
			RoomNumber: battleSchool,
			RoomName:   "Battle School",
			RoomDesc:   "A place where young children play a game. At least, that is what the media says. The reality for this little ones is far bleaker than anyone could imagine.",
		},
	}
}

func (g *game) createMonsters() {
	g.monsters = map[string]*lurk.Character{
		colonelGraph: {
			Type:       lurk.TypeCharacter,
			Name:       colonelGraph,
			Attack:     20,
			Defense:    100,
			Regen:      100,
			Health:     100,
			Gold:       0,
			RoomNum:    battleSchool,
			PlayerDesc: "An older man, starting to let himself go, but sturdy non the less.",
		},
	}
}

func (g *game) handleMessage(msg *lurk.Message, conn net.Conn) {

}
func (g *game) handleChangeRoom(changeRoom *lurk.ChangeRoom, conn net.Conn) {

}
func (g *game) handleFight(conn net.Conn) {

}
func (g *game) handlePVPFight(pvp *lurk.PVPFight, conn net.Conn) {

}
func (g *game) handleLoot(loot *lurk.Loot, conn net.Conn) {

}
func (g *game) handleCharacter(char *lurk.Character, conn net.Conn) {

}
func (g *game) handleLeave(player string, conn net.Conn) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.users[player] = nil
	_ = g.sendAccept(conn, lurk.TypeLeave) // don't need it to work if they are exiting.
}
