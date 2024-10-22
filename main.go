package main

import (
	"log"

	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	bg     *ebiten.Image
	person *ebiten.Image
	cat    *ebiten.Image
	window *ebiten.Image
)

type game struct {
}

func (g *game) Update() error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.DrawImage(bg, nil)
	screen.DrawImage(person, nil)
	// Cat
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(560, 0)
	screen.DrawImage(cat, op)
	// Window
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(10.8, 2.7)
	op.GeoM.Translate(100, 450)
	screen.DrawImage(window, op)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	img, _, err := ebitenutil.NewImageFromFile("./assets/bg.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	bg = img
	img, _, err = ebitenutil.NewImageFromFile("assets/person.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	person = img
	img, _, err = ebitenutil.NewImageFromFile("assets/cat.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	cat = img
	img, _, err = ebitenutil.NewImageFromFile("assets/window.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	window = img
	ebiten.SetWindowTitle("Visual Novel Game")
	ebiten.SetWindowSize(1280, 720)
	g := &game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
