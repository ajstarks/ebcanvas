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
	rgb(screen)
}

var screenWidth = 1000
var screenHeight = 1000

func rgb(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight
	colortab := []string{
		"orange",
		"rgb(100)",
		"rgb(100,50)",
		"rgb(100,50,2)",
		"rgb(100,50,2,100)",
		"hsv(0,70,50)",
		"hsv(0,70,50,50)",
		"#aa",
		"#aabb",
		"#aabbcc",
		"#aabbcc64",
		"rgb()",
		"hsv()",
		"#",
		"#error",
		"nonsense",
	}
	var x, y float32
	x, y = 50, 95
	canvas.Background(color.NRGBA{255, 255, 255, 255})
	for _, c := range colortab {
		canvas.EText(x-10, y, 3, c, color.NRGBA{0, 0, 0, 255})
		canvas.Circle(x, y+1, 2, ebcanvas.ColorLookup(c))
		y -= 6
	}

}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("rgb")

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
