package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
)

type App struct{}

func (g *App) Update() error {
	return nil
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenWidth, screenHeight = ebcanvas.DisplayScale(outsideWidth, outsideHeight)
	return screenWidth, screenHeight
}

func (g *App) Draw(screen *ebiten.Image) {
	eclipse(screen)
}

var screenWidth = 1000
var screenHeight = 1000

func eclipse(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight
	black := color.NRGBA{0, 0, 0, 255}
	white := color.NRGBA{255, 255, 255, 255}
	var r float32 = 5.0
	var y float32 = 50.0
	var x float32 = 10.0
	for x = 10.0; x < 100.0; x += 15 {
		canvas.Circle(x, 50, r+0.5, white)
		canvas.Circle(x, y, r, black)
		y -= 2
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("eclipse")
	ebiten.SetTPS(0)
	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
