package main

import (
	_ "embed"
	"log"

	"github.com/Clayal10/enders_game/cmd/client/app/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &game.Game{}
	if err := game.SetResources(); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(1080, 720)
	ebiten.SetWindowTitle("Lurk Client")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
