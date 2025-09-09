package server

import (
	"fmt"
	"net"

	"github.com/Clayal10/enders_game/lib/cross"
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
	battleSchoolDesc           = "A place where young children play a game. At least, that is what the media says. The reality is that they will manipulate and contort our lives just to see what we can handle."
	battleSchoolBarracksDesc   = "The room filled with small children, most of them scared, but none of them trying to show their weakness. It looks like the remaining bunk is here at the front, lucky me."
	battleSchoolGameRoomDesc   = "Many older boys are hunched over the game table, none pay much attention to my entrance."
	battleSchoolBattleRoomDesc = "A room, 100 cubic meters in size, defying the laws of gravity. With a gate on either side of the room, us children are able to wage war against each other for honor, all the while perfecting zero G movement. "
	formicStarSystemDesc       = "Out here in the cold, dark vastness of space, a world filled with billions of alien life forms lay idle."
	formicHomeWorldDesc        = "In all of the universe, one could not find a more perfect machine working under the surface of this planet. The queen instructs, and the workers follow. Flawlessly. To see this creature is to be in awe and trembling fear at the same time."
	erosDesc                   = "The secret base for International Fleet Command operations. The surface is blacked out, covered in solar panels. The inhabitants stay below the surface in the smooth tunnels crafted by the formic race many years ago."
	earthDesc                  = "A world doomed. A planet that needs a savior. To go back now is to let the wretched Formics win."
)

// Entity names.
const (
	// Friends
	colonelGraph = "Colonel Graph"
	bean         = "Bean"
	// Enemies
	bonzo       = "Bonito de Madrid"
	formicFleet = "Formic Fleet"
	hiveQueen   = "Hive Queen"
)

// Room Numbers
const (
	battleSchool           = 1 // Central hub / hallways / entrance / exit for battle school.
	battleSchoolBarracks   = 2
	battleSchoolGameRoom   = 3
	battleSchoolBattleRoom = 4
	eros                   = 5 // Hidden until defeating bonzo
	shakespeare            = 6 // Hidden until defeating formics.
	formicStarSystem       = 7
	formicHomeWorld        = 8
	earth                  = 9 // Hidden until defeating or losing to bonzo.
)

func (g *game) createRooms() {
	g.rooms = map[uint16]*room{
		battleSchool: {
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: battleSchool,
				RoomName:   "Battle School",
				RoomDesc:   battleSchoolDesc,
			},
			connections: []*lurk.Connection{
				{
					Type:       lurk.TypeConnection,
					RoomNumber: battleSchoolBarracks,
					RoomName:   "The Barracks",
					RoomDesc:   battleSchoolBarracksDesc,
				},
				{
					Type:       lurk.TypeConnection,
					RoomNumber: battleSchoolGameRoom,
					RoomName:   "The Game Room",
					RoomDesc:   battleSchoolGameRoomDesc,
				},
				{
					Type:       lurk.TypeConnection,
					RoomNumber: battleSchoolBattleRoom,
					RoomName:   "The Battle Room",
					RoomDesc:   battleSchoolBattleRoomDesc,
				},
				{
					Type:       lurk.TypeConnection,
					RoomNumber: eros,
					RoomName:   "Eros",
					RoomDesc:   erosDesc,
				},
			},
		},
		battleSchoolBarracks: {
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: battleSchoolBarracks,
				RoomName:   "The Barracks",
				RoomDesc:   battleSchoolBarracksDesc,
			},
			connections: []*lurk.Connection{
				{
					Type:       lurk.TypeConnection,
					RoomNumber: battleSchool,
					RoomName:   "Battle School",
					RoomDesc:   battleSchoolDesc,
				},
			},
		},
		battleSchoolGameRoom: {
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: battleSchoolGameRoom,
				RoomName:   "The Game Room",
				RoomDesc:   battleSchoolGameRoomDesc,
			},
			connections: []*lurk.Connection{
				{
					Type:       lurk.TypeConnection,
					RoomNumber: battleSchool,
					RoomName:   "Battle School",
					RoomDesc:   battleSchoolDesc,
				},
			},
		},
		battleSchoolBattleRoom: {
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: battleSchoolBattleRoom,
				RoomName:   "The Battle Room",
				RoomDesc:   battleSchoolBattleRoomDesc,
			},
			connections: []*lurk.Connection{
				{
					Type:       lurk.TypeConnection,
					RoomNumber: battleSchool,
					RoomName:   "Battle School",
					RoomDesc:   battleSchoolDesc,
				},
			},
		},
		formicStarSystem: {
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: formicStarSystem,
				RoomName:   "Formic Star System",
				RoomDesc:   formicStarSystemDesc,
			},
			connections: []*lurk.Connection{
				{
					Type:       lurk.TypeConnection,
					RoomNumber: formicHomeWorld,
					RoomName:   "Formic Home World",
					RoomDesc:   formicHomeWorldDesc,
				},
			},
		},
		eros: {
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: eros,
				RoomName:   "Eros",
				RoomDesc:   erosDesc,
			},
		},
		earth: {
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: earth,
				RoomName:   "Earth",
				RoomDesc:   earthDesc,
			},
		},
	}
}
func (g *game) createMonsters() {
	g.monsters = map[string]*lurk.Character{
		colonelGraph: {
			Type: lurk.TypeCharacter,
			Name: colonelGraph,
			Flags: map[string]bool{
				lurk.Alive:   true,
				lurk.Monster: true,
			},
			Attack:     20,
			Defense:    100,
			Regen:      100,
			Health:     100,
			Gold:       0,
			RoomNum:    battleSchool,
			PlayerDesc: "An older man, starting to let himself go, but sturdy non the less.",
		},
		bean: {
			Type: lurk.TypeCharacter,
			Name: bean,
			Flags: map[string]bool{
				lurk.Alive:   true,
				lurk.Monster: true, // Maybe monster?
			},
			Attack:     10,
			Defense:    200,
			Regen:      100,
			Health:     100,
			Gold:       0,
			RoomNum:    battleSchoolBarracks,
			PlayerDesc: "The littlest one in battle school. You would be mistaken to think that is an indication of his power, though.",
		},
	}
}

func (g *game) handleMessage(msg *lurk.Message, player string) {}
func (g *game) handleChangeRoom(changeRoom *lurk.ChangeRoom, conn net.Conn, player string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	user, ok := g.users[player]
	if !ok {
		return g.sendError(conn, cross.Other, fmt.Sprintf("%v: error in changing room", cross.ErrUserNotInServer.Error()))
	}

	currentRoom := g.rooms[user.c.RoomNum]
	hasConnection := false
	for _, connection := range currentRoom.connections {
		hasConnection = connection.RoomNumber == changeRoom.RoomNumber
		if hasConnection {
			break
		}
	}

	if !hasConnection {
		return g.sendError(conn, cross.BadRoom, fmt.Sprintf("%v: error in changing room", cross.ErrRoomsNotConnected.Error()))
	}

	room, ok := g.rooms[changeRoom.RoomNumber]
	if !ok {
		return g.sendError(conn, cross.BadRoom, fmt.Sprintf("%v: error in changing room", cross.ErrInvalidRoomNumber.Error()))
	}

	user.c.RoomNum = room.r.RoomNumber
	return g.sendRoom(room, player, conn)
}
func (g *game) handleFight(fight *lurk.Fight, player string)        {}
func (g *game) handlePVPFight(pvp *lurk.PVPFight, player string)    {}
func (g *game) handleLoot(loot *lurk.Loot, player string)           {}
func (g *game) handleCharacter(char *lurk.Character, player string) {}
func (g *game) handleLeave(player string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.users, player)
}
