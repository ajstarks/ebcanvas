// ebdeck: show deck markup
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ajstarks/deck"
	"github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type App struct {
	slideNumber int
	nslides     int
	deckname    string
	d           deck.Deck
}

// command line options
type options struct {
	sansfont      string
	serifont      string
	monofont      string
	symbolfont    string
	layers        string
	pages         string
	pagesize      string
	fontdir       string
	gridpct       float64
	width, height int
}

// PageDimen describes page dimensions
// the unit field is used to convert to pt.
type PageDimen struct {
	width, height, unit float64
}

const (
	mm2pt        = 2.83464 // mm to pt conversion
	linespacing  = 1.8
	listspacing  = 1.5
	listwrap     = 95.0
	defaultColor = "rgb(128,128,128)"
)

var (
	btime                     time.Time
	codemap                   = strings.NewReplacer("\t", "    ") // convert tyabs to spaces
	opts                      options                             // command line options
	screenWidth, screenHeight int                                 // screen width, height

	imagecache = map[string]image.Image{}

	fontmap = map[string]*text.GoTextFaceSource{ // canonical names for fonts
		"sans":   ebcanvas.CurrentFont,
		"serif":  nil,
		"mono":   nil,
		"symbol": nil,
	}

	// pagemap defines page dimensions
	pagemap = map[string]PageDimen{
		"Letter":     {792, 612, 1},
		"Legal":      {1008, 612, 1},
		"Tabloid":    {1224, 792, 1},
		"ArchA":      {864, 648, 1},
		"Widescreen": {1152, 648, 1},
		"4R":         {432, 288, 1},
		"Index":      {360, 216, 1},
		"A2":         {420, 594, mm2pt},
		"A3":         {420, 297, mm2pt},
		"A4":         {297, 210, mm2pt},
		"A5":         {210, 148, mm2pt},
	}
)

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (a *App) Draw(screen *ebiten.Image) {
	ebdeck(a, screen)
}

func (a *App) Update() error {

	// if the deckfile has changed, reload
	t, err := modtime(a.deckname)
	if err == nil && t.After(btime) {
		a.dodeck(0, a.nslides, 0, 0)
	}

	// mouse wheel position
	_, wy := ebiten.Wheel()
	switch {
	// quit
	case inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		os.Exit(0)
	// grid
	case inpututil.IsKeyJustPressed(ebiten.KeyBracketLeft):
		opts.gridpct = 5
	case inpututil.IsKeyJustPressed(ebiten.KeyBracketRight):
		opts.gridpct = 0
	// refresh
	case inpututil.IsKeyJustPressed(ebiten.KeyR):
		a.dodeck(0, a.nslides, 0, 0)
	// home key -> first slide
	case inpututil.IsKeyJustPressed(ebiten.KeyHome):
		a.slideNumber = 0
	// end key -> last slide
	case inpututil.IsKeyJustPressed(ebiten.KeyEnd):
		a.slideNumber = a.nslides
	// move backwards
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyPageDown) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || wy > 0:
		a.slideNumber--
	// move forward
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyPageUp) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || wy < 0:
		a.slideNumber++
	}
	return nil
}

// modtime returns the modification time of a file
func modtime(filename string) (time.Time, error) {
	if filename == "" {
		return time.Time{}, nil
	}
	s, err := os.Stat(filename)
	return s.ModTime(), err
}

// process slides
func process(a *App, canvas *ebcanvas.Canvas) {
	ebiten.SetWindowTitle(fmt.Sprintf("%s: page %d of %d", a.deckname, a.slideNumber+1, a.nslides+1))
	slide := a.d.Slide[a.slideNumber]

	var bg color.NRGBA
	if slide.Bg == "" {
		bg = color.NRGBA{255, 255, 255, 255}
	} else {
		bg = ebcanvas.ColorLookup(slide.Bg)
	}
	if slide.Fg == "" {
		slide.Fg = "black"
	}
	canvas.Background(bg)

	// process each element according to the layer list
	layerlist := strings.Split(opts.layers, ":")
	for il := range layerlist {
		switch layerlist[il] { // for each element type loop through each set
		case "image":
			for _, i := range slide.Image {
				img, ok := imagecache[i.Name]
				if !ok {
					continue
				}
				dimage(canvas, img, i)
			}
		case "text":
			for _, t := range slide.Text {
				if t.Color == "" {
					t.Color = slide.Fg
				}
				dtext(canvas, t)
			}
		case "list":
			for _, li := range slide.List {
				list(canvas, li)
			}
		case "ellipse":
			for _, e := range slide.Ellipse {
				ellipse(canvas, e)
			}
		case "line":
			for _, l := range slide.Line {
				line(canvas, l)
			}
		case "rect":
			for _, r := range slide.Rect {
				rect(canvas, r)
			}
		case "poly":
			for _, p := range slide.Polygon {
				poly(canvas, p)
			}
		case "arc":
			for _, a := range slide.Arc {
				arc(canvas, a)
			}
		case "curve":
			for _, c := range slide.Curve {
				curve(canvas, c)
			}
		}
	}
	// add a grid, if specified
	if opts.gridpct > 0 {
		gc := ebcanvas.ColorLookup(slide.Fg)
		gc.A = 100
		canvas.Grid(0, 0, 100, 100, 0.1, float32(opts.gridpct), gc)
	}
}

