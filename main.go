package main

import (
	"image"
	"log"

	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1280
	screenHeight = 720
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
	// Calculate the scaling factors for width and height
	scaleX := float64(screenWidth) / float64(bg.Bounds().Dx())
	scaleY := float64(screenHeight) / float64(bg.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	screen.DrawImage(bg, op) // It is a little bit smaller than the window size
	screen.DrawImage(cat, nil)
	// Person
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(1280, 0)
	op.ColorScale.Scale(0.5, 0.5, 0.5, 0.5)
	screen.DrawImage(person, op)
	// Window
	// Source rectangle, 9 equal parts
	srcRects := []image.Rectangle{}
	xs := []int{0, 30, 70, 100}
	ys := []int{0, 30, 70, 100}
	for x := range 3 {
		for y := range 3 {
			rect := image.Rect(xs[x], ys[y], xs[x+1], ys[y+1])
			srcRects = append(srcRects, rect)
		}
	}
	// Destination rectangle, 9 parts but stretched.
	// corners remain the same 30x30
	dstRects := []image.Rectangle{}
	xs = []int{0, 30, 1080 - 30, 1080}
	ys = []int{0, 30, 260 - 30, 260}
	for x := range 3 {
		for y := range 3 {
			rect := image.Rect(xs[x], ys[y], xs[x+1], ys[y+1])
			dstRects = append(dstRects, rect)
		}
	}
	for i := range 9 {
		srcRect := srcRects[i]
		dstRect := dstRects[i]
		// The method SubImage takes an argument image.Rectangle and returns an
		// image cut out from that area.
		subImage := window.SubImage(srcRect).(*ebiten.Image)
		scaleX := float64(dstRect.Dx()) / float64(srcRect.Dx())
		scaleY := float64(dstRect.Dy()) / float64(srcRect.Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(100, 450) // New "origin" for the image
		op.GeoM.Translate(float64(dstRect.Min.X), float64(dstRect.Min.Y))
		screen.DrawImage(subImage, op)
	}
	// Small image for name -- TODO refactor code (Create helper function)
	dstRects = []image.Rectangle{}
	xs = []int{0, 30, 200 - 30, 200}
	ys = []int{0, 30, 70 - 30, 70}
	for x := range 3 {
		for y := range 3 {
			rect := image.Rect(xs[x], ys[y], xs[x+1], ys[y+1])
			dstRects = append(dstRects, rect)
		}
	}
	for i := range 9 {
		srcRect := srcRects[i]
		dstRect := dstRects[i]
		subImage := window.SubImage(srcRect).(*ebiten.Image)
		scaleX := float64(dstRect.Dx()) / float64(srcRect.Dx())
		scaleY := float64(dstRect.Dy()) / float64(srcRect.Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(150, 400)
		op.GeoM.Translate(float64(dstRect.Min.X), float64(dstRect.Min.Y))
		screen.DrawImage(subImage, op)
	}
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
	ebiten.SetWindowSize(screenWidth, screenHeight)
	g := &game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
