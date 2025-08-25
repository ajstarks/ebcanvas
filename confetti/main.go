package main

import (
	"fmt"
	"image/color"
	"math/rand"
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
	rgb(screen)
}

var screenWidth = 1000
var screenHeight = 1000

func rn(n int) float32 {
	return float32(rand.Intn(n))
}

func rn8(n int) uint8 {
	return uint8(rand.Intn(n))
}

func rgb(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight

	nshapes := 100
	maxsize := 10
	canvas.Screen.Fill(color.NRGBA{0, 0, 0, 255})
	for i := 0; i < nshapes; i++ {
		color := color.NRGBA{rn8(255), rn8(255), rn8(255), rn8(255)}
		x, y := rn(100), rn(100)
		w, h := rn(maxsize), rn(maxsize)
		if i%2 == 0 {
			canvas.Circle(x, y, w, color)
		} else {
			canvas.CenterRect(x, y, w, h, color)
		}
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("confetti")
	ebiten.SetTPS(0)
	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