// list processes lists
func list(canvas *ebcanvas.Canvas, list deck.List) {
	c := ebcanvas.ColorLookup(list.Color)
	c.A = setopacity(list.Opacity)
	var xp, yp, ls, ts float32
	xp = float32(list.Xp)
	yp = float32(list.Yp)
	ts = float32(list.Sp)
	ls = float32(list.Lp)
	if list.Font == "" {
		list.Font = "sans"
	}
	if ls == 0 {
		ls = listspacing
	}
	ebcanvas.CurrentFont = fontmap[list.Font]
	var t string
	for i, item := range list.Li {
		t = item.ListText
		if item.Font != "" {
			ebcanvas.CurrentFont = fontmap[item.Font]
		}
		if item.Color != "" {
			c = ebcanvas.ColorLookup(item.Color)
		}
		if list.Type == "number" {
			number(canvas, i, xp, yp, ts, c)
		}
		if list.Type == "bullet" {
			bullet(canvas, xp, yp, ts, c)
		}
		if list.Align == "center" {
			canvas.CText(xp, yp, ts, t, c)
		} else {
			canvas.Text(xp, yp, ts, t, c)
		}
		yp -= ls * ts * listspacing
	}
}

// arc makes arcs
func arc(canvas *ebcanvas.Canvas, a deck.Arc) {
	if a.Color == "" {
		a.Color = defaultColor
	}
	c := ebcanvas.ColorLookup(a.Color)
	c.A = setopacity(a.Opacity)
	canvas.Arc(float32(a.Xp), float32(a.Yp), float32(a.Wp/2), float32(a.A1), float32(a.A2), c)
}

// curve makea a quad bezier curve
func curve(canvas *ebcanvas.Canvas, curve deck.Curve) {
	if curve.Color == "" {
		curve.Color = defaultColor
	}
	c := ebcanvas.ColorLookup(curve.Color)
	c.A = setopacity(curve.Opacity)
	x1, y1 := float32(curve.Xp1), float32(curve.Yp1)
	x2, y2 := float32(curve.Xp2), float32(curve.Yp2)
	x3, y3 := float32(curve.Xp3), float32(curve.Yp3)
	sw := float32(curve.Sp)
	canvas.StrokedCurve(x1, y1, x2, y2, x3, y3, sw, c)
}

// rect makes rectangles and squares
func rect(canvas *ebcanvas.Canvas, r deck.Rect) {
	if r.Color == "" {
		r.Color = defaultColor
	}
	c := ebcanvas.ColorLookup(r.Color)
	c.A = setopacity(r.Opacity)
	x, y, w, h := float32(r.Xp), float32(r.Yp), float32(r.Wp), float32(r.Hp)
	if r.Hr == 100 {
		canvas.Square(x, y, w, c)
	} else {
		canvas.CenterRect(x, y, w, h, c)
	}
}

// poly makes a filled polygon
func poly(canvas *ebcanvas.Canvas, p deck.Polygon) {
	xs := strings.Split(p.XC, " ")
	ys := strings.Split(p.YC, " ")
	if len(xs) != len(ys) {
		return
	}
	if len(xs) < 3 || len(ys) < 3 {
		return
	}
	xp := make([]float32, len(xs))
	yp := make([]float32, len(ys))
	for i := range xs {
		x, err := strconv.ParseFloat(xs[i], 64)
		if err != nil {
			xp[i] = 0
		} else {
			xp[i] = float32(x)
		}
		y, err := strconv.ParseFloat(ys[i], 64)
		if err != nil {
			yp[i] = 0
		} else {
			yp[i] = float32(y)
		}
	}
	if p.Color == "" {
		p.Color = defaultColor
	}
	c := ebcanvas.ColorLookup(p.Color)
	c.A = setopacity(p.Opacity)
	canvas.Polygon(xp, yp, c)
}

