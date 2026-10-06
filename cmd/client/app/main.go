package main

import (
	_ "embed"
	"log"

	"github.com/Clayal10/enders_game/cmd/client/app/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &game.Game{}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Lurk Client")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
