package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	screenWidth  int = 1000
	screenHeight int = 1000
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	mondiran(screen)
}

func mondiran(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Height = screenHeight
	canvas.Width = screenWidth

	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	blue := color.RGBA{0, 0, 255, 255}
	red := color.RGBA{255, 0, 0, 255}
	yellow := color.RGBA{255, 255, 0, 255}

	var third float32 = 100.0 / 3
	var border float32 = 1
	halft := third / 2
	qt := third / 4
	t2 := third * 2
	tq := 100.0 - qt
	t2h := t2 + halft

	canvas.Background(white)
	canvas.CenterRect(halft, halft, third, third, blue)      // lower left blue square
	canvas.CenterRect(t2, t2, t2, t2, red)                   // big red
	canvas.CenterRect(tq, qt, halft, halft, yellow)          // small yellow lower right
	canvas.Line(0, 0, 100, 0, border, black)                 // top border
	canvas.Line(0, 0, 0, 100, border, black)                 // left border
	canvas.Line(100, border/2, 100, 100, border, black)      // right border
	canvas.Line(0, 100, 100, 100, border, black)             // bottom border
	canvas.Line(t2h, halft, t2h+halft, halft, border, black) // top of yellow square
	canvas.Line(third, 100, third, 0, border, black)         //  first column border
	canvas.Line(t2h, 0, t2h, third, border, black)           // left of small right squares
	canvas.Line(0, third, 100, third, border, black)         // top of bottom squares
	canvas.Line(0, t2, third, t2, border, black)             // border between left white squares
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mondrian")
	if err := ebiten.RunGame(&Game{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
