package server

import (
	"fmt"
	"log"
	"net"
	"time"

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
	battleSchoolBarracksDesc   = "The room filled with small children, most of them scared, but none of them trying to show their weakness."
	battleSchoolGameRoomDesc   = "Many older boys are hunched over the game table, just trying to show off to each other. You may be able to gain some experience if someone would give you the chance."
	battleSchoolBattleRoomDesc = "A room, 100 cubic meters in size, defying the laws of gravity. With a gate on either side of the room, us children are able to wage war against each other for honor, all the while practicing zero G movement."
	formicStarSystemDesc       = "Out here in the cold, dark vastness of space, a world filled with billions of alien life forms lay idle."
	formicHomeWorldDesc        = "In all of the universe, one could not find a more perfect machine working under the surface of this planet. The queen instructs, and the workers follow. Flawlessly. To see this creature is to be in awe and trembling fear at the same time."
	erosDesc                   = "The secret base for International Fleet Command operations. The surface is blacked out, covered in solar panels. The inhabitants stay below the surface in the smooth tunnels crafted by the formic race many years ago."
	earthDesc                  = "A world doomed. A planet that needs a savior. To go back now is to let the wretched Formics win."
	shakespeareDesc            = "The next frontier for human expansion. With the buggers eliminated, we can take their land and breed the next generation of humans and crops."
	rotterdamDesc              = "A city of ruins. The streets are filled with starved children fighting to the death."
)

// Entity names.
const (
	// Friends
	colonelGraph = "Colonel Graph"
	bean         = "Bean"
	mazer        = "Mazer Rackham"
	petra        = "Petra Arkanian"
	// Enemies
	bonzo       = "Bonito de Madrid"
	formicFleet = "Formic Fleet"
	achilles    = "Achilles de Flandres"
	peter       = "Peter Wiggin"
	hiveQueen   = "Hive Queen"
	// both
	hiveQueenCocoon = "Hive Queen Cacoon"
)

// Room Numbers
const (
	battleSchool           uint16 = 1 // Central hub / hallways / entrance / exit for battle school.
	battleSchoolBarracks   uint16 = 2
	battleSchoolGameRoom   uint16 = 3
	battleSchoolBattleRoom uint16 = 4
	formicStarSystem       uint16 = 5
	rotterdam              uint16 = 6

	eros            uint16 = 11 // Hidden until defeating bonzo
	shakespeare     uint16 = 12 // Hidden until defeating formics.
	earth           uint16 = 13 // Hidden until defeating or losing to bonzo.
	formicHomeWorld uint16 = 14
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
				{ // secret.
					Type:       lurk.TypeConnection,
					RoomNumber: eros,
					RoomName:   "Eros",
					RoomDesc:   erosDesc,
				},
				{ // secret.
					Type:       lurk.TypeConnection,
					RoomNumber: earth,
					RoomName:   "Earth",
					RoomDesc:   earthDesc,
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
				{
					Type:       lurk.TypeConnection,
					RoomNumber: shakespeare,
					RoomName:   "Shakespeare Colony",
					RoomDesc:   shakespeareDesc,
				},
				{
					Type:       lurk.TypeConnection,
					RoomNumber: eros,
					RoomName:   "Eros",
					RoomDesc:   erosDesc,
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
			connections: []*lurk.Connection{ // secret.
				{
					Type:       lurk.TypeConnection,
					RoomNumber: shakespeare,
					RoomName:   "Shakespeare Colony",
					RoomDesc:   shakespeareDesc,
				},
				{
					Type:       lurk.TypeConnection,
					RoomNumber: formicStarSystem,
					RoomName:   "Formic Star System",
					RoomDesc:   formicStarSystemDesc,
				},
				{
					Type:       lurk.TypeConnection,
					RoomNumber: battleSchool,
					RoomName:   "Battle School",
					RoomDesc:   battleSchoolDesc,
				},
			},
		},
		earth: { // No escape.
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: earth,
				RoomName:   "Earth",
				RoomDesc:   earthDesc,
			},
			connections: []*lurk.Connection{
				{
					Type:       lurk.TypeConnection,
					RoomNumber: rotterdam,
					RoomName:   "Rotterdam, The Netherlands",
					RoomDesc:   rotterdamDesc,
				},
			},
		},
		rotterdam: {
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: rotterdam,
				RoomName:   "Rotterdam, The Netherlands",
				RoomDesc:   rotterdamDesc,
			},
			connections: []*lurk.Connection{
				{
					Type:       lurk.TypeConnection,
					RoomNumber: earth,
					RoomName:   "Earth",
					RoomDesc:   earthDesc,
				},
			},
		},
		shakespeare: {
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: shakespeare,
				RoomName:   "Shakespeare Colony",
				RoomDesc:   shakespeareDesc,
			},
		},
		formicHomeWorld: {
			r: &lurk.Room{
				Type:       lurk.TypeRoom,
				RoomNumber: formicHomeWorld,
				RoomName:   "Formic Home World",
				RoomDesc:   formicHomeWorldDesc,
			},
			connections: []*lurk.Connection{
				{
					Type:       lurk.TypeConnection,
					RoomNumber: formicStarSystem,
					RoomName:   "Formic Star System",
					RoomDesc:   formicStarSystemDesc,
				},
			},
		},
	}
}

