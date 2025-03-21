package ebcanvas

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const Pi = 3.14159265358979323846264338327950288419716939937510582097494459

type Canvas struct {
	Width, Height int
	Screen        *ebiten.Image
}

var (
	mplusFaceSource *text.GoTextFaceSource
)

// LoadFont loads the default font
func LoadFont() error {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		return err
	}
	mplusFaceSource = s
	return nil
}

// pct calculates percentage of a value
func pct(p float32, m float32) float32 {
	return ((p / 100.0) * m)
}

// dimen calculates dimensions based on percentages
func dimen(xp, yp, w, h float32) (float32, float32) {
	return pct(xp, w), pct(100-yp, h)
}

// Absolute methods

// arc draws a filled arc centered at (cx,cy) with radius r, between angle a1 and a2
func arc(screen *ebiten.Image, cx, cy, r, a1, a2 float32, fillcolor color.NRGBA) {
	var p vector.Path
	p.Arc(cx, cy, r, a1, a2, vector.CounterClockwise)
	vector.DrawFilledPath(screen, &p, fillcolor, true, vector.FillRuleEvenOdd)
}

// degreesToRadians converts degrees (0-360 counter-clockwise) to the
// radian measures used by ebiten/vector
func degreesToRadians(deg float32) float32 {
	return (360 - deg) * (Pi / 180)
}

// strokedarc strokes an arc centered at (cx,cy) with radius r,
// between angles a1 and a2 (degrees 0-360, counter-clockwise)
func strokedarc(screen *ebiten.Image, cx, cy, r, a1, a2, size float32, strokecolor color.NRGBA) {
	var p vector.Path
	op := vector.StrokeOptions{Width: size}
	p.Arc(cx, cy, r, a1, a2, vector.CounterClockwise)
	vector.StrokePath(screen, &p, strokecolor, true, &op)
}

// btext draws text beginning at (x,y)
func btext(screen *ebiten.Image, x, y float64, size float64, s string, textcolor color.NRGBA) {
	ff := &text.GoTextFace{Source: mplusFaceSource, Size: size}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y-size)
	op.ColorScale.ScaleWithColor(textcolor)
	text.Draw(screen, s, ff, op)
}

// ctext draws text centered at (x,y)
func ctext(screen *ebiten.Image, x, y float64, size float64, s string, textcolor color.NRGBA) {
	ff := &text.GoTextFace{Source: mplusFaceSource, Size: size}
	tw := text.Advance(s, ff)
	op := &text.DrawOptions{}
	op.GeoM.Translate(x-(tw/2), y-size)
	op.ColorScale.ScaleWithColor(textcolor)
	text.Draw(screen, s, ff, op)
}

// etext draws text with end point at (x,y)
func etext(screen *ebiten.Image, x, y float64, size float64, s string, textcolor color.NRGBA) {
	ff := &text.GoTextFace{Source: mplusFaceSource, Size: size}
	tw := text.Advance(s, ff)
	op := &text.DrawOptions{}
	op.GeoM.Translate(x-tw, y-size)
	op.ColorScale.ScaleWithColor(textcolor)
	text.Draw(screen, s, ff, op)
}

// rtext draws rotated text (angle theta (radians)), starting at (x,y)
func rtext(screen *ebiten.Image, x, y float64, size, theta float64, s string, textcolor color.NRGBA) {
	ff := &text.GoTextFace{Source: mplusFaceSource, Size: size}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y-size)
	op.GeoM.Rotate(theta)
	op.ColorScale.ScaleWithColor(textcolor)
	text.Draw(screen, s, ff, op)
}

// whitespace determines if a rune is whitespace
func whitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

func textwrap(screen *ebiten.Image, x, y, w, size float64, s string, color color.NRGBA) {
	const factor = 0.3
	leading := size * 1.2
	ff := &text.GoTextFace{Source: mplusFaceSource, Size: size}
	wordspacing := text.Advance("M", ff)
	xp := x
	yp := y
	edge := x + w
	words := strings.FieldsFunc(s, whitespace)
	for _, s := range words {
		tw := text.Advance(s, ff)
		btext(screen, xp, yp, size, s, color)
		xp += tw + (wordspacing * factor)
		if xp >= edge {
			xp = x
			yp += leading
		}
	}
}

// centerRect draws a filled rectangle centered at (x,y) with dimensions (w,h)
func centerRect(screen *ebiten.Image, x, y, w, h float32, fillcolor color.NRGBA) {
	px, py := x-(w/2), y-(h/2)
	vector.DrawFilledRect(screen, px, py, w, h, fillcolor, true)
}

