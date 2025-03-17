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

	ec "github.com/ajstarks/ebcanvas"
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
	population int64
}

// One election "frame"
type election struct {
	title    string
	min, max int64
	data     []egrid
}

// command line options
type options struct {
	width, height               int
	top, left, rowsize, colsize float64
	bgcolor, textcolor          string
}

var (
	elections                 []election
	opts                      options
	partyColors               = map[string]string{"r": "red", "d": "blue", "i": "gray"}
	screenWidth, screenHeight int
)

func (a *App) Update() error {
	_, wy := ebiten.Wheel()

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		a.electionNumber = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) ||
		wy > 0 {
		a.electionNumber--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) ||
		wy < 0 {
		a.electionNumber++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		a.electionNumber = a.ne
	}
	return nil
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (a *App) Draw(screen *ebiten.Image) {
	elect(a, screen)
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

// atoi64 converts a string to an 64-bit integer
func atoi64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// readData reads election data into the data structure
func readData(r io.Reader) (election, error) {
	var d egrid
	var data []egrid
	var min, max int64
	title := ""
	scanner := bufio.NewScanner(r)
	min, max = math.MaxInt64, -math.MaxInt64
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
		d.population = atoi64(fields[4])
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

// process walks the data, making the visualization
func process(canvas *ec.Canvas, e election) {
	amin := area(float64(e.min))
	amax := area(float64(e.max))
	beginPage(canvas, opts.bgcolor)
	showtitle(canvas, e.title, opts.top+15, opts.textcolor)
	for _, d := range e.data {
		x := opts.left + (float64(d.row) * opts.colsize)
		y := opts.top - (float64(d.col) * opts.rowsize)
		r := ec.MapRange(area(float64(d.population)), amin, amax, 2, opts.colsize)
		circle(canvas, x, y, r, partyColors[d.party])
		ctext(canvas, x, y-0.5, 1.2, d.name, "white")
	}
	endPage(canvas)
}

// showtitle shows the title and subhead
func showtitle(canvas *ec.Canvas, s string, top float64, textcolor string) {
	fields := strings.Fields(s) // year, democratic, republican, third-party (optional)
	if len(fields) < 3 {
		return
	}
	suby := top - 7
	ctext(canvas, 50, top, 3.6, fields[0]+" US Presidential Election", textcolor)
	legend(canvas, 20, suby, 2.0, fields[1], partyColors["d"], textcolor)
	legend(canvas, 80, suby, 2.0, fields[2], partyColors["r"], textcolor)
	if len(fields) > 3 {
		legend(canvas, 50, suby, 2.0, fields[3], partyColors["i"], textcolor)
	}
}

// circle makes a circle
func circle(canvas *ec.Canvas, x, y, r float64, color string) {
	cx, cy, cr := float32(x), float32(y), float32(r)
	canvas.Circle(cx, cy, cr/2, ec.ColorLookup(color))
}

// ctext makes centered text
func ctext(canvas *ec.Canvas, x, y, size float64, s string, color string) {
	tx, ty, ts := float32(x), float32(y), float32(size)
	canvas.CText(tx, ty, ts, s, ec.ColorLookup(color))
}

// ltext makes left-aligned text
func ltext(canvas *ec.Canvas, x, y, size float64, s string, color string) {
	tx, ty, ts := float32(x), float32(y), float32(size)
	canvas.Text(tx, ty, ts, s, ec.ColorLookup(color))
}

// legend makes the subtitle
func legend(canvas *ec.Canvas, x, y, ts float64, s string, color, textcolor string) {
	ltext(canvas, x, y, ts, s, textcolor)
	circle(canvas, x-ts, y+ts/3, ts/2, color)
}

// beginPage starts a page
func beginPage(canvas *ec.Canvas, bgcolor string) {
	canvas.Background(ec.ColorLookup(bgcolor))
}

// endPage ends a page
func endPage(canvas *ec.Canvas) {
	ctext(canvas, 50, 5, 1.5, "The area of a circle denotes state population: source U.S. Census", "gray")
}

// elect processes election data
func elect(a *App, screen *ebiten.Image) {
	canvas := new(ec.Canvas)
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
	flag.Float64Var(&opts.left, "left", 7, "map left value (canvas %)")
	flag.Float64Var(&opts.rowsize, "rowsize", 9, "rowsize (canvas %)")
	flag.Float64Var(&opts.colsize, "colsize", 7, "column size (canvas %)")
	flag.IntVar(&screenWidth, "width", 1200, "canvas width")
	flag.IntVar(&screenHeight, "height", 900, "canvas height")
	flag.StringVar(&opts.bgcolor, "bgcolor", "black", "background color")
	flag.StringVar(&opts.textcolor, "textcolor", "white", "text color")
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
	if err := ec.LoadFont(); err != nil {
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