var monsterHealth = map[string]int16{
	colonelGraph:    50,
	bean:            100,
	petra:           100,
	mazer:           100,
	bonzo:           75,
	formicFleet:     10000,
	hiveQueen:       1000,
	achilles:        1000,
	peter:           1000,
	hiveQueenCocoon: 1,
}

// The amount of gold you gain by defeating each monster
var monsterGold = map[string]uint16{
	colonelGraph:    10,
	bean:            15,
	petra:           20,
	mazer:           100,
	bonzo:           50,
	formicFleet:     1000,
	hiveQueen:       1000,
	achilles:        1000,
	peter:           1000,
	hiveQueenCocoon: 64535,
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
			Health:     50,
			Gold:       0,
			RoomNum:    battleSchool,
			PlayerDesc: "An older man, starting to let himself go, but sturdy non the less.",
		},
		bean: {
			Type: lurk.TypeCharacter,
			Name: bean,
			Flags: map[string]bool{
				lurk.Alive:   true,
				lurk.Monster: true,
			},
			Attack:     10,
			Defense:    100,
			Regen:      100,
			Health:     100,
			Gold:       0,
			RoomNum:    battleSchoolBattleRoom,
			PlayerDesc: "The littlest one in battle school. You would be mistaken to think that is an indication of his power, though.",
		},
		petra: {
			Type: lurk.TypeCharacter,
			Name: petra,
			Flags: map[string]bool{
				lurk.Alive:   true,
				lurk.Monster: true,
			},
			Attack:     20,
			Defense:    80,
			Regen:      100,
			Health:     100,
			Gold:       0,
			RoomNum:    battleSchoolGameRoom,
			PlayerDesc: "The only girl in battle school, but she can be more dangerous that most of the boys. She could be an important teacher at this point.",
		},
		mazer: {
			Type: lurk.TypeCharacter,
			Name: mazer,
			Flags: map[string]bool{
				lurk.Alive:   true,
				lurk.Monster: true,
			},
			Attack:     100,
			Defense:    100,
			Regen:      0,
			Health:     100,
			Gold:       0,
			RoomNum:    eros,
			PlayerDesc: "Once believed to be dead, the greatest commander in all of history has shown up again. It seems his only intention is to train the next great commander of history. He will accomplish his goal or kill someone in the process.",
		},
		bonzo: {
			Type: lurk.TypeCharacter,
			Name: bonzo,
			Flags: map[string]bool{
				lurk.Alive:   true,
				lurk.Monster: true,
			},
			Attack:     100,
			Defense:    50,
			Regen:      50,
			Health:     75,
			Gold:       0,
			RoomNum:    battleSchoolBattleRoom,
			PlayerDesc: "Benito de Madrid; pretty boy. He will fight till the death for his families honor. To cross Bonzo is to can be the worst mistake you will make in your potentially short life.",
		},
		formicFleet: {
			Type: lurk.TypeCharacter,
			Name: formicFleet,
			Flags: map[string]bool{
				lurk.Alive:   true,
				lurk.Monster: true,
			},
			Attack:     50,
			Defense:    50,
			Regen:      0,
			Health:     1000,
			Gold:       0,
			RoomNum:    formicStarSystem,
			PlayerDesc: "A fleet of not thousands, or tens of thousands, but millions of individual formic creatures. They seems to move as if instructed by a single mind, perhaps a queen.",
		},
		hiveQueen: {
			Type: lurk.TypeCharacter,
			Name: hiveQueen,
			Flags: map[string]bool{
				lurk.Alive:   true,
				lurk.Monster: true,
			},
			Attack:     0,
			Defense:    0,
			Regen:      0,
			Health:     1000,
			Gold:       0,
			RoomNum:    formicHomeWorld,
			PlayerDesc: "The epitome of beauty and horror. There isn't a more terrifying creature imaginable by man. All the propaganda back on earth does not do justice to the fear that this creature invokes in one's heart. At the same time though, there is nothing more beautiful. You can feel her presence in your own, her mind in yours. To kill this creature is to kill your own self.",
		},
		achilles: {
			Type: lurk.TypeCharacter,
			Name: achilles,
			Flags: map[string]bool{
				lurk.Alive:   true,
				lurk.Monster: true,
			},
			Attack:     100,
			Defense:    100,
			Regen:      50,
			Health:     1000,
			Gold:       0,
			RoomNum:    rotterdam,
			PlayerDesc: "This boy seems to have taken control of the streets. Starving children cling to him as their papa. However, few claim he is must more than that...",
		},
		peter: {
			Type: lurk.TypeCharacter,
			Name: peter,
			Flags: map[string]bool{
				lurk.Alive:   true,
				lurk.Monster: true,
			},
			Attack:     100,
			Defense:    100,
			Regen:      50,
			Health:     1000,
			Gold:       0,
			RoomNum:    earth,
			PlayerDesc: "The boy who will take over the world. Peter will gain control of all those in his grasp, will you be his enemy or foe?",
		},
		hiveQueenCocoon: {
			Type: lurk.TypeCharacter,
			Name: hiveQueenCocoon,
			Flags: map[string]bool{
				lurk.Alive: true,
			},
			Attack:     0,
			Defense:    0,
			Regen:      0,
			Health:     1,
			Gold:       0,
			RoomNum:    shakespeare,
			PlayerDesc: "The next hive queen. Will you restore their race?",
		},
	}
}