// ellipse makes circles (for now)
func ellipse(canvas *ebcanvas.Canvas, e deck.Ellipse) {
	if e.Hr != 100 {
		return
	}
	if e.Color == "" {
		e.Color = defaultColor
	}
	c := ebcanvas.ColorLookup(e.Color)
	c.A = setopacity(e.Opacity)
	canvas.Circle(float32(e.Xp), float32(e.Yp), float32(e.Wp/2), c)
}

// line makes lines
func line(canvas *ebcanvas.Canvas, l deck.Line) {
	if l.Color == "" {
		l.Color = defaultColor
	}
	c := ebcanvas.ColorLookup(l.Color)
	c.A = setopacity(l.Opacity)
	canvas.Line(float32(l.Xp1), float32(l.Yp1), float32(l.Xp2), float32(l.Yp2), float32(l.Sp), c)
}

// dtext processes text
func dtext(canvas *ebcanvas.Canvas, t deck.Text) {
	if t.Font == "" {
		t.Font = "sans"
	}
	x, y, ts := float32(t.Xp), float32(t.Yp), float32(t.Sp)
	c := ebcanvas.ColorLookup(t.Color)
	c.A = setopacity(t.Opacity)
	ebcanvas.CurrentFont = fontmap[t.Font]

	s := t.Tdata
	if t.Type == "block" {
		canvas.TextWrap(x, y, float32(t.Wp), ts, s, c)
		return
	}
	if len(t.File) > 0 {
		tl := strings.Split(includefile(t.File), "\n")
		if t.Type == "code" {
			ebcanvas.CurrentFont = fontmap["mono"]
			ch := float64(len(tl)) * linespacing * float64(ts)
			canvas.CornerRect(x-ts, y+(ts*2), float32(t.Wp), float32(ch), color.NRGBA{240, 240, 240, 255})
		}
		textlines(canvas, x, y, ts, tl, c)
		return
	}
	if t.Rotation > 0 {
		canvas.RText(x, y, float32(t.Rotation), ts, s, c)
		return
	}
	switch t.Align {
	case "c", "middle", "mid", "center":
		canvas.CText(x, y, ts, s, c)
	case "e", "right", "end":
		canvas.EText(x, y, ts, s, c)
	default:
		canvas.Text(x, y, ts, s, c)
	}
}

// dimage processes deck images
func dimage(canvas *ebcanvas.Canvas, img image.Image, i deck.Image) {
	var sc float32
	sc = 100.0
	if i.Scale > 0 {
		sc = float32(i.Scale)
	}
	if i.Height == 0 {
		sc = float32(i.Width)
	}
	canvas.CenterImage(float32(i.Xp), float32(i.Yp), sc, img)
}

// setopacity sets the alpha value:
// 0 == default value (opaque)
// -1 == fully transparent
// > 0 set opacity percent
func setopacity(v float64) uint8 {
	var o uint8
	switch {
	case v < 0:
		o = 0
	case v > 0:
		o = uint8(255.0 * (v / 100))
	case v == 0:
		o = 255
	}
	return o
}

// bullet draws a bullet for a list item.
func bullet(canvas *ebcanvas.Canvas, x, y, size float32, c color.NRGBA) {
	canvas.Circle(x-size, y+size/2, size/4, c)
}

// number adds a number for a list item.
func number(canvas *ebcanvas.Canvas, n int, x, y, size float32, c color.NRGBA) {
	canvas.EText(x-size/2, y, size, fmt.Sprintf("%d.", n+1), c)
}

// textlines shows a series of lines of text
func textlines(canvas *ebcanvas.Canvas, x, y, size float32, s []string, c color.NRGBA) {
	yp := y
	for _, t := range s {
		canvas.Text(x, yp, size, t, c)
		yp -= linespacing * size
	}
}

// includefile reads the content of a file into a string
func includefile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return ""
	}
	return codemap.Replace(string(data))
}

// ebdeck processes deck data
func ebdeck(a *App, screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight

	a.nslides = len(a.d.Slide) - 1
	if a.slideNumber > a.nslides {
		a.slideNumber = 0
	}
	if a.slideNumber < 0 {
		a.slideNumber = a.nslides
	}
	process(a, canvas)
}

// imageinfo returns an image from a named file
func imageInfo(s string) image.Image {
	f, err := os.Open(s)
	defer f.Close()
	if err != nil {
		return nil
	}
	im, _, err := image.Decode(f)
	if err != nil {
		return nil
	}
	return im
}

