// elections: show election results on a state grid
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
	linespacing  = 1.4
	listspacing  = 2.0
	fontfactor   = 1.0
	listwrap     = 95.0
	defaultColor = "rgb(128,128,128)"
)

var (
	opts                      options // command line options
	screenWidth, screenHeight int     // screen width, height

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
	_, wy := ebiten.Wheel()

	switch {
	// quit
	case inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		os.Exit(0)
	// home key -> first result
	case inpututil.IsKeyJustPressed(ebiten.KeyHome):
		a.slideNumber = 0
	// end key -> last result
	case inpututil.IsKeyJustPressed(ebiten.KeyEnd):
		a.slideNumber = a.nslides
	// move backwards
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyPageDown) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || wy < 0:
		a.slideNumber--
	// move forwardq
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) ||
		inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) ||
		inpututil.IsKeyJustPressed(ebiten.KeyPageUp) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || wy > 0:
		a.slideNumber++
	}
	return nil
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

// process makes an election result
func process(a *App, canvas *ebcanvas.Canvas) {
	var bg color.NRGBA
	slide := a.d.Slide[a.slideNumber]
	if slide.Bg == "" {
		bg = color.NRGBA{255, 255, 255, 255}
	} else {
		bg = ebcanvas.ColorLookup(slide.Bg)
	}
	if slide.Fg == "" {
		slide.Fg = "black"
	}
	canvas.Background(bg)
	layerlist := strings.Split(opts.layers, ":")
	// process each element
	for il := range layerlist {
		switch layerlist[il] {
		case "image":
			for _, i := range slide.Image {
				img, ok := imagecache[i.Name]
				if !ok {
					continue
				}
				doimage(canvas, img, i)
			}
		case "text":
			for _, t := range slide.Text {
				if t.Color == "" {
					t.Color = slide.Fg
				}
				dotext(canvas, t)
			}
		case "list":
			for _, li := range slide.List {
				dolist(canvas, li)
			}
		case "ellipse":
			for _, e := range slide.Ellipse {
				doellipse(canvas, e)
			}
		case "line":

			for _, l := range slide.Line {
				if l.Color == "" {
					l.Color = slide.Fg
				}

				doline(canvas, l)
			}
		case "rect":
			for _, r := range slide.Rect {
				dorect(canvas, r)
			}
		case "poly":
			for _, p := range slide.Polygon {
				dopoly(canvas, p)
			}
		case "arc":
			for _, a := range slide.Arc {
				doarc(canvas, a)
			}
		case "curve":
			for _, c := range slide.Curve {
				docurve(canvas, c)
			}
		}
	}
	if opts.gridpct > 0 {
		gc := ebcanvas.ColorLookup(slide.Fg)
		gc.A = 100
		canvas.Grid(0, 0, 100, 100, 0.1, float32(opts.gridpct), gc)
	}
}

// bullet draws a bullet
func bullet(canvas *ebcanvas.Canvas, x, y, size float32, c color.NRGBA) {
	canvas.Circle(x-size, y+size/2, size/4, c)
}

func number(canvas *ebcanvas.Canvas, n int, x, y, size float32, c color.NRGBA) {
	canvas.EText(x-size/2, y, size, fmt.Sprintf("%d.", n+1), c)
}

func dolist(canvas *ebcanvas.Canvas, list deck.List) {
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
		ls = linespacing
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
		yp -= ls * ts * 2
	}
}

func doimage(canvas *ebcanvas.Canvas, img image.Image, i deck.Image) {
	sc := 100.0
	if i.Scale > 0 {
		sc = i.Scale
	}
	if i.Height == 0 {
		sc = float64(i.Width)
	}
	canvas.CenterImage(float32(i.Xp), float32(i.Yp), float32(sc), img)
}

func doarc(canvas *ebcanvas.Canvas, a deck.Arc) {
	if a.Color == "" {
		a.Color = defaultColor
	}
	c := ebcanvas.ColorLookup(a.Color)
	c.A = setopacity(a.Opacity)
	cx, cy, r, a1, a2 := float32(a.Xp), float32(a.Yp), float32(a.Wp/2), float32(a.A1), float32(a.A2)
	canvas.Arc(cx, cy, r, a1, a2, c)
}

func docurve(canvas *ebcanvas.Canvas, c deck.Curve) {
	if c.Color == "" {
		c.Color = defaultColor
	}
	clr := ebcanvas.ColorLookup(c.Color)
	clr.A = setopacity(c.Opacity)
	x1, y1 := float32(c.Xp1), float32(c.Yp1)
	x2, y2 := float32(c.Xp2), float32(c.Yp2)
	x3, y3 := float32(c.Xp3), float32(c.Yp3)
	sw := float32(c.Sp)
	canvas.StrokedCurve(x1, y1, x2, y2, x3, y3, sw, clr)
}

func dorect(canvas *ebcanvas.Canvas, r deck.Rect) {
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

func dopoly(canvas *ebcanvas.Canvas, p deck.Polygon) {
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

func doellipse(canvas *ebcanvas.Canvas, e deck.Ellipse) {
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

func doline(canvas *ebcanvas.Canvas, l deck.Line) {
	if l.Color == "" {
		l.Color = defaultColor
	}
	c := ebcanvas.ColorLookup(l.Color)
	c.A = setopacity(l.Opacity)
	canvas.Line(float32(l.Xp1), float32(l.Yp1), float32(l.Xp2), float32(l.Yp2), float32(l.Sp), c)
}

func dotext(canvas *ebcanvas.Canvas, t deck.Text) {
	if t.Font == "" {
		t.Font = "sans"
	}
	x, y, ts := float32(t.Xp), float32(t.Yp), float32(t.Sp)
	c := ebcanvas.ColorLookup(t.Color)
	c.A = setopacity(t.Opacity)
	ebcanvas.CurrentFont = fontmap[t.Font]

	if t.Type == "block" {
		canvas.TextWrap(x, y, float32(t.Wp), ts, t.Tdata, c)
		return
	}

	switch t.Align {
	case "c", "middle", "mid", "center":
		canvas.CText(x, y, ts, t.Tdata, c)
	case "e", "right", "end":
		canvas.EText(x, y, ts, t.Tdata, c)
	default:
		canvas.Text(x, y, ts, t.Tdata, c)
	}
}

// elect processes election data
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

// imageinfo returns an image dimensions
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

func loadDeckFont(dname, name string) {
	f, err := ebcanvas.LoadFontName(path.Join(opts.fontdir, name) + ".ttf")
	if err != nil {
		fontmap[dname] = ebcanvas.CurrentFont
		fmt.Fprintf(os.Stderr, "%v (default font used for %s)\n", err, dname)
	}
	fontmap[dname] = f
}

func (a *App) dodeck(r io.ReadCloser, begin, end int, pw, ph float64) {
	d, err := deck.ReadDeck(r, 0, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(5)
	}
	// cache all images
	ns := len(d.Slide)
	for i := range ns {
		slide := d.Slide[i]
		for j := range slide.Image {
			iname := slide.Image[j].Name
			img := imageInfo(iname)
			if img == nil {
				continue
			}
			imagecache[iname] = img
		}
	}
	a.d = d
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

	var err error

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

	var r io.ReadCloser
	files := flag.Args()
	if len(files) < 1 {
		r = os.Stdin
	} else {
		r, err = os.Open(files[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(2)
		}
	}

	ebcanvas.CurrentFont = fontmap["sans"]
	a := new(App)
	a.dodeck(r, begin, end, pw, ph)
	screenWidth, screenHeight = int(pw), int(ph)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ebdeck")
	if err := ebiten.RunGame(a); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	}
}
