package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type App struct {
	page int
}

func (a *App) Update() error {
	_, wy := ebiten.Wheel()

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		a.page = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || wy > 0 {
		a.page++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || wy < 0 {
		a.page--
	}
	if a.page < 0 {
		a.page = 0
	}
	return nil
}

func (g *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (app *App) Draw(screen *ebiten.Image) {
	work(app, screen)
}

var screenWidth = 1000
var screenHeight = 1000

var palette = []color.NRGBA{
	color.NRGBA{128, 128, 128, 255},
	color.NRGBA{255, 0, 0, 255},
	color.NRGBA{0, 255, 0, 255},
	color.NRGBA{0, 0, 255, 255},
	color.NRGBA{173, 154, 255, 255},
	color.NRGBA{255, 0, 255, 255},
}

func work(app *App, screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight

	white := color.NRGBA{255, 255, 255, 255}
	black := color.NRGBA{0, 0, 0, 255}
	datacolor := color.NRGBA{128, 0, 0, 255}
	bgcolor := white
	txcolor := black
	dotcolor := datacolor

	var bm float32 = 20
	var lf float32 = 10
	//var angle float64 = ebcanvas.Pi / 4
	dx := []float32{lf, 30, 40, 50, 60, 70, lf}
	dy := []float32{30, 80, 40, 30, 25, bm, bm}

	canvas.Background(bgcolor)
	canvas.Text(20, 90, 4, "Testing...", txcolor)
	canvas.CText(50, 10, 3, fmt.Sprintf("Page: %d", app.page), txcolor)
	for i := 0; i < len(dx); i++ {
		canvas.Circle(dx[i], dy[i], 1.5, dotcolor)
		canvas.CText(dx[i], bm-3, 2, fmt.Sprintf("%.0f", dx[i]), txcolor)
	}
	paintColor := palette[app.page%len(palette)]
	if app.page%2 == 0 {
		canvas.Polygon(dx, dy, paintColor)
	} else {
		canvas.StrokedPolygon(dx, dy, 0.2, paintColor)
	}
	var tw float32 = 30.0
	var tx float32 = 40.0
	var ts float32 = 2.5
	tm := "I am the resurrection and the life. He that believeth in me though he were dead, yet shall he live. And whosoever liveth and believeth in me shall never die." // "For since by man came death, by man came also the resurrection. For as in Adam all die, in Christ shall all be made alive."
	canvas.TextWrap(tx, 95, tw, ts, tm, txcolor)
	canvas.TextWrapStrict(tx, 70, tw, ts, tm, txcolor)
	linecolor := txcolor
	linecolor.A = 100
	canvas.VLine(tx, 0, 100, 0.1, linecolor)
	canvas.VLine(tx+tw, 0, 100, 0.1, linecolor)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("hello")

	if err := ebcanvas.LoadFont(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