// cornerRect draws a filled rectangle with upperleft at (x,y) with dimensions (w,h)
func cornerRect(screen *ebiten.Image, x, y, w, h float32, fillcolor color.NRGBA) {
	px, py := x, y
	vector.DrawFilledRect(screen, px, py, w, h, fillcolor, true)
}

// circle draws a filled circle centered at (x,y), with radius r
func circle(screen *ebiten.Image, cx, cy, r float32, fillcolor color.NRGBA) {
	vector.DrawFilledCircle(screen, cx, cy, r, fillcolor, true)
}

// line draws a line between (x1,y1) and (x2,y2)
func line(screen *ebiten.Image, x1, y1, x2, y2, sw float32, strokecolor color.NRGBA) {
	vector.StrokeLine(screen, x1, y1, x2, y2, sw, strokecolor, true)
}

// polygon draws a filled polygon using the points in x and y
func polygon(screen *ebiten.Image, x, y []float32, fillcolor color.NRGBA) {
	l := len(x)
	if l != len(y) {
		return
	}
	if l < 3 {
		return
	}
	var p vector.Path
	p.MoveTo(x[0], y[0])
	for i := 1; i < l; i++ {
		p.LineTo(x[i], y[i])
	}
	vector.DrawFilledPath(screen, &p, fillcolor, true, vector.FillRuleNonZero)
}

// quadcurve draws a filled quadradic bezier curve beginning at (x1,y1),
// with control point at (x2,y2), ending at (x3,y3)
func quadcurve(screen *ebiten.Image, x1, y1, x2, y2, x3, y3 float32, fillcolor color.NRGBA) {
	var p vector.Path
	p.MoveTo(x1, y1)
	p.QuadTo(x2, y2, x3, y3)
	vector.DrawFilledPath(screen, &p, fillcolor, true, vector.FillRuleEvenOdd)
}

// strokeduadcurve strokes a quadradic bezier curve beginning at (x1,y1),
// with control point at (x2,y2), ending at (x3,y3)
func strokedquadcurve(screen *ebiten.Image, x1, y1, x2, y2, x3, y3, sw float32, strokecolor color.NRGBA) {
	var p vector.Path
	p.MoveTo(x1, y1)
	p.QuadTo(x2, y2, x3, y3)
	op := vector.StrokeOptions{Width: sw}
	vector.StrokePath(screen, &p, strokecolor, true, &op)
}

// cubecurve makes a filled cubic Bezier curve beginning at (x1, y1),
// control points at (x2,y2) and (x3,y3), ending at (x4,y4)
func cubecurve(screen *ebiten.Image, x1, y1, x2, y2, x3, y3, x4, y4 float32, fillcolor color.NRGBA) {
	var p vector.Path
	p.MoveTo(x1, y1)
	p.CubicTo(x2, y2, x3, y3, x4, y4)
	vector.DrawFilledPath(screen, &p, fillcolor, true, vector.FillRuleEvenOdd)
}

// strokedcubecurve strokes a cubic Bezier curve beginning at (x1, y1),
// control points at (x2,y2) and (x3,y3), ending at (x4,y4)
func strokedcubecurve(screen *ebiten.Image, x1, y1, x2, y2, x3, y3, x4, y4, sw float32, strokecolor color.NRGBA) {
	var p vector.Path
	p.MoveTo(x1, y1)
	p.CubicTo(x2, y2, x3, y3, x4, y4)
	op := vector.StrokeOptions{Width: sw}
	vector.StrokePath(screen, &p, strokecolor, true, &op)
}

// showimage places an image with the upper left corner at (x,y)
func showimage(screen *ebiten.Image, x, y float32, scale float64, img image.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	im := ebiten.NewImageFromImage(img)
	screen.DrawImage(im, op)
}

// Percentage based methods: (x, y and measures range from 0-100%),
// The coordinate system is the classical orientation
// (x increasing left to right, y increasing bottom to top)
// bottom left corner is (0,0), top right corner is (100,100), middle of the canvas is (50,50)
// Measures such as text size as scaled by the canvas width.

// Image methods

// CenterImage places an image centered at (x,y) at the specified scale (0-100)
// using percent-based coordinates and measures
func (c *Canvas) CenterImage(x, y float32, scale float32, img image.Image) {
	cw, ch := float32(c.Width), float32(c.Height)
	x, y = dimen(x, y, cw, ch)
	scale /= 100
	fimw, fimh := float32(img.Bounds().Max.X)*scale, float32(img.Bounds().Max.Y)*scale
	showimage(c.Screen, x-(fimw/2), y-(fimh/2), float64(scale), img)
}

