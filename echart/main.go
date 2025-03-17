// echart -- command line chart using ebcanvas chart package
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/ebcanvas"
	"github.com/ajstarks/ebcanvas/chart"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var opts chartOptions
var screenWidth, screenHeight int
var data chart.ChartBox

type App struct{}

func (g *App) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	return nil
}

func (g *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

type chartOptions struct {
	top, bottom, left, right                                                               float64
	barwidth, linewidth, linespacing, dotsize, textsize, piesize, ty, frameOp, opacity     float64
	bgcolor, dcolor, labelcolor, valuecolor, chartitle, yaxfmt, yrange, fontname, valuefmt string
	xlabel                                                                                 int
	zb, line, bar, hbar, scatter, area, pie, lego, dot, wbar, showtitle, showgrid          bool
}

// perr prints a filename and message to stderr
func perr(msg, file string) {
	io.WriteString(os.Stderr, file+": "+msg+"\n")
}

// string to floating point
func stof(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// yr parses the the yrange (max,min,step) string
func yr(yrange string, dmin, dmax float64) (float64, float64, float64) {
	var min, max, step float64
	min = dmin
	max = dmax
	step = max - min/5
	v := strings.Split(yrange, ",")
	switch len(v) {
	case 1:
		min = stof(v[0])
	case 2:
		min = stof(v[0])
		max = stof(v[1])
	case 3:
		min = stof(v[0])
		max = stof(v[1])
		step = stof(v[2])
	}
	return min, max, step
}

func (g *App) Draw(screen *ebiten.Image) {
	gchart(screen)
}

// gchart draws a chart
func gchart(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight

	// Define the colors
	datacolor := ebcanvas.ColorLookup(opts.dcolor)
	labelcolor := ebcanvas.ColorLookup(opts.labelcolor)
	bgcolor := ebcanvas.ColorLookup(opts.bgcolor)

	// Set the chart attributes
	data.Zerobased = opts.zb
	data.Top, data.Bottom = opts.top, opts.bottom
	data.Left, data.Right = opts.left, opts.right

	canvas.Background(bgcolor)

	// Draw the data
	data.Color = datacolor
	if opts.frameOp > 0 {
		data.Frame(canvas, opts.frameOp)
	}
	if opts.line {
		data.Line(canvas, opts.linewidth)
	}
	if opts.bar {
		data.Bar(canvas, opts.barwidth)
	}
	if opts.scatter {
		data.Scatter(canvas, opts.dotsize)
	}
	if opts.hbar {
		data.HBar(canvas, opts.barwidth, opts.linespacing, opts.textsize, opts.valuefmt, opts.valuecolor)
	}
	if opts.wbar {
		data.WBar(canvas, opts.linespacing, opts.textsize, opts.opacity, opts.valuefmt, opts.valuecolor)
	}
	if opts.dot {
		data.Dot(canvas, opts.dotsize)
	}
	if opts.area {
		data.Area(canvas, opts.opacity)
	}
	if opts.pie {
		data.Pie(canvas, opts.piesize)
	}
	if opts.lego {
		data.Lego(canvas, opts.dotsize)
	}

	// Draw labels, axes if specified
	data.Color = labelcolor
	if opts.line || opts.bar || opts.scatter || opts.area || opts.dot {
		data.Label(canvas, opts.textsize, opts.xlabel, opts.valuefmt, opts.valuecolor)
		if len(opts.yrange) > 0 {
			yaxmin, yaxmax, yaxstep := yr(opts.yrange, data.Minvalue, data.Maxvalue)
			data.YAxis(canvas, opts.textsize, yaxmin, yaxmax, yaxstep, opts.yaxfmt, opts.showgrid)
		}
	}

	if opts.hbar || opts.wbar {
		if len(opts.yrange) > 0 {
			yaxmin, yaxmax, yaxstep := yr(opts.yrange, data.Minvalue, data.Maxvalue)
			data.XAxis(canvas, opts.textsize, yaxmin, yaxmax, yaxstep, opts.yaxfmt, opts.showgrid)
		}
	}

	// Draw the chart titles
	if len(opts.chartitle) > 0 {
		data.Title = opts.chartitle
	}
	if opts.showtitle && len(data.Title) > 0 {
		data.CTitle(canvas, opts.textsize*2, opts.ty)
	}

}

func cmdUsage() {
	usage := `
gchart [options] file...

Options     Default               Description
.....................................................................
-area        false                make an area chart
-bar         false                make a bar chart
-dot         false                make a dot chart
-hbar        false                make a horizontal bar chart
-wbar        false                make a horizontal word bar chart
-lego        false                make a lego chart
-line        false                make a line chart
-pie         false                make a pie chart
-scatter     false                make a scatter chart
.....................................................................
-color       "lightsteelblue"     data color
-labelcolor  "rgb(100,100,100)"   label color
-valuecolor  "rgb(128,0,0)"       value color
-opacity     40                   opacity for area and wbar charts
-frame       0                    frame opacity
-font        ""                   specify font file (\"\": default)
.....................................................................
-h           1000                 canvas height
-w           1000                 canvas width
-left        20                   chart left
-top         80                   chart top
-bottom      20                   chart bottom
-right       80                   chart right
.....................................................................
-barwidth    0.5                  bar width
-dotsize     0.5                  bar width
-linewidth   0.25                 line width
-ls          2                    line spacing
-piesize     20                   pie chart radius
-textsize    1.5                  text size
.....................................................................
-chartitle   ""                   chart title
-ty          5                    title position relative to the top
-xlabel      1                    x-xaxis label interval (0: no labels)
-yfmt        "%v"                 yaxis format
-vfmt        ""                   value format ("": no values)
-yrange      ""                   y axis range (min,max,step)
.....................................................................
-grid        false                show y axis grid
-title       false                show the title
-zero        true                 zero minumum
......................................................................

`
	io.WriteString(os.Stderr, usage)
}

func main() {
	// Command line options

	// chart types
	flag.BoolVar(&opts.lego, "lego", false, "lego chart")
	flag.BoolVar(&opts.area, "area", false, "area chart")
	flag.BoolVar(&opts.bar, "bar", false, "bar chart")
	flag.BoolVar(&opts.dot, "dot", false, "dot chart")
	flag.BoolVar(&opts.line, "line", false, "line chart")
	flag.BoolVar(&opts.hbar, "hbar", false, "horizontal bar")
	flag.BoolVar(&opts.wbar, "wbar", false, "horizontal word bar")
	flag.BoolVar(&opts.scatter, "scatter", false, "scatter chart")
	flag.BoolVar(&opts.pie, "pie", false, "show a pie chart")
	// chart element sizes
	flag.Float64Var(&opts.barwidth, "barwidth", 0.5, "bar width")
	flag.Float64Var(&opts.dotsize, "dotsize", 0.5, "dot size")
	flag.Float64Var(&opts.linewidth, "linewidth", 0.25, "line width")
	flag.Float64Var(&opts.linespacing, "ls", opts.barwidth*4, "line spacing")
	flag.Float64Var(&opts.piesize, "piesize", 20, "pie chart radius")
	flag.Float64Var(&opts.textsize, "textsize", 1.5, "text size")
	// canvas sizes
	flag.IntVar(&screenWidth, "w", 1000, "canvas width")
	flag.IntVar(&screenHeight, "h", 1000, "canvas height")
	// chart positions
	flag.Float64Var(&opts.top, "top", 80, "chart top")
	flag.Float64Var(&opts.bottom, "bottom", 20, "chart bottom")
	flag.Float64Var(&opts.left, "left", 20, "chart left")
	flag.Float64Var(&opts.right, "right", 80, "chart right")
	// titles and axis settings
	flag.Float64Var(&opts.ty, "ty", 5, "title position relative to the top")
	flag.IntVar(&opts.xlabel, "xlabel", 1, "x-axis label interval")
	flag.StringVar(&opts.yrange, "yrange", "", "y axis range (min,max,step)")
	flag.StringVar(&opts.chartitle, "chartitle", "", "chart title")
	flag.StringVar(&opts.valuefmt, "vfmt", "", "value format (\"\": no values)")
	flag.StringVar(&opts.yaxfmt, "yfmt", "%v", "yaxis format (\"\" no y axis)")
	// colors and opacities
	flag.StringVar(&opts.dcolor, "color", "lightsteelblue", "color")
	flag.StringVar(&opts.bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&opts.fontname, "font", "", "font name")
	flag.StringVar(&opts.labelcolor, "labelcolor", "rgb(100,100,100)", "label color")
	flag.StringVar(&opts.valuecolor, "valuecolor", "rgb(128,100,0)", "value color")
	flag.Float64Var(&opts.frameOp, "frame", 0, "frame opacity (0: no frame)")
	flag.Float64Var(&opts.opacity, "opacity", 40, "% opacity for area and wbar charts")
	// on-off flags
	flag.BoolVar(&opts.showtitle, "title", true, "show the title")
	flag.BoolVar(&opts.showgrid, "grid", false, "show y axis grid")
	flag.BoolVar(&opts.zb, "zero", true, "zero minumum")
	flag.Usage = cmdUsage
	flag.Parse()

	var input io.Reader
	var ferr, err error
	var infile string

	// Read from stdin or specified file
	if len(flag.Args()) == 0 {
		input = os.Stdin
		infile = "stdin"
	} else {
		infile = flag.Args()[0]
		input, ferr = os.Open(infile)
		if ferr != nil {
			perr("unable to open ", infile)
			os.Exit(1)
		}
	}
	// read the data
	data, err = chart.DataRead(input)
	if err != nil {
		perr("unable to read ", infile)
		os.Exit(2)
	}
	// specify at least one of line, bar, hbar, scatter, area, pie, lego
	if !(opts.line || opts.scatter || opts.bar || opts.dot || opts.wbar || opts.area || opts.hbar || opts.lego || opts.pie) {
		perr("pick a chart type (-line, -bar, -hbar, -area, -scatter, -lego, -pie)", infile)
		os.Exit(3)
	}
	// make the chart
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("gchart")

	err = ebcanvas.LoadFont()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