const defaultWriteTimeout = time.Second

func (g *game) handleMessage(msg *lurk.Message, conn net.Conn) (err error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	recipient, ok := g.users[msg.Recipient]
	if !ok {
		return g.sendError(conn, cross.NoTarget, fmt.Sprintf("%v: error in sending message", cross.ErrUserNotInServer.Error()))
	}

	if err := recipient.conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout)); err != nil {
		return err
	}
	if _, err = recipient.conn.Write(lurk.Marshal(msg)); err == nil {
		log.Printf("%s sent message to %s\n", msg.Sender, msg.Recipient)
	}
	return err
}

func (g *game) handleChangeRoom(changeRoom *lurk.ChangeRoom, conn net.Conn, player string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	user, ok := g.users[player]
	// Checks for user.
	if !ok {
		return g.sendError(conn, cross.Other, fmt.Sprintf("%v: error in changing room", cross.ErrUserNotInServer.Error()))
	}

	// Check for valid request.
	currentRoom := g.rooms[user.c.RoomNum]
	hasConnection := false
	for _, connection := range currentRoom.connections {
		if hasConnection = connection.RoomNumber == changeRoom.RoomNumber &&
			user.allowedRoom[changeRoom.RoomNumber]; hasConnection {
			break
		}
	}

	if !hasConnection {
		return g.sendError(conn, cross.BadRoom, fmt.Sprintf("%v: error in changing room", cross.ErrRoomsNotConnected.Error()))
	}

	newRoom, ok := g.rooms[changeRoom.RoomNumber]
	if !ok {
		return g.sendError(conn, cross.BadRoom, fmt.Sprintf("%v: error in changing room", cross.ErrInvalidRoomNumber.Error()))
	}

	// Send new room to user.
	if user.c.RoomNum = newRoom.r.RoomNumber; user.c.RoomNum == battleSchoolBarracks {
		user.c.Flags[lurk.Alive] = true
		user.c.Health = initialHealth
	}

	if err := g.sendRoom(newRoom, player, conn); err != nil {
		return err
	}

	// Message others in the room that they have left and those in the room they are going to.
	for name, u := range g.users {
		msg := ""
		if rn := currentRoom.r.RoomNumber; u.c.RoomNum == rn {
			if !u.allowedRoom[rn] {
				msg = fmt.Sprintf("%s has been sent orders out of here.", user.c.Name)
			}
			if err := g.sendCharacterUpdate(user.c, u.conn, name, msg); err != nil {
				log.Printf("%s: error when sending character updates to %s", err, name)
			}
			// NOTE: This will send an updated character to the user.
		} else if rn := newRoom.r.RoomNumber; u.c.RoomNum == rn {
			if err := g.sendCharacterUpdate(user.c, u.conn, name, ""); err != nil {
				log.Printf("%s: error when sending character updates to %s", err, name)
			}
		}
	}

	return nil
}