// CornerImage places an image with the upper left corner at (x,y) t the specified scale (0-100)
// using percent-based coordinates and measures
func (c *Canvas) CornerImage(x, y float32, scale float64, img image.Image) {
	cw, ch := float32(c.Width), float32(c.Height)
	x, y = dimen(x, y, cw, ch)
	showimage(c.Screen, x, y, scale/100, img)
}

// Image places an image centered at (x,y) (shorthand for CenterImage)
// using percent-based coordinates and measures
func (c *Canvas) Image(x, y float32, scale float32, img image.Image) {
	c.CenterImage(x, y, scale, img)
}

// Shape methods

// Arc draws an filled arc centered at (cx,cy) with radius r,
// using percent-based coordinates and measures,
// between angles a1 and a2 (0-360 degrees, counter-clockwise)
func (c *Canvas) Arc(cx, cy, r, a1, a2 float32, fillcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy = dimen(cx, cy, cw, ch)
	r = pct(r, cw)
	a1 = degreesToRadians(a1)
	a2 = degreesToRadians(a2)
	arc(c.Screen, cx, cy, r, a1, a2, fillcolor)
}

// StrokedArc draws an stroked arc centered at (cx,cy) with radius r,
// using percent-based coordinates and measures,
// between angles a1 and a2 (0-360 degrees, counter-clockwise)
func (c *Canvas) StrokedArc(cx, cy, r, a1, a2, size float32, strokecolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy = dimen(cx, cy, cw, ch)
	r = pct(r, cw)
	size = pct(size, cw)
	a1 = degreesToRadians(a1)
	a2 = degreesToRadians(a2)
	strokedarc(c.Screen, cx, cy, r, a1, a2, size, strokecolor)
}

// Wedge fills a wedge centered at (x,y), with radius r, using percentage-based
// coordinates and measures.  The wedge is filled with fillcolor between
// angles a1 and a2 (degrees between 0-360)
func (c *Canvas) Wedge(cx, cy, r, a1, a2 float32, fillcolor color.NRGBA) {
	px := make([]float32, 3)
	py := make([]float32, 3)
	px[0], py[0] = cx, cy
	px[1], py[1] = c.PolarDegrees(cx, cy, r, a1)
	px[2], py[2] = c.PolarDegrees(cx, cy, r, a2)
	c.Polygon(px, py, fillcolor)
	c.Arc(cx, cy, r, a1, a2, fillcolor)
}

// CenterRect draws a filled rectangle centered at (x,y) with dimensions (w,h)
// using percent-based coordinates and measures
func (c *Canvas) CenterRect(x, y, w, h float32, fillcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	w = pct(w, cw)
	h = pct(h, ch)
	x, y = dimen(x, y, cw, ch)
	centerRect(c.Screen, x, y, w, h, fillcolor)
}

// CornerRect draws a filled rectangle centered at (x,y) with dimensions (w,h)
// using percent-based coordinates and measures
func (c *Canvas) CornerRect(x, y, w, h float32, fillcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	w = pct(w, cw)
	h = pct(h, ch)
	x, y = dimen(x, y, cw, ch)
	cornerRect(c.Screen, x, y, w, h, fillcolor)
}

// Rect draws a filled rectangle centered at (x,y) with dimensions (w,h)
// using percent-based coordinates and measures
func (c *Canvas) Rect(x, y, w, h float32, fillcolor color.NRGBA) {
	c.CenterRect(x, y, w, h, fillcolor)
}

// Circle draws a filled circle centered at (x,y), with radius r
// using percent-based coordinates and measures
func (c *Canvas) Circle(cx, cy, r float32, fillcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy = dimen(cx, cy, cw, ch)
	r = pct(r, cw)
	circle(c.Screen, cx, cy, r, fillcolor)
}

// HLine makes a horizonal line beginning at (x,y) extending to the right
func (c *Canvas) HLine(x1, y1, size, sw float32, strokecolor color.NRGBA) {
	c.Line(x1, y1, x1+size, y1, sw, strokecolor)
}

// Line draws a line between (x1,y1) and (x2,y2)
// using percent-based coordinates and measures
func (c *Canvas) Line(x1, y1, x2, y2, sw float32, strokecolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	x1, y1 = dimen(x1, y1, cw, ch)
	x2, y2 = dimen(x2, y2, cw, ch)
	sw = pct(sw, cw)
	line(c.Screen, x1, y1, x2, y2, sw, strokecolor)
}

