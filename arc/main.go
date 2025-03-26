package main

import (
	"fmt"
	"image"
	"image/color"
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
	work(screen)
}

var screenWidth = 1000
var screenHeight = 1000
var douglass image.Image

func work(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight

	var (
		cx        float32
		cy        float32
		r         float32
		deg       float32
		ts        float32 = 1.5
		fillcolor         = color.NRGBA{0, 0, 128, 255}
		dotcolor          = color.NRGBA{128, 0, 0, 255}
		txcolor           = fillcolor
		bgcolor           = color.NRGBA{255, 255, 255, 255}
	)

	canvas.Background(bgcolor)
	r = 3.3
	cx = 2.5
	cy = 50.0
	wcolor := fillcolor
	wcolor.A = 100

	canvas.CText(50, 90, 3, "Arcs, Wedges, Rotated Text", txcolor)
	for deg = 0; deg < 360; deg += 30 {
		canvas.CText(cx, cy-r*2, ts, fmt.Sprintf("%.0f°", deg), txcolor)
		canvas.Circle(cx, cy, r/10, dotcolor)
		canvas.StrokedArc(cx, cy, r, 0, deg, 0.25, fillcolor)
		canvas.Wedge(cx, cy, r, 0, deg, wcolor)
		cx += r * 2.5
	}
	canvas.StrokedArc(50, 75, 10, 0, 90, 1, dotcolor)
	canvas.StrokedArc(50, 75, 10, 90, 180, 1, fillcolor)
	canvas.StrokedArc(50, 75, 10, 180, 270, 1, dotcolor)
	canvas.StrokedArc(50, 75, 10, 270, 360, 1, fillcolor)

	canvas.Wedge(50, 75, 10, 0, 90, dotcolor)
	canvas.Wedge(50, 75, 10, 180, 270, txcolor)

	for deg = 0; deg < 360; deg += 30 {
		canvas.RText(50, 20, deg, 2.5, fmt.Sprintf("rotated: %.0f°", deg), txcolor)
	}

}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetTPS(0)
	ebiten.SetWindowTitle("arcs")

	if err := ebcanvas.LoadFont(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
