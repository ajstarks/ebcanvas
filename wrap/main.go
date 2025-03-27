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
	tx float32
	tw float32
	ts float32
}

func (a *App) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		a.ts = 2.5
		a.tx = 30
		a.tw = 40
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		a.ts += 0.5
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		a.ts -= 0.5
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		a.tx += 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		a.tx -= 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		a.tw += 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		a.tw -= 1
	}
	if a.tx == 0 {
		a.ts = 2.5
		a.tx = 30
		a.tw = 40
	}
	if a.ts < 0 {
		a.ts = 2
	}
	if a.ts > 10 {
		a.ts = 2
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
var black = color.NRGBA{0, 0, 0, 255}
var white = color.NRGBA{255, 255, 255, 255}
var tm string = `If there is no struggle there is no progress.
Those who profess to favor freedom and yet deprecate agitation
are men who want crops without plowing up the ground;
they want rain without thunder and lightning.
They want the ocean without the awful roar of its many waters.`

func work(a *App, screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight
	bgcolor := black
	txcolor := white
	canvas.Background(bgcolor)
	wrap(canvas, false, a.tx, 95, a.tw, a.ts, tm, txcolor)
	wrap(canvas, true, a.tx, 45, a.tw, a.ts, tm, txcolor)
	canvas.EText(95, 2, 2, fmt.Sprintf("tx=%v tw=%v ts=%v", a.tx, a.tw, a.ts), txcolor)
}

func wrap(canvas *ebcanvas.Canvas, strict bool, x, y, w, size float32, s string, color color.NRGBA) {
	rcolor := color
	rcolor.A = 50
	canvas.CornerRect(x, 100, w, 100, rcolor)
	if strict {
		canvas.TextWrapStrict(x, y, w, size, s, color)
	} else {
		canvas.TextWrap(x, y, w, size, s, color)
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("wrap")

	if err := ebcanvas.LoadFont(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