// Polygon draws a filled polygon using the points in x and y
// using percent-based coordinates and measures
func (c *Canvas) Polygon(x, y []float32, fillcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	for i := 0; i < len(x); i++ {
		x[i], y[i] = dimen(x[i], y[i], cw, ch)
	}
	polygon(c.Screen, x, y, fillcolor)
}

// StrokedPolygon strokes a polygon of the specified size and color, using the points in x and y,
// using percent-based coordinates and measures
func (c *Canvas) StrokedPolygon(x, y []float32, size float32, strokecolor color.NRGBA) {
	l := len(x) - 1
	for i := 0; i < l; i++ {
		c.Line(x[i], y[i], x[i+1], y[i+1], size, strokecolor)
	}
	c.Line(x[l], y[l], x[0], y[0], size, strokecolor)
}

// Quadcurve draws a filled quadradic bezier curve beginning at (x1,y1),
// with control point at (x2,y2). ending at (x3,y3),
// using percent-based coordinates and measures
func (c *Canvas) QuadCurve(x1, y1, x2, y2, x3, y3 float32, fillcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	x1, y1 = dimen(x1, y1, cw, ch)
	x2, y2 = dimen(x2, y2, cw, ch)
	x3, y3 = dimen(x3, y3, cw, ch)
	quadcurve(c.Screen, x1, y1, x2, y2, x3, y3, fillcolor)
}

// QuadStrokedCurve strokes a quadradic bezier curve beginning at (x1,y1),
// with control point at (x2,y2). ending at (x3,y3),
// using percent-based coordinates and measures
func (c *Canvas) QuadStrokedCurve(x1, y1, x2, y2, x3, y3, size float32, strokecolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	x1, y1 = dimen(x1, y1, cw, ch)
	x2, y2 = dimen(x2, y2, cw, ch)
	x3, y3 = dimen(x3, y3, cw, ch)
	size = pct(size, cw)
	strokedquadcurve(c.Screen, x1, y1, x2, y2, x3, y3, size, strokecolor)
}

// CubeCurve makes a filled cubic Bezier curve beginning at (x1,y1),
// with control points at (x2,y2). ending at (x3,y3), ending at (x4,y4)
// using percent-based coordinates and measures
func (c *Canvas) CubeCurve(x1, y1, x2, y2, x3, y3, x4, y4 float32, strokecolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	x1, y1 = dimen(x1, y1, cw, ch)
	x2, y2 = dimen(x2, y2, cw, ch)
	x3, y3 = dimen(x3, y3, cw, ch)
	x4, y4 = dimen(x4, y4, cw, ch)
	cubecurve(c.Screen, x1, y1, x2, y2, x3, y3, x4, y4, strokecolor)
}

// CubeCurve strokes a cubic Bezier curve beginning at (x1,y1),
// with control points at (x2,y2). ending at (x3,y3), ending at (x4,y4)
// using percent-based coordinates and measures
func (c *Canvas) StrokedCubeCurve(x1, y1, x2, y2, x3, y3, x4, y4, size float32, strokecolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	x1, y1 = dimen(x1, y1, cw, ch)
	x2, y2 = dimen(x2, y2, cw, ch)
	x3, y3 = dimen(x3, y3, cw, ch)
	x4, y4 = dimen(x4, y4, cw, ch)
	size = pct(size, cw)
	strokedcubecurve(c.Screen, x1, y1, x2, y2, x3, y3, x4, y4, size, strokecolor)
}

// Curve is a shorthand for QuadCurve
func (c *Canvas) Curve(x1, y1, x2, y2, x3, y3 float32, fillcolor color.NRGBA) {
	c.QuadCurve(x1, y1, x2, y2, x3, y3, fillcolor)
}

// StrokedCurve is a shorthand for QuadStrokeCurve
func (c *Canvas) StrokedCurve(x1, y1, x2, y2, x3, y3, size float32, strokecolor color.NRGBA) {
	c.QuadStrokedCurve(x1, y1, x2, y2, x3, y3, size, strokecolor)
}

// Square draws a filled square centered at (x,y), sides at size
func (c *Canvas) Square(x, y, w float32, fillcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	x, y = dimen(x, y, cw, ch)
	w = pct(w, ch)
	h := pct(100, w)
	centerRect(c.Screen, x, y, w, h, fillcolor)
}

// VLines draws a vertical line begging ar (x,y) moving up size
func (c *Canvas) VLine(x1, y1, size, sw float32, strokecolor color.NRGBA) {
	c.Line(x1, y1, x1, y1+size, sw, strokecolor)
}

// Text methods

