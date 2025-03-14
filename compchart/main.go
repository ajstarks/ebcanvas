// component charts
package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"

	"github.com/ajstarks/ebcanvas"
	"github.com/ajstarks/ebcanvas/chart"
	"github.com/hajimehoshi/ebiten/v2"
)

var screenWidth = 1000
var screenHeight = 1000

type App struct{}

func (g *App) Update() error {
	return nil
}

func (g *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *App) Draw(screen *ebiten.Image) {
	comp(screen)
}

func comp(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Width = screenWidth
	canvas.Height = screenHeight
	canvas.Screen = screen

	sr, err := os.Open("sine.d")
	if err != nil {
		return
	}
	cr, err := os.Open("cosine.d")
	if err != nil {
		return
	}
	sine, err := chart.DataRead(sr)
	if err != nil {
		return
	}
	cosine, err := chart.DataRead(cr)
	if err != nil {
		return
	}
	cosine.Zerobased = false
	sine.Zerobased = false
	canvas.Background(color.NRGBA{255, 255, 255, 255})
	cosine.Frame(canvas, 5)
	sine.Label(canvas, 1.5, 10, "", "gray")
	cosine.YAxis(canvas, 1.2, -1.0, 1.0, 0.5, "%0.2f", true)
	cosine.Color = color.NRGBA{0, 128, 0, 255}
	sine.Color = color.NRGBA{128, 0, 0, 255}
	cosine.Scatter(canvas, 0.5)
	sine.Scatter(canvas, 0.5)

	sine.Left = 10
	sine.Right = sine.Left + 40
	sine.Top, cosine.Top = 30, 30
	sine.Bottom, cosine.Bottom = 10, 10

	sine.CTitle(canvas, 2, 2)
	sine.Frame(canvas, 10)
	sine.Scatter(canvas, 0.25)

	offset := 45.0
	cosine.Left = sine.Left + offset
	cosine.Right = sine.Right + offset

	cosine.CTitle(canvas, 2, 2)
	cosine.Frame(canvas, 10)
	cosine.Scatter(canvas, 0.25)
}

func main() {

	flag.IntVar(&screenWidth, "width", 1000, "canvas width")
	flag.IntVar(&screenHeight, "height", 1000, "canvas height")
	flag.Parse()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("sine+cosine")

	if err := ebcanvas.LoadFont(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
