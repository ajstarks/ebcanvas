// elections: show election results on a state grid
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type App struct {
	electionNumber int
	ne             int
}

// Data file structure
type egrid struct {
	name       string
	party      string
	row        int
	col        int
	population int
}

// One election "frame"
type election struct {
	title    string
	min, max int
	data     []egrid
}

// command line options
type options struct {
	width, height               int
	top, left, rowsize, colsize float64
	bgcolor, textcolor, shape   string
}

var (
	elections                 []election
	opts                      options
	screenWidth, screenHeight int
	partyColors               = map[string]string{
		"r":  "red",
		"d":  "blue",
		"i":  "gray",
		"w":  "peru",
		"dr": "purple",
		"f":  "green",
	}
)

// maprange maps one range into another
func maprange(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func (a *App) Update() error {
	_, wy := ebiten.Wheel()

	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		os.Exit(0)

	case inpututil.IsKeyJustPressed(ebiten.KeyHome):
		a.electionNumber = 0

	case inpututil.IsKeyJustPressed(ebiten.KeyEnd):
		a.electionNumber = a.ne

	case inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyPageDown) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || wy < 0:
		a.electionNumber--

	case inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyPageUp) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || wy > 0:
		a.electionNumber++
	}
	return nil
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (a *App) Draw(screen *ebiten.Image) {
	elect(a, screen)
}

func million(n int) string {
	s := strconv.FormatInt(int64(n), 10)
	p := len(s)
	return s[0:p-6] + "," + s[p-6:p-3] + "," + s[p-3:p]
}

// area computes the area of a circle
func area(v float64) float64 {
	return math.Sqrt((v / math.Pi)) * 2
}

// atoi converts a string to an integer
func atoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

// readData reads election data into the data structure
func readData(r io.Reader) (election, error) {
	var d egrid
	var data []egrid
	var min, max int
	title := ""
	scanner := bufio.NewScanner(r)
	min, max = math.MaxInt32, -math.MaxInt32
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 { // skip blank lines
			continue
		}
		if t[0] == '#' && len(t) > 2 { // get the title
			title = t[2:]
			continue
		}
		fields := strings.Split(t, "\t")
		if len(fields) < 5 { // skip incomplete records
			continue
		}
		// name,col,row,party,population
		d.name = fields[0]
		d.col = atoi(fields[1])
		d.row = atoi(fields[2])
		d.party = fields[3]
		d.population = atoi(fields[4])
		data = append(data, d)
		// compute min, max
		if d.population > max {
			max = d.population
		}
		if d.population < min {
			min = d.population
		}
	}
	var e election
	e.title = title
	e.min = min
	e.max = max
	e.data = data
	return e, scanner.Err()
}

func process(canvas *ebcanvas.Canvas, e election) {
	beginPage(canvas, opts.bgcolor)
	fmin, fmax := float64(e.min), float64(e.max)
	amin, amax := area(fmin), area(fmax)
	sumpop := 0
	for _, d := range e.data {
		sumpop += d.population
		x := opts.left + (float64(d.row) * opts.colsize)
		y := opts.top - (float64(d.col) * opts.rowsize)
		fpop := float64(d.population)
		apop := area(fpop)

		// defaults
		txcolor := "white"
		txsize := 1.2
		font := "sans"
		name := d.name

		switch opts.shape {
		case "c": // circle
			r := maprange(apop, amin, amax, 2, opts.colsize)
			circle(canvas, x, y, r, partyColors[d.party])
		case "h": // hexagom
			r := maprange(apop, amin, amax, 2, opts.colsize)
			hexagon(canvas, x, y, r/2, partyColors[d.party])
		case "s": // square
			r := maprange(fpop, fmin, fmax, 2, opts.colsize)
			square(canvas, x, y, r, partyColors[d.party])
		case "l": // lines
			r := maprange(apop, amin, amax, 2, opts.colsize)
			polylines(canvas, x, y, r/2, 0.25, partyColors[d.party])
			txcolor = partyColors[d.party]
		case "p": // plain text
			txcolor = partyColors[d.party]
			txsize = maprange(fpop, fmin, fmax, 2, opts.colsize*0.75)
			//		case "g": // geographic
			//			txcolor = partyColors[d.party]
			//			name = statemap[d.name]
			//			font = "symbol"
			//			txsize = maprange(fpop, fmin, fmax, 2, opts.colsize)
		default:
			r := maprange(apop, amin, amax, 2, opts.colsize)
			circle(canvas, x, y, r, partyColors[d.party])
		}
		ctext(canvas, x, y-0.5, txsize, name, font, txcolor)
	}
	showtitle(canvas, e.title, sumpop, opts.top+15, opts.textcolor)
	endPage(canvas)
}

func partycand(s, def string) (string, string) {
	var party, cand string
	f := strings.Split(s, ":")
	if len(f) > 1 {
		party = f[1]
		cand = f[0]
	} else {
		party = def
		cand = s
	}
	return party, cand
}