// setpagesize parses the page size string (wxh)
func setpagesize(s string) (float64, float64) {
	var width, height float64
	var err error
	d := strings.FieldsFunc(s, func(c rune) bool { return !unicode.IsNumber(c) })
	if len(d) != 2 {
		return 0, 0
	}
	width, err = strconv.ParseFloat(d[0], 64)
	if err != nil {
		return 0, 0
	}
	height, err = strconv.ParseFloat(d[1], 64)
	if err != nil {
		return 0, 0
	}
	return width, height
}

// pagerange returns the begin and end using a "-" string
func pagerange(s string) (int, int) {
	p := strings.Split(s, "-")
	if len(p) != 2 {
		return 0, 0
	}
	b, berr := strconv.Atoi(p[0])
	e, err := strconv.Atoi(p[1])
	if berr != nil || err != nil {
		return 0, 0
	}
	if b > e {
		return 0, 0
	}
	return b, e
}

// coord makes coordinates
func coord(canvas *ebcanvas.Canvas, x, y float64) {
	s := fmt.Sprintf("(%.0f, %.0f)", x, y)
	canvas.CText(float32(x), float32(y-1.5), 1.5, s, color.NRGBA{0, 0, 0, 255})
}

// setfontdir determines the font directory:
// if the string argument is non-empty, use that, otherwise
// use the contents of the DECKFONT environment variable,
// if that is not set, or empty, use $HOME/deckfonts
func setfontdir(s string) string {
	if len(s) > 0 {
		return s
	}
	envdef := os.Getenv("DECKFONTS")
	if len(envdef) > 0 {
		return envdef
	}
	return path.Join(os.Getenv("HOME"), "deckfonts")
}

// loadDeckFont gets fonts from the font directory
func loadDeckFont(dname, name string) {
	f, err := ebcanvas.LoadFontName(path.Join(opts.fontdir, name) + ".ttf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	fontmap[dname] = f
}

func (a *App) updateDeck() (io.ReadCloser, error) {
	var err error
	var r io.ReadCloser
	if a.deckname == "" {
		r = os.Stdin
		a.deckname = "Standard-Input"
	} else {
		r, err = os.Open(a.deckname)
	}
	return r, err
}

// dodeck reads a deck, caching all images,
func (a *App) dodeck(begin, end int, pw, ph float64) {
	r, err := a.updateDeck()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	d, err := deck.ReadDeck(r, 0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	// cache all images
	ns := len(d.Slide)
	a.nslides = ns
	a.d = d
	for i := range ns {
		slide := d.Slide[i]
		for j := range slide.Image {
			iname := slide.Image[j].Name
			img := imageInfo(iname)
			if img == nil {
				continue
			}
			_, ok := imagecache[iname]
			if !ok {
				imagecache[iname] = img
			}
		}
	}

	r.Close()
}

func main() {
	// parse command line options
	flag.StringVar(&opts.sansfont, "sans", "PublicSans-Regular", "sans font")
	flag.StringVar(&opts.monofont, "mono", "Inconsolata-Medium", "mono font")
	flag.StringVar(&opts.serifont, "serif", "Charter-Regular", "sans font")
	flag.StringVar(&opts.symbolfont, "symbol", "ZapfDingbats", "sans font")
	flag.StringVar(&opts.layers, "layers", "image:rect:ellipse:curve:arc:line:poly:text:list", "Layer order")
	flag.StringVar(&opts.pagesize, "pagesize", "Letter", "pagesize: w,h, or one of: Letter, Legal, Tabloid, A3, A4, A5, ArchA, 4R, Index, Widescreen")
	flag.StringVar(&opts.pages, "pages", "1-1000000", "page range (first-last)")
	flag.StringVar(&opts.fontdir, "fontdir", setfontdir(""), "directory for fonts")
	flag.Float64Var(&opts.gridpct, "grid", 0, "grid size (0 for no grid)")
	flag.Parse()

	loadDeckFont("sans", opts.sansfont)
	loadDeckFont("serif", opts.serifont)
	loadDeckFont("mono", opts.monofont)
	loadDeckFont("symbol", opts.symbolfont)
	pw, ph := setpagesize(opts.pagesize)
	begin, end := pagerange(opts.pages)
	if pw == 0 && ph == 0 {
		p, ok := pagemap[opts.pagesize]
		if !ok {
			p = pagemap["Letter"]
		}
		pw = p.width * p.unit
		ph = p.height * p.unit
	}
	ebcanvas.CurrentFont = fontmap["sans"]

	// read decks from a named file or stdin
	a := new(App)
	files := flag.Args()
	if len(files) < 1 {
		a.deckname = ""
	} else {
		a.deckname = files[0]
	}
	var err error
	btime, err = modtime(a.deckname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	screenWidth, screenHeight = int(pw), int(ph)
	a.dodeck(begin, end, pw, ph)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(a); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	}
}
