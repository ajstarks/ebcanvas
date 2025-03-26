package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Point struct {
	X, Y float32
}

type App struct{}

func (a *App) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	return nil
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (a *App) Draw(screen *ebiten.Image) {
	play(screen)
}

const pi = 3.14159265359

var (
	screenWidth  int
	screenHeight int
	tcolor       = color.NRGBA{128, 0, 0, 50}
	fcolor       = color.NRGBA{0, 0, 128, 50}
	stcolor      = color.NRGBA{128, 0, 0, 255}
	sfcolor      = color.NRGBA{0, 0, 128, 255}
	bgcolor      = color.NRGBA{255, 255, 255, 255}
	gridcolor    = color.NRGBA{128, 128, 128, 50}
	labelcolor   = color.NRGBA{50, 50, 50, 255}
	earth        image.Image
	drawgrid     bool
)

func loadAssets() error {
	r, err := os.Open("earth.jpg")
	if err != nil {
		return err
	}
	earth, _, err = image.Decode(r)
	if err != nil {
		return err
	}
	return nil
}

func play(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight

	var colx float32
	var lw float32 = 0.2
	var labelsize float32 = 2
	var wrap float32 = 15.0
	titlesize := labelsize * 2
	subsize := labelsize * 0.7
	wrapfmt := "This is text wrapped at a specified width (%1.f%% of the canvas width)."
	apimsg := "A canvas API for ebiten, using high-level objects, and a percentage-based coordinate system. (https://github.com/ebcanvas)"

	// Title
	canvas.Background(bgcolor)

	colx = 20
	canvas.Text(10, 92, titlesize, "Ebiten Canvas API", labelcolor)
	canvas.TextWrap(50, 96, 25, titlesize*.3, apimsg, labelcolor)

	// Lines
	canvas.CText(colx, 80, labelsize, "Line", labelcolor)
	canvas.Line(10, 70, colx+5, 65, lw, stcolor)
	canvas.Coord(10, 70, subsize, "P0", labelcolor)
	canvas.Coord(colx+5, 65, subsize, "P1", labelcolor)

	canvas.Line(colx, 70, 35, 75, lw, sfcolor)
	canvas.Coord(colx, 70, subsize, "P0", labelcolor)
	canvas.Coord(35, 75, subsize, "P1", labelcolor)

	// Circle
	cx1 := colx - 10
	cx2 := colx + 10
	canvas.CText(cx1, 55, labelsize, "Circle", labelcolor)
	canvas.Circle(cx1, 45, 5, fcolor)
	canvas.Coord(cx1, 45, subsize, "center", labelcolor)

	// Arc
	canvas.CText(cx2, 55, labelsize, "Arc", labelcolor)
	canvas.Arc(cx2, 45, 5, 0, 180, tcolor)
	canvas.StrokedArc(cx2, 45, 5, 0, 180, lw, stcolor)
	canvas.Coord(cx2, 45, subsize, "center", labelcolor)

	// Text
	tx := cx1 + (cx2-cx1)/2
	canvas.CText(colx, 30, labelsize, "Text", labelcolor)
	canvas.Text(tx, 25, subsize, "Begin-aligned", labelcolor)
	canvas.Circle(tx, 25, subsize/4, labelcolor)
	canvas.CText(tx, 20, subsize, "Centered", labelcolor)
	canvas.Circle(tx, 20, subsize/4, labelcolor)
	canvas.EText(tx, 15, subsize, "End-aligned", labelcolor)
	canvas.Circle(tx, 15, subsize/4, labelcolor)
	canvas.TextWrap(tx, 10, wrap, subsize, fmt.Sprintf(wrapfmt, wrap), labelcolor)
	canvas.Circle(tx, 10, subsize/4, labelcolor)
	canvas.RText(cx1, 5, 45, subsize, "Rotated", labelcolor)
	canvas.Circle(cx1, 5, subsize/4, labelcolor)
	// Quadradic Bezier
	start := Point{X: 45, Y: 65}
	c1 := Point{X: 70, Y: 85}
	end := Point{X: 70, Y: 65}
	canvas.CText(60, 80, labelsize, "Quadratic Bezier Curve", labelcolor)
	canvas.StrokedCurve(start.X, start.Y, c1.X, c1.Y, end.X, end.Y, lw, stcolor)
	canvas.Curve(start.X, start.Y, c1.X, c1.Y, end.X, end.Y, tcolor)
	canvas.Coord(start.X, start.Y, subsize, "start", labelcolor)
	canvas.Coord(c1.X, c1.Y, subsize, "control", labelcolor)
	canvas.Coord(end.X, end.Y, subsize, "end", labelcolor)

	colx += 40
	// Cubic Bezier
	start = Point{X: 45, Y: 40}
	c1 = Point{X: 45, Y: 55}
	c2 := Point{X: colx, Y: 50}
	end = Point{X: 70, Y: 40}
	canvas.CText(colx, 55, labelsize, "Cubic Bezier Curve", labelcolor)
	canvas.StrokedCubeCurve(start.X, start.Y, c1.X, c1.Y, c2.X, c2.Y, end.X, end.Y, lw, sfcolor)
	canvas.CubeCurve(start.X, start.Y, c1.X, c1.Y, c2.X, c2.Y, end.X, end.Y, fcolor)
	canvas.Coord(start.X, start.Y, subsize, "start", labelcolor)
	canvas.Coord(end.X, end.Y, subsize, "end", labelcolor)
	canvas.Coord(c1.X, c1.Y, subsize, "control 1", labelcolor)
	canvas.Coord(c2.X, c2.Y, subsize, "control 2", labelcolor)

	// Polygon
	canvas.CText(colx, 30, labelsize, "Polygon", labelcolor)
	xp := []float32{45, 60, 70, 70, 60, 45}
	yp := []float32{25, 20, 25, 5, 10, 5}
	for i := 0; i < len(xp); i++ {
		canvas.Coord(xp[i], yp[i], subsize, fmt.Sprintf("P%d", i), labelcolor)
	}
	canvas.StrokedPolygon(xp, yp, lw, stcolor)
	canvas.Polygon(xp, yp, tcolor)

	colx += 30
	// Rectangles
	canvas.CText(colx, 80, labelsize, "Rectangle", labelcolor)
	canvas.CenterRect(colx, 70, 5, 15, fcolor)
	canvas.Coord(colx, 70, subsize, "center", labelcolor)

	// Square
	canvas.CText(colx, 55, labelsize, "Square", labelcolor)
	canvas.Square(colx, 45, 10, tcolor)
	canvas.Coord(colx, 45, subsize, "center", labelcolor)

	// Image
	canvas.CText(colx, 30, labelsize, "Image", labelcolor)
	canvas.Image(colx, 15, (float32(screenWidth)/1000)*10, earth)
	canvas.Coord(colx, 15, subsize, "", color.NRGBA{255, 255, 255, 255})

	if drawgrid {
		canvas.Grid(0, 0, 100, 100, 0.1, 5, gridcolor)
	}
}
func main() {

	flag.BoolVar(&drawgrid, "grid", false, "draw a grid")
	flag.IntVar(&screenWidth, "width", 1600, "canvas width")
	flag.IntVar(&screenHeight, "height", 1000, "canvas height")
	flag.Parse()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("play")

	var err error
	earth, err = ebcanvas.LoadImage("earth.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	if err = ebcanvas.LoadFont(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	if err = ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
