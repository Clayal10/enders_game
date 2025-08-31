package server

import "github.com/Clayal10/enders_game/lib/lurk"

const (
	gameDescription = ` 
 ____  __ _  ____  ____  ____  _ ____     ___   __   _  _  ____ 
(  __)(  ( \(    \(  __)(  _ \(// ___)   / __) / _\ ( \/ )(  __)
 ) _) /    / ) D ( ) _)  )   /  \___ \  ( (_ \/    \/ \/ \ ) _) 
(____)\_)__)(____/(____)(__\_)  (____/   \___/\_/\_/\_)(_/(____)

The world has been ravaged by the most feared and despised being known to man, the formic. When it comes down to preventing their second massacre, will you be the one to step up and destroy them?`
)

func createRooms() map[uint16]*lurk.Room {
	return map[uint16]*lurk.Room{
		1: {
			Type:       lurk.TypeRoom,
			RoomNumber: 1,
			RoomName:   "Test Room",
			RoomDesc:   "A room with nothing in it?? Could it be the vast emptiness of space?? Maybe it's just this server's lack of real assets!",
		},
	}
}
