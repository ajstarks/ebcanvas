package main

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
)

type App struct{}

func (g *App) Update() error {
	return nil
}

func (g *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *App) Draw(screen *ebiten.Image) {
	polar(screen)
}

var screenWidth = 1000
var screenHeight = 1000

var (
	topcolor = color.RGBA{255, 0, 0, 100}
	botcolor = color.RGBA{0, 0, 255, 100}
	bgcolor  = color.RGBA{0, 0, 0, 255}
)

func polar(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight

	canvas.Background(bgcolor)
	var theta, radius float32
	for radius = 2; radius < 50; radius += 2 {
		for theta = 180; theta <= 360; theta += 15 { // degrees
			x, y := canvas.PolarDegrees(50, 50, radius, theta)
			canvas.Circle(x, y, radius/12, topcolor)
		}
		for theta = math.Pi / 16; theta < math.Pi; theta += math.Pi / 16 { // radians
			x, y := canvas.Polar(50, 50, radius, theta)
			canvas.Circle(x, y, radius/12, botcolor)
		}
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("polar")

	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