// showtitle shows the title and subhead
func showtitle(canvas *ebcanvas.Canvas, s string, pop int, top float64, textcolor string) {
	fields := strings.Fields(s) // year, democratic, republican, third-party (optional)
	if len(fields) < 2 {
		return
	}
	suby := top - 7
	ctext(canvas, 50, top, 3.6, fields[0]+" US Presidential Election", "sans", textcolor)
	ctext(canvas, 85, 5, 1.5, "Population: "+million(pop), "sans", textcolor)
	var party string
	var cand string
	if len(fields) > 1 {
		party, cand = partycand(fields[1], "d")
		legend(canvas, 20, suby, 2.0, cand, partyColors[party], textcolor)
	}
	if len(fields) > 2 {
		party, cand = partycand(fields[2], "r")
		legend(canvas, 80, suby, 2.0, cand, partyColors[party], textcolor)
	}
	if len(fields) > 3 {
		party, cand = partycand(fields[3], "i")
		legend(canvas, 50, suby, 2.0, cand, partyColors[party], textcolor)
	}

}

// circle makes a circle
func circle(canvas *ebcanvas.Canvas, x, y, r float64, color string) {
	cx, cy, cr := float32(x), float32(y), float32(r)
	canvas.Circle(cx, cy, cr/2, ebcanvas.ColorLookup(color))
}

// ctext makes centered text
func ctext(canvas *ebcanvas.Canvas, x, y, size float64, s string, fontname string, color string) {
	tx, ty, ts := float32(x), float32(y), float32(size)
	canvas.CText(tx, ty, ts, s, ebcanvas.ColorLookup(color))
}

// ltext makes left-aligned text
func ltext(canvas *ebcanvas.Canvas, x, y, size float64, s string, color string) {
	tx, ty, ts := float32(x), float32(y), float32(size)
	canvas.Text(tx, ty, ts, s, ebcanvas.ColorLookup(color))
}

// square makes a square centered ar (x,y), width w.
func square(canvas *ebcanvas.Canvas, x, y, w float64, color string) {
	canvas.Square(float32(x), float32(y), float32(w), ebcanvas.ColorLookup(color))
}

// pangles computes the points of a polygon based on a series of angles
func pangles(cx, cy, r float64, angles []float64) ([]float32, []float32) {
	px := make([]float32, len(angles))
	py := make([]float32, len(angles))
	aspect := float64(screenWidth) / float64(screenHeight)
	for i, a := range angles {
		t := a * (math.Pi / 180)
		px[i] = float32(cx + (r * math.Cos(t)))
		py[i] = float32(cy + ((r * aspect) * math.Sin(t)))
	}
	return px, py
}

// hexagon makes a filled hexagon centered at (cx, cy), size is the subscribed circle radius r
func hexagon(canvas *ebcanvas.Canvas, cx, cy, r float64, color string) {
	px, py := pangles(cx, cy, r, []float64{30, 90, 150, 210, 270, 330})
	canvas.Polygon(px, py, ebcanvas.ColorLookup(color))
}

// polylines makes a outlined hexagon, centered at (cx, cy), size is the subscribed circle radius r
func polylines(canvas *ebcanvas.Canvas, cx, cy, r, lw float64, color string) {
	px, py := pangles(cx, cy, r, []float64{30, 90, 150, 210, 270, 330})
	lx := len(px) - 1
	linewidth := float32(lw)
	for i := 0; i < lx; i++ {
		canvas.Line(px[i], py[i], px[i+1], py[i+1], linewidth, ebcanvas.ColorLookup(color))
	}
	canvas.Line(px[0], py[0], px[lx], py[lx], linewidth, ebcanvas.ColorLookup(color))
}

// legend makes the subtitle
func legend(canvas *ebcanvas.Canvas, x, y, ts float64, s string, color, textcolor string) {
	ltext(canvas, x, y, ts, s, textcolor)
	circle(canvas, x-ts, y+ts/3, ts/2, color)
}

// beginPage starts a page
func beginPage(canvas *ebcanvas.Canvas, bgcolor string) {
	canvas.Background(ebcanvas.ColorLookup(bgcolor))
}

// endPage ends a page
func endPage(canvas *ebcanvas.Canvas) {
	ctext(canvas, 50, 5, 1.5, "The area of a circle denotes state population: source U.S. Census", "sans", "gray")
}

// elect processes election data
func elect(a *App, screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight
	a.ne = len(elections) - 1
	if a.electionNumber > a.ne {
		a.electionNumber = 0
	}
	if a.electionNumber < 0 {
		a.electionNumber = a.ne
	}
	process(canvas, elections[a.electionNumber])
}

func main() {
	// parse command line options
	flag.Float64Var(&opts.top, "top", 75, "map top value (canvas %)")
	flag.Float64Var(&opts.left, "left", 15, "map left value (canvas %)")
	flag.Float64Var(&opts.rowsize, "rowsize", 9, "rowsize (canvas %)")
	flag.Float64Var(&opts.colsize, "colsize", 7, "column size (canvas %)")
	flag.IntVar(&screenWidth, "width", 1200, "canvas width")
	flag.IntVar(&screenHeight, "height", 900, "canvas height")
	flag.StringVar(&opts.bgcolor, "bgcolor", "black", "background color")
	flag.StringVar(&opts.textcolor, "textcolor", "white", "text color")
	flag.StringVar(&opts.shape, "shape", "c", "shape for states:\n\"c\": circle,\n\"h\": hexagon,\n\"s\": square\n\"l\": line\n\"g\": geographic\n\"p\": plain text")

	flag.Parse()

	// Read in the data
	for _, f := range flag.Args() {
		r, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		e, err := readData(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		elections = append(elections, e)
		r.Close()
	}
	if len(elections) < 1 {
		fmt.Fprintln(os.Stderr, "no data read")
		os.Exit(1)
	}
	if err := ebcanvas.LoadFont(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("elections")
	if err := ebiten.RunGame(&App{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	}

}
