package client

import (
	"fmt"

	"github.com/Clayal10/enders_game/lib/lurk"
)

// updateClientState will take a slice of LurkMessage interface objects and return a ClientUpdate with
// the proper text fields.
func (c *Client) updateClientState(lurkMessages []lurk.LurkMessage) {
	for _, msg := range lurkMessages {
		switch msg.GetType() {
		case lurk.TypeGame:
			game := msg.(*lurk.Game)
			c.Game = game
			c.State.Info += fmt.Sprintf("Stat Limit: %v\nInitial Points: %v\n%s\n", game.StatLimit, game.InitialPoints, game.GameDesc)
		case lurk.TypeVersion:
			version := msg.(*lurk.Version)
			extensionBytes := 0
			for _, v := range version.Extensions {
				extensionBytes += len(v)
			}
			c.State.Info += fmt.Sprintf("LURK Version %v.%v | %v Bytes of Extensions\n", version.Major, version.Minor, extensionBytes)
		case lurk.TypeCharacter:
			character := msg.(*lurk.Character)
			c.State.characters = append(c.State.characters, character)
			c.State.uniqueCharacters[character.Name] = len(c.State.characters) - 1
			if character.Name == c.character.Name {
				c.character = character
			}
			c.stringifyCharacters()
		case lurk.TypeRoom:
			room := msg.(*lurk.Room)
			c.State.room = room
			c.State.stringifyRooms()
		case lurk.TypeConnection:
			connection := msg.(*lurk.Connection)
			c.State.connections = append(c.State.connections, connection)
			c.State.stringifyRooms()
		case lurk.TypeMessage:
			message := msg.(*lurk.Message)
			c.State.Info += lineBreak
			if message.Narration {
				c.State.Info += fmt.Sprintf(narratorTemplate, message.Sender, message.Recipient, message.Text)
				continue
			}
			c.State.Info += fmt.Sprintf(messageTemplate, message.Sender, message.Recipient, message.Text)
		case lurk.TypeError:
			e := msg.(*lurk.Error)
			c.State.Info += lineBreak
			c.State.Info += fmt.Sprintf(errorTemplate, e.ErrCode, e.ErrMessage)
		}
	}
}

func (c *Client) stringifyCharacters() {
	c.State.Players = fmt.Sprintf(userTemplate, c.character.Name, c.character.Attack, c.character.Defense, c.character.Health)
	namesInList := map[string]bool{}
	for i, character := range c.State.characters {
		if character.RoomNum != c.character.RoomNum || namesInList[character.Name] || c.State.uniqueCharacters[character.Name] != i || character.Name == c.character.Name {
			continue
		}
		namesInList[character.Name] = true
		if character.Flags[lurk.Monster] {
			c.State.Players += fmt.Sprintf(monsterTemplate, character.Name, character.Attack, character.Defense, character.Health)
			continue
		}
		c.State.Players += fmt.Sprintf(characterTemplate, character.Name, character.Attack, character.Defense, character.Health)
	}
}

func (state *ClientState) stringifyRooms() {
	state.Rooms = ""
	roomsInList := map[uint16]bool{}
	state.Rooms += fmt.Sprintf(roomTemplate, state.room.RoomNumber, state.room.RoomName, state.room.RoomDesc)
	roomsInList[state.room.RoomNumber] = true
	for _, connection := range state.connections {
		if roomsInList[connection.RoomNumber] {
			continue
		}
		state.Rooms += fmt.Sprintf(connectionTemplate, connection.RoomNumber, connection.RoomName, connection.RoomDesc)
	}
}

const characterTemplate = `
%s
  | Attack: %v
  | Defense: %v
  | Health: %v
  `

const monsterTemplate = `
<span style="color: red;">%s</span>
  | Attack: %v
  | Defense: %v
  | Health: %v
  `

const userTemplate = `
<span style="color: green;">%s</span>
  | Attack: %v
  | Defense: %v
  | Health: %v
  `
const errorTemplate = `
<span style="color: red;">Error #%d</span>: %s
`

const roomTemplate = `
(Current Room) %v: %s
-> %s
`

const connectionTemplate = `
%v: %s
-> %s
`

const messageTemplate = `
%s => %s: %s`

const narratorTemplate = `
<span style="color: purple;">%s</span> => %s: %s`

const lineBreak = `
==================================================
`
