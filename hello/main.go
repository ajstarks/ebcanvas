package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type App struct{}

func (a *App) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	return nil
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (a *App) Draw(screen *ebiten.Image) {
	hello(screen)
}

func loadimage(filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	earth, _, err = image.Decode(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	return nil
}

var earth image.Image
var screenWidth = 1000
var screenHeight = 666

func hello(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight

	bgcolor := color.NRGBA{0, 0, 0, 255}
	txcolor := color.NRGBA{255, 255, 255, 255}

	canvas.Background(bgcolor)
	canvas.Image(50, 0, 100, earth)
	canvas.CText(50, 80, 10, "hello, world", txcolor)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("hello")

	if err := loadimage("earth.jpg"); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	if err := ebcanvas.LoadFont(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
