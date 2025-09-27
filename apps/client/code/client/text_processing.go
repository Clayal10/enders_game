package client

import (
	"fmt"

	"github.com/Clayal10/enders_game/lib/lurk"
)

// getClientUpdate will take a slice of LurkMessage interface objects and return a ClientUpdate with
// the proper text fields.
func (c *Client) getClientUpdate(lurkMessages []lurk.LurkMessage) *ClientUpdate {
	cu := &ClientUpdate{
		Id: c.id,
	}
	for _, msg := range lurkMessages {
		switch msg.GetType() {
		case lurk.TypeGame:
			game := msg.(*lurk.Game)
			c.Game = game
			cu.Info += fmt.Sprintf("Stat Limit: %v\nInitial Points: %v\n%s\n", game.StatLimit, game.InitialPoints, game.GameDesc)
		case lurk.TypeVersion:
			version := msg.(*lurk.Version)
			extensionBytes := 0
			for _, v := range version.Extensions {
				extensionBytes += len(v)
			}
			cu.Info += fmt.Sprintf("LURK Version %v.%v | %v Bytes of Extensions\n", version.Major, version.Minor, extensionBytes)
		case lurk.TypeCharacter:
			character := msg.(*lurk.Character)
			cu.Players += fmt.Sprintf("%s | Attack: %v Defense: %v Health: %v Monster?: %v\n", character.Name, character.Attack, character.Defense, character.Health, character.Flags[lurk.Monster])
		}
	}
	return cu
}