func (g *game) handleFight(conn net.Conn, player string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	user, ok := g.users[player]
	if !ok {
		return g.sendError(conn, cross.Other, fmt.Sprintf("%v: error in fighting", cross.ErrUserNotInServer.Error()))
	}

	if !user.c.Flags[lurk.Alive] {
		return g.sendError(conn, cross.NoFight, player+", you cannot fight when you are dead")
	}

	currentRoom := g.rooms[user.c.RoomNum]

	var fights uint16 = 0
	for _, monster := range g.monsters {
		if monster.RoomNum != user.c.RoomNum || !monster.Flags[lurk.Alive] {
			continue
		}
		if monster.Name == hiveQueenCocoon {
			return g.sendError(conn, cross.Other, "If you wish to destroy the next hive queen, you must PVP fight.")
		}

		g.lastActivity[monster.Name] = time.Now()
		lurk.CalculateFight(user.c, monster)
		fights++

		if user.c.Flags[lurk.Alive] {
			user.c.Gold += monsterGold[monster.Name]
		}

		if monster.Name == hiveQueen && !monster.Flags[lurk.Alive] {
			user.killedQueen = true
		}
		if monster.Name == formicFleet && !monster.Flags[lurk.Alive] {
			user.killedFleet = true
		}

		g.startHealTimer(monster)
		if err := g.sendCharacters(currentRoom, conn); err != nil {
			return err
		}
	}

	if fights == 0 {
		return g.sendError(conn, cross.NoFight, fmt.Sprintf(
			"No live monsters to fight in the room %v", currentRoom.r.RoomName,
		))
	}

	if user.c.Flags[lurk.Alive] {
		return nil
	}
	log.Printf("%v died in a fight", user.c.Name)

	_, err := conn.Write(lurk.Marshal(&lurk.Message{
		Recipient: player,
		Sender:    narrator,
		Narration: true,
		Text:      "You have lost in battle. Regenerate your health to fight again.",
	}))

	return err
}

