package game

import (
	"embed"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed resources/*
var resourceDirectory embed.FS

type Game struct {
	*resources
}

type resources struct {
	images *ebiten.DrawImageOptions
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Write your game's logical update.
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) getResources() *resources {
	if g.resources != nil {
		return g.resources
	}
	resourceDirectory.ReadFile("space_gopher.png")
}
