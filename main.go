package main

import (
	"embed"
	"image"
	"log"
	"os"

	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

//go:embed assets/*.jpeg
var fsys embed.FS

var (
	bg     *ebiten.Image
	person *ebiten.Image
	cat    *ebiten.Image
	window *ebiten.Image
)

var (
	fontFace *text.GoTextFace
	ticks    = 0
)

type game struct {
}

func (g *game) Update() error {
	ticks++
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ticks = 0
	}
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
		op.ColorScale.ScaleAlpha(0.5)
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
	// Dialogue
	textop := &text.DrawOptions{}
	textop.LineSpacing = 30 * 1.5
	textop.ColorScale.Scale(0, 0, 0, 1)
	glyphs := text.AppendGlyphs(nil, "Hello world! I am a cat. No name yet.", fontFace, &textop.LayoutOptions)
	length := ticks / 5 // Send 1 character every 5 ticks
	for i, g := range glyphs {
		if i > length {
			break // Ends when the number of displayed characters is exceeded
		}
		textop.GeoM.Reset() // Reuse textop with all characters and reset only GeoM
		textop.GeoM.Translate(200, 480)
		textop.GeoM.Translate(g.X, g.Y)
		screen.DrawImage(g.Image, &textop.DrawImageOptions)
	}
	// Name box
	textop = &text.DrawOptions{}
	w, h := text.Measure("Cat", fontFace, 30*1.5)
	textop.GeoM.Translate(250-w/2, 430-h/2)
	textop.ColorScale.Scale(0, 0, 0, 1)
	text.Draw(screen, "Cat", fontFace, textop)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	// Background
	var err error
	bg, err = loadEmbededAsset(fsys, "assets/bg.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	// Person
	person, err = loadEmbededAsset(fsys, "assets/person.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	// Cat
	cat, err = loadEmbededAsset(fsys, "assets/cat.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	// Window
	window, err = loadEmbededAsset(fsys, "assets/window.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	// Fonts
	f, err := os.Open("assets/NotoSans-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	src, err := text.NewGoTextFaceSource(f)
	if err != nil {
		log.Fatal(err)
	}
	fontFace = &text.GoTextFace{Source: src, Size: 30}

	ebiten.SetWindowTitle("Visual Novel Game")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	g := &game{}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func loadEmbededAsset(fsys embed.FS, path string) (*ebiten.Image, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	img = ebiten.NewImageFromImage(img)
	ebiImg, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return nil, err
	}
	return ebiImg, nil
}
