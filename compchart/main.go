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
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	sr2, err := os.Open("sine2.d")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	sine, err := chart.DataRead(sr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	}
	sine2, err := chart.DataRead(sr2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(4)
	}

	canvas.Background(color.NRGBA{255, 255, 255, 255})
	minv, maxv := -2.0, 2.0
	dotsize := 0.4
	frameOpacity := 5.0
	sine2.Zerobased, sine.Zerobased = false, false
	sine.Minvalue, sine.Maxvalue = minv, maxv
	sine2.Minvalue, sine2.Maxvalue = minv, maxv

	sine.Frame(canvas, frameOpacity)
	sine.Label(canvas, 1.5, 10, "", "gray")
	sine.YAxis(canvas, 1.5, minv, maxv, 0.5, "%0.2f", true)

	sine2.Color = color.NRGBA{0, 128, 0, 255}
	sine.Color = color.NRGBA{128, 0, 0, 255}
	sine2.Scatter(canvas, dotsize)
	sine.Scatter(canvas, dotsize)

	sine.Left = 10
	sine.Right = sine.Left + 35
	sine.Top, sine2.Top = 30, 30
	sine.Bottom, sine2.Bottom = 10, 10
	dotsize /= 2
	frameOpacity *= 2

	sine.CTitle(canvas, 2, 2)
	sine.Frame(canvas, frameOpacity)
	sine.Scatter(canvas, dotsize)

	offset := 45.0
	sine2.Left = sine.Left + offset
	sine2.Right = sine.Right + offset

	sine2.CTitle(canvas, 2, 2)
	sine2.Frame(canvas, frameOpacity)
	sine2.Scatter(canvas, dotsize)
}

func main() {

	flag.IntVar(&screenWidth, "width", 1000, "canvas width")
	flag.IntVar(&screenHeight, "height", 1000, "canvas height")
	flag.Parse()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("composite chart")

	if err := ebcanvas.LoadFont(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