// Is only called in thread safe function.
func (g *game) handleHiveQueenFight(user *user, conn net.Conn) error {
	hq := g.monsters[hiveQueenCocoon]
	if user.c.RoomNum != shakespeare {
		return g.sendError(conn, cross.NoFight, fmt.Sprintf("user %s is not in the same room as you", hq.Name))
	}

	lurk.CalculateFight(user.c, hq)
	if !hq.Flags[lurk.Alive] {
		log.Printf("%s killed the hive queen cocoon\n", user.c.Name)
		// Add more gameplay here.
		if _, err := conn.Write(lurk.Marshal(&lurk.Message{
			Recipient: user.c.Name,
			Sender:    narrator,
			Narration: true,
			Text:      "You have committed true Xenocide.",
		})); err != nil {
			return err
		}
	}
	return g.sendCharacters(g.rooms[shakespeare], conn)
}

func (g *game) handlePVPFight(pvp *lurk.PVPFight, conn net.Conn, player string) (err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	user, ok := g.users[player]
	if !ok {
		return g.sendError(conn, cross.Other, fmt.Sprintf("%v: error in PVP fighting", cross.ErrUserNotInServer.Error()))
	}

	if pvp.TargetName == hiveQueenCocoon {
		return g.handleHiveQueenFight(user, conn)
	}

	target, ok := g.users[pvp.TargetName]
	if !ok {
		return g.sendError(conn, cross.Other, fmt.Sprintf("%v: error in PVP fighting", cross.ErrUserNotInServer.Error()))
	}

	if user.c.RoomNum != target.c.RoomNum {
		return g.sendError(conn, cross.NoFight, fmt.Sprintf("user %s is not in the same room as you", target.c.Name))
	}

	if !user.c.Flags[lurk.Alive] {
		return g.sendError(conn, cross.NoFight, player+", you cannot fight when you are dead")
	}
	if !target.c.Flags[lurk.Alive] {
		return g.sendError(conn, cross.NoFight, target.c.Name+" is already dead!")
	}

	if _, err = target.conn.Write(lurk.Marshal(&lurk.Message{
		Recipient: target.c.Name,
		Sender:    narrator,
		Text:      fmt.Sprintf("You have been engaged in combat by %v!", user.c.Name),
		Narration: true,
	})); err != nil {
		return err
	}

	lurk.CalculateFight(user.c, target.c)

	if err = g.sendCharacters(g.rooms[user.c.RoomNum], conn); err != nil {
		return err
	}
	if err = g.sendCharacters(g.rooms[target.c.RoomNum], target.conn); err != nil {
		return err
	}

	if !user.c.Flags[lurk.Alive] {
		log.Printf("%v died in a fight", user.c.Name)

		_, err = conn.Write(lurk.Marshal(&lurk.Message{
			Recipient: player,
			Sender:    narrator,
			Narration: true,
			Text:      "You have lost in battle. Regenerate your health to fight again.",
		}))
	}

	if !target.c.Flags[lurk.Alive] {
		log.Printf("%v died in a fight", target.c.Name)

		_, err = target.conn.Write(lurk.Marshal(&lurk.Message{
			Recipient: player,
			Sender:    narrator,
			Narration: true,
			Text:      "You have lost in battle. Regenerate your health to fight again.",
		}))
	}

	return err
}
func (g *game) handleLoot(loot *lurk.Loot, player string) {}
func (g *game) handleLeave(player string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	defer delete(g.users, player)

	user, ok := g.users[player]
	if !ok {
		return
	}

	oldRoom := user.c.RoomNum
	user.c.RoomNum = 0

	for _, other := range g.users {
		if other.c.RoomNum != oldRoom {
			continue
		}
		if err := g.sendCharacterUpdate(user.c, other.conn, other.c.Name, fmt.Sprintf("%s left the server!", player)); err != nil {
			log.Printf("%v: error when updating others of leaving the server", err.Error())
		}
	}
}
