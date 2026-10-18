package game

import (
	"embed"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed resources/*
var resourceFS embed.FS

type elemAttr int

const (
	scaleX elemAttr = iota
	rotationX
	translationX
	rotationY
	scaleY
	translationY
)

type elem struct {
	i, j int
}

var elementMap = map[elemAttr]elem{
	scaleX:       {0, 0},
	rotationX:    {0, 1},
	translationX: {0, 2},
	rotationY:    {1, 0},
	scaleY:       {1, 1},
	translationY: {1, 2},
}

func getElement(e elemAttr) (int, int) {
	return elementMap[e].i, elementMap[e].j
}

const (
	gopher          = "resources/space_gopher.png"
	gopherSize      = 0.1
	gopherSpeed     = 4
	gopherSpeedSlow = 2
)

type Game struct {
	*resources
}

type resources struct {
	images  map[string]*ebiten.Image            // key is file name
	objects map[string]*ebiten.DrawImageOptions // key is same for images=
}

func (g *Game) Update() error {
	// The Gopher
	op := g.resources.objects[gopher]
	var speed float64 = gopherSpeed
	// Speed
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
		speed = gopherSpeedSlow
	}

	// Movement
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		op.GeoM.Translate(0, -speed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		op.GeoM.Translate(0, speed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		op.GeoM.Translate(-speed, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		op.GeoM.Translate(speed, 0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0xA, 0xC, 0xFF})

	// The Gopher
	op := g.resources.objects[gopher]

	i, j := getElement(scaleX)
	op.GeoM.SetElement(i, j, 0.1)
	i, j = getElement(scaleY)
	op.GeoM.SetElement(i, j, 0.1)
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
	objects := make(map[string]*ebiten.DrawImageOptions)

	// Make into some iterable loop if we have more images
	file, err := resourceFS.Open(gopher)
	if err != nil {
		return err
	}
	image, _, err := image.Decode(file)
	if err == nil {
		ebImage := ebiten.NewImageFromImage(image)
		images[gopher] = ebImage
		// allocate a new image option to be associated when drawing the image
		objects[gopher] = &ebiten.DrawImageOptions{}
	}

	g.resources = &resources{
		images:  images,
		objects: objects,
	}

	return err
}
