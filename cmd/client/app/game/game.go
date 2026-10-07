package game

import (
	"embed"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed resources/*
var resourceFS embed.FS

const gopher = "resources/space_gopher.png"

type Game struct {
	*resources
}

type resources struct {
	images map[string]*ebiten.Image // key is file name
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
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.resources.images[gopher], op)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) SetResources() error {
	if g.resources != nil {
		return nil
	}

	images := make(map[string]*ebiten.Image)

	// Make into some iterable loop if we have more images
	file, err := resourceFS.Open(gopher)
	if err != nil {
		return err
	}
	image, _, err := image.Decode(file)
	if err == nil {
		ebImage := ebiten.NewImageFromImage(image)
		images[gopher] = ebImage
	}

	g.resources = &resources{images: images}

	return err
}
