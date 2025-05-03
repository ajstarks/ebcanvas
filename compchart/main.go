// component charts
package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/ajstarks/ebcanvas"
	"github.com/ajstarks/ebcanvas/chart"
	"github.com/hajimehoshi/ebiten/v2"
)

var screenWidth = 1000
var screenHeight = 1000
var sine1 chart.ChartBox
var sine2 chart.ChartBox

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

func datagen() {
	var s1buf, s2buf bytes.Buffer
	begin := 0.0
	end := math.Pi * 4
	incr := 0.1
	// write to data sets, y=sin(x) and y=2*sin(x)
	for x := begin; x <= end; x += incr {
		if x == begin { // add data set labels
			fmt.Fprintln(&s1buf, "# y=sin(x)")
			fmt.Fprintln(&s2buf, "# y=2*sin(x)")
		}
		fmt.Fprintf(&s1buf, fmt.Sprintf("%.2f\t%f\n", x, math.Sin(x)))
		fmt.Fprintf(&s2buf, fmt.Sprintf("%.2f\t%f\n", x, 2*math.Sin(x)))
	}

	// read in data sets
	sr1 := bytes.NewReader(s1buf.Bytes())
	sr2 := bytes.NewReader(s2buf.Bytes())
	var s1err, s2err error
	sine1, s1err = chart.DataRead(sr1)
	if s1err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", s1err)
		os.Exit(1)
	}
	sine2, s2err = chart.DataRead(sr2)
	if s2err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", s2err)
		os.Exit(2)
	}
}

func comp(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Width = screenWidth
	canvas.Height = screenHeight
	canvas.Screen = screen

	datagen()

	minv, maxv := -2.0, 2.0
	dotsize := 0.4
	frameOpacity := 5.0

	// set attributes for each data set
	sine1.Zerobased = false
	sine2.Zerobased = false

	sine1.MinMax(minv, maxv)
	sine2.MinMax(minv, maxv)

	canvas.Background(color.NRGBA{255, 255, 255, 255})
	sine1.Frame(canvas, frameOpacity)
	sine1.Label(canvas, 1.5, 10, "", "gray")
	sine1.YAxis(canvas, 1.5, minv, maxv, 0.5, "%0.2f", true)

	sine1.Color = color.NRGBA{128, 0, 0, 255}
	sine2.Color = color.NRGBA{0, 128, 0, 255}

	// chart the data on the same frame
	sine1.Scatter(canvas, dotsize)
	sine2.Scatter(canvas, dotsize)

	// using the same data sets, make separate charts
	sine1.Left = 10
	sine1.Right = sine1.Left + 30

	sine1.Top = 30
	sine2.Top = 30

	sine1.Bottom = 10
	sine2.Bottom = 10

	dotsize /= 2
	frameOpacity *= 2

	sine1.CTitle(canvas, 2, 2)
	sine1.Frame(canvas, frameOpacity)
	sine1.Scatter(canvas, dotsize)

	offset := 50.0
	sine2.Left = sine1.Left + offset
	sine2.Right = sine1.Right + offset

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
