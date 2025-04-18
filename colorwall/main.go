// colorwall: inspired by Ellsworth Kelly's "Colors for a Large Wall, 1951'
package main

import (
	"fmt"
	"os"

	"github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	layout = [][]string{
		{"#000000", "#eeeeee", "#735976", "#eeeeee", "#000000", "#af5d23", "#eeeeee", "#366e93"}, // row 1
		{"#eeeeee", "#03342f", "#000000", "#eeeeee", "#ccb04d", "#eeeeee", "#a74e4a", "#000000"}, // row 2
		{"#000000", "#eeeeee", "#eeeeee", "#391a32", "#eeeeee", "#eeeeee", "#eeeeee", "#af5d23"}, // row 3
		{"#8a1f1b", "#eeeeee", "#366e93", "#eeeeee", "#5e825e", "#000000", "#391a32", "#eeeeee"}, // row 4
		{"#eeeeee", "#391a32", "#000000", "#eeeeee", "#eeeeee", "#8a1f1b", "#eeeeee", "#122e63"}, // row 5
		{"#03342f", "#eeeeee", "#eeeeee", "#366e93", "#eeeeee", "#eeeeee", "#03342f", "#000000"}, // row 6
		{"#eeeeee", "#a74e4a", "#5e825e", "#eeeeee", "#000000", "#735976", "#eeeeee", "#eeeeee"}, // row 7
		{"#000000", "#eeeeee", "#391a32", "#ccb04d", "#eeeeee", "#000000", "#a74e4a", "#000000"}, // row 8
	}

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
	colorwall(screen)
}

func colorwall(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Height = screenHeight
	canvas.Width = screenWidth

	var x, y, left, right, top, bottom, xincr, yincr, nc, nr float32
	left, right, bottom, top = 25, 85, 20, 80
	nr, nc = 8, 8

	xincr = (right - left) / nr
	yincr = (top - bottom) / nc
	bgcolor := ebcanvas.ColorLookup("#dddddd")
	basecolor := ebcanvas.ColorLookup("#bbbbbb")

	canvas.Background(bgcolor)
	canvas.Square(51.25, 53.75, 60.5, basecolor)
	y = top
	for i := 0; i < len(layout); i++ {
		row := layout[i]
		x = left
		for j := 0; j < len(row); j++ {
			canvas.Square(x, y, yincr-0.1, ebcanvas.ColorLookup(row[j]))
			x += xincr
		}
		y -= yincr
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("colorwall: inspired by Ellsworth Kelly's Colors for a Large Wall, 1951")
	if err := ebiten.RunGame(&Game{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
