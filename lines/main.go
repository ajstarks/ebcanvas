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
	lines(screen)
}

var screenWidth = 1000
var screenHeight = 1000

func lines(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight

	var x, y, lw, ls float32
	lw = 0.1
	ls = 1

	canvas.Background(color.NRGBA{255, 255, 255, 255})
	for y = 5; y <= 95; y += 5 {
		canvas.Line(50, 50, 95, y, lw, color.NRGBA{128, 0, 0, 128})
		canvas.Line(50, 50, 5, y, lw, color.NRGBA{0, 0, 128, 128})
		canvas.Coord(95, y, ls, "", color.NRGBA{0, 0, 0, 255})
		canvas.Coord(5, y, ls, "", color.NRGBA{0, 0, 0, 255})
		lw += 0.1
	}

	lw = 0.1
	for x = 5; x <= 95; x += 5 {
		canvas.Line(50, 50, x, 95, lw, color.NRGBA{0, 128, 0, 128})
		canvas.Line(50, 50, x, 5, lw, color.NRGBA{0, 0, 0, 128})
		canvas.Coord(x, 95, ls, "", color.NRGBA{0, 0, 0, 255})
		canvas.Coord(x, 5, ls, "", color.NRGBA{0, 0, 0, 255})
		lw += 0.1
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("hello")

	if err := ebcanvas.LoadFont(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	ebiten.SetTPS(0)
	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