// Text draws text contained in s, beginning at (x,y), at the specifed size
// using percent-based coordinates and measures
func (c *Canvas) Text(x, y, size float32, s string, textcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy := dimen(x, y, cw, ch)
	size = pct(size, cw)
	btext(c.Screen, float64(cx), float64(cy), float64(size), s, textcolor)
}

// CText draws text contained in s centered at (x,y), at the specified size
// using percent-based coordinates and measures
func (c *Canvas) CText(x, y, size float32, s string, textcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy := dimen(x, y, cw, ch)
	size = pct(size, cw)
	ctext(c.Screen, float64(cx), float64(cy), float64(size), s, textcolor)
}

// TextMid is an alternative name for CText
func (c *Canvas) TextMid(x, y, size float32, s string, textcolor color.NRGBA) {
	c.CText(x, y, size, s, textcolor)
}

// EText draws text contained in s with end point at (x,y) at the specified size
// using percent-based coordinates and measures
func (c *Canvas) EText(x, y, size float32, s string, textcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy := dimen(x, y, cw, ch)
	size = pct(size, cw)
	etext(c.Screen, float64(cx), float64(cy), float64(size), s, textcolor)
}

// TextEnd is an alternative name for EText
func (c *Canvas) TextEnd(x, y, size float32, s string, textcolor color.NRGBA) {
	c.EText(x, y, size, s, textcolor)
}

// RText draws rotated text at (x,y), rotated at the specified angle
func (c *Canvas) RText(x, y, angle, size float32, s string, textcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy := dimen(x, y, cw, ch)
	size = pct(size, cw)
	theta := float64(angle) * (Pi / 180)
	rtext(c.Screen, float64(cx), float64(cy), theta, float64(size), s, textcolor)
}

func (c *Canvas) TextWrap(x, y, w, size float32, s string, textcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy := dimen(x, y, cw, ch)
	size = pct(size, cw)
	w = pct(w, cw)
	textwrap(c.Screen, float64(cx), float64(cy), float64(w), float64(size), s, textcolor)
}

// Utility Methods

// Background fills the canvas with the specified color
func (c *Canvas) Background(fillcolor color.NRGBA) {
	c.Screen.Fill(fillcolor)
}

// Grid draws a grid starting at (x,y), dimensions at (w,h),
// A gridline is drawn at the specifed interval.
func (c *Canvas) Grid(x, y, w, h, size, interval float32, strokecolor color.NRGBA) {
	for xp := x; xp <= x+w; xp += interval {
		c.Line(xp, y, xp, y+h, size, strokecolor) // vertical line
	}
	for yp := y; yp <= y+h; yp += interval {
		c.Line(x, yp, x+w, yp, size, strokecolor) // horizontal line
	}
}

// PolarDegrees returns the Cartesian coordinates (x, y) from polar coordinates
// with compensation for canvas aspect ratio
// center at (cx, cy), radius r, and angle theta (degrees)
func (c *Canvas) PolarDegrees(cx, cy, r, theta float32) (float32, float32) {
	fr := float64(r)
	ft := float64(theta * (math.Pi / 180))
	aspect := float64(c.Width / c.Height)
	px := fr * math.Cos(ft)
	py := (fr * aspect) * math.Sin(ft)
	return cx + float32(px), cy + float32(py)
}

// Polar returns the Cartesian coordinates (x, y) from polar coordinates
// with compensation for canvas aspect ratio
// center at (cx, cy), radius r, and angle theta (radians)
func (c *Canvas) Polar(cx, cy, r, theta float32) (float32, float32) {
	fr := float64(r)
	ft := float64(theta)
	aspect := float64(c.Width / c.Height)
	px := fr * math.Cos(ft)
	py := (fr * aspect) * math.Sin(ft)
	return cx + float32(px), cy + float32(py)
}

// Coord shows the specified coordinate, using percentage-based coordinates
// the (x, y) label is above the point, with a label below
func (c *Canvas) Coord(x, y, size float32, s string, fillcolor color.NRGBA) {
	c.Circle(x, y, size/4, fillcolor)
	b := []byte("(")
	b = strconv.AppendFloat(b, float64(x), 'g', -1, 32)
	b = append(b, ',')
	b = strconv.AppendFloat(b, float64(y), 'g', -1, 32)
	b = append(b, ')')
	c.CText(x, y+size, size, string(b), fillcolor)
	if len(s) > 0 {
		c.CText(x, y-(size*1.33), size*0.66, s, fillcolor)
	}
}

// MapRange maps a value between low1 and high1, return the corresponding value between low2 and high2
func MapRange(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}
