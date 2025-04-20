// test the chart package
package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"

	ec "github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
)

// NameValue defines data
type NameValue struct {
	name  string
	note  string
	value float64
}

// ChartOptions define all the components of a chart
type ChartOptions struct {
	showtitle, showscatter, showarea, showframe, showlegend, showbar bool
	title, legend, color                                             string
	xlabelInterval                                                   int
}

var chartopts ChartOptions
var screenWidth = 1000
var screenHeight = 1000
var data []NameValue

type App struct{}

func (g *App) Update() error {
	return nil
}

func (g *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *App) Draw(screen *ebiten.Image) {
	chart(screen)
}

func minmax(data []NameValue) (float64, float64) {
	min := data[0].value
	max := data[0].value
	for _, d := range data {
		if d.value > max {
			max = d.value
		}
		if d.value < min {
			min = d.value
		}
	}
	return min, max
}

// DataRead reads tab separated values into a NameValue slice
func DataRead(r io.Reader) ([]NameValue, error) {
	var d NameValue
	var data []NameValue
	var err error
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 { // skip blank lines
			continue
		}
		if t[0] == '#' && len(t) > 2 { // process titles
			// title = strings.TrimSpace(t[1:])
			continue
		}
		fields := strings.Split(t, "\t")
		if len(fields) < 2 {
			continue
		}
		if len(fields) == 3 {
			d.note = fields[2]
		} else {
			d.note = ""
		}
		d.name = fields[0]
		d.value, err = strconv.ParseFloat(fields[1], 64)
		if err != nil {
			d.value = 0
		}
		data = append(data, d)
	}
	err = scanner.Err()
	return data, err
}

func xaxis(canvas *ec.Canvas, x, y, width, height float32, interval int, data []NameValue) {
	for i, d := range data {
		xp := float32(ec.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		if interval > 0 && i%interval == 0 {
			canvas.TextMid(xp, y-3, 1.5, d.name, color.NRGBA{0, 0, 0, 255})
			canvas.Line(xp, y, xp, height, 0.1, color.NRGBA{0, 0, 0, 128})
		}
	}
	canvas.Line(x, height, width, height, 0.1, color.NRGBA{0, 0, 0, 128})
	canvas.Line(width, height, width, y, 0.1, color.NRGBA{0, 0, 0, 128})
}

func frame(canvas *ec.Canvas, x, y, width, height float32, color color.NRGBA) {
	canvas.CornerRect(x, height, width-x, height-y, color)
}

func dotchart(canvas *ec.Canvas, x, y, width, height float32, data []NameValue, datacolor color.NRGBA) {
	min, max := minmax(data)
	for i, d := range data {
		xp := float32(ec.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		yp := float32(ec.MapRange(d.value, min, max, float64(y), float64(height)))
		canvas.Circle(xp, yp, 0.3, datacolor)
	}
}

func barchart(canvas *ec.Canvas, x, y, width, height float32, data []NameValue, datacolor color.NRGBA) {
	min, max := minmax(data)
	for i, d := range data {
		xp := float32(ec.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		yp := float32(ec.MapRange(d.value, min, max, float64(y), float64(height)))
		canvas.VLine(xp, y, yp-y, 0.1, datacolor)
	}
}

func areachart(canvas *ec.Canvas, x, y, width, height float32, data []NameValue, datacolor color.NRGBA) {
	min, max := minmax(data)
	l := len(data)

	ax := make([]float32, l+2)
	ay := make([]float32, l+2)
	ax[0] = x
	ay[0] = y
	ax[l+1] = width
	ay[l+1] = y

	for i, d := range data {
		xp := float32(ec.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		yp := float32(ec.MapRange(d.value, min, max, float64(y), float64(height)))
		ax[i+1] = xp
		ay[i+1] = yp
	}
	datacolor.A = 128
	canvas.Polygon(ax, ay, datacolor)
}

func chart(screen *ebiten.Image) {
	black := color.NRGBA{0, 0, 0, 255}
	datacolor := ec.ColorLookup(chartopts.color)
	bgcolor := color.NRGBA{255, 255, 255, 255}
	framecolor := color.NRGBA{0, 0, 0, 20}
	canvas := new(ec.Canvas)
	canvas.Width = screenWidth
	canvas.Height = screenHeight
	canvas.Screen = screen

	canvas.Background(bgcolor)
	if chartopts.showtitle {
		canvas.Text(10, 90, 3, chartopts.title, black)
	}
	if chartopts.showlegend {
		canvas.Text(10, 84, 2.5, chartopts.legend, datacolor)
		canvas.HLine(20, 85, 2, 1, datacolor)
	}
	if chartopts.xlabelInterval > 0 {
		xaxis(canvas, 10, 15, 90, 70, chartopts.xlabelInterval, data)
	}
	if chartopts.showframe {
		frame(canvas, 10, 15, 90, 70, framecolor)
	}
	if chartopts.showscatter {
		dotchart(canvas, 10, 15, 90, 70, data, datacolor)
	}
	if chartopts.showarea {
		areachart(canvas, 10, 15, 90, 70, data, datacolor)
	}
	if chartopts.showbar {
		barchart(canvas, 10, 15, 90, 70, data, datacolor)
	}
}

func main() {

	flag.IntVar(&screenWidth, "width", 1200, "canvas width")
	flag.IntVar(&screenHeight, "height", 900, "canvas height")
	flag.IntVar(&chartopts.xlabelInterval, "xlabel", 0, "show x axis")

	flag.StringVar(&chartopts.title, "chartitle", "", "chart title")
	flag.StringVar(&chartopts.legend, "chartlegend", "", "chart legend")
	flag.StringVar(&chartopts.color, "color", "maroon", "chart data color")

	flag.BoolVar(&chartopts.showtitle, "title", true, "show title")
	flag.BoolVar(&chartopts.showlegend, "legend", false, "show legend")
	flag.BoolVar(&chartopts.showbar, "bar", false, "show bar chart")
	flag.BoolVar(&chartopts.showarea, "area", false, "show area chart")
	flag.BoolVar(&chartopts.showscatter, "scatter", false, "show scatter chart")
	flag.BoolVar(&chartopts.showframe, "frame", false, "show frame")
	flag.Parse()

	var err error
	data, err = DataRead(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("chartest")

	if err = ec.LoadFont(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	if err = ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

}
