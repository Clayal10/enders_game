package client

import (
	"fmt"
	"strings"

	"github.com/Clayal10/enders_game/pkg/lurk"
)

// Each field correlates to a section of the UI that should be updated in the JS.
// This struct will be json'd and each field should be added to the innerHTML of the HTML element,
// and should not overwrite the data already in it.
type ClientState struct {
	Info        string `json:"info"`
	Rooms       string `json:"rooms"`
	Connections string `json:"connections"`
	Players     string `json:"players"`
	Id          int64  `json:"id"`

	characters map[string]*lurk.Character // key is name

	room *lurk.Room
}

func newClientState(id int64) *ClientState {
	return &ClientState{
		Id:         id,
		characters: map[string]*lurk.Character{},
	}
}

// updateClientState will take a slice of LurkMessage interface objects and modifies the ClientUpdate field with
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
			c.State.characters[character.Name] = character
			if character.Name == c.character.Name {
				c.character = character
			}
			c.stringifyCharacters()
		case lurk.TypeRoom:
			room := msg.(*lurk.Room)
			c.State.room = room
			c.State.resetState()
			c.State.Rooms += fmt.Sprintf(roomTemplate, c.State.room.RoomNumber, c.State.room.RoomName, c.State.room.RoomDesc)
		case lurk.TypeConnection:
			connection := msg.(*lurk.Connection)

			newConnection := fmt.Sprintf(connectionTemplate, connection.RoomNumber, connection.RoomName, connection.RoomDesc)
			if !strings.Contains(c.State.Connections, newConnection) {
				c.State.Connections += newConnection
			}
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
	c.State.Players = ""
	c.State.Players = fmt.Sprintf(userTemplate, c.character.Name, c.character.Attack, c.character.Defense, c.character.Regen, c.character.Health, c.character.Gold)
	for _, character := range c.State.characters {
		if character.RoomNum != c.character.RoomNum || character.Name == c.character.Name {
			continue
		}

		switch {
		case !character.Flags[lurk.Alive]:
			c.State.Players += fmt.Sprintf(deadEntity, character.Name)
		case character.Flags[lurk.Monster]:
			c.State.Players += fmt.Sprintf(monsterTemplate, character.Name, character.Attack, character.Defense, character.Regen, character.Health)
		default:
			c.State.Players += fmt.Sprintf(characterTemplate, character.Name, character.Attack, character.Defense, character.Regen, character.Health, character.Gold)
		}
	}
}

func (state *ClientState) resetState() {
	state.Rooms = ""
	state.Connections = ""
	state.Players = ""
	state.characters = map[string]*lurk.Character{}
}

const characterTemplate = `
%s
  | Attack: %v
  | Defense: %v
  | Regen: %v
  | Health: %v
  | Gold: %v
  `

const monsterTemplate = `
<span style="color: red;">%s</span>
  | Attack: %v
  | Defense: %v
  | Regen: %v
  | Health: %v
  `

const userTemplate = `
<span style="color: green;">%s</span>
  | Attack: %v
  | Defense: %v
  | Regen: %v
  | Health: %v
  | Gold: %v
  `

const deadEntity = `
<span style="background-color: red; color: white;">%s</span>
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
