package ebcanvas

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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

// strokedarc strokes an arc centered at (cx,cy) with radius r, between angle a1 and a2
func strokedarc(screen *ebiten.Image, cx, cy, r, a1, a2, size float32, strokecolor color.NRGBA) {
	var p vector.Path
	op := vector.StrokeOptions{Width: size}
	p.Arc(cx, cy, r, a1, a2, vector.CounterClockwise)
	vector.StrokePath(screen, &p, strokecolor, true, &op)
}

// btext draws text beginning at (x,y)
func btext(screen *ebiten.Image, x, y float64, size float64, s string, textcolor color.NRGBA) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y-size)
	op.ColorScale.ScaleWithColor(textcolor)
	text.Draw(screen, s, &text.GoTextFace{Source: mplusFaceSource, Size: size}, op)
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

// centerRect draws a filled rectangle centered at (x,y) with dimensions (w,h)
func centerRect(screen *ebiten.Image, x, y, w, h float32, fillcolor color.NRGBA) {
	px, py := x-(w/2), y-(h/2)
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
	lx := len(x)
	if lx != len(y) {
		return
	}
	if lx < 3 {
		return
	}
	var p vector.Path
	p.MoveTo(x[0], y[0])
	for i := 1; i < lx; i++ {
		p.LineTo(x[i], y[i])
	}
	p.Close()
	vector.DrawFilledPath(screen, &p, fillcolor, true, vector.FillRuleEvenOdd)
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

func cubecurve(screen *ebiten.Image, x1, y1, x2, y2, x3, y3, x4, y4 float32, fillcolor color.NRGBA) {
	var p vector.Path
	p.MoveTo(x1, y1)
	p.CubicTo(x2, y2, x3, y3, x4, y4)
	vector.DrawFilledPath(screen, &p, fillcolor, true, vector.FillRuleEvenOdd)
}

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

// Arc draws an filled arc centered at (cx,cy) with radius r, between angle a1 and a2
// using percent-based coordinates and measures
func (c *Canvas) Arc(cx, cy, r, a1, a2 float32, fillcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy = dimen(cx, cy, cw, ch)
	r = pct(r, cw)
	arc(c.Screen, cx, cy, r, a1, a2, fillcolor)
}

// StrokedArc draws an stroked arc centered at (cx,cy) with radius r, between angle a1 and a2
// using percent-based coordinates and measures
func (c *Canvas) StrokedArc(cx, cy, r, a1, a2, size float32, strokecolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy = dimen(cx, cy, cw, ch)
	r = pct(r, cw)
	size = pct(size, cw)
	strokedarc(c.Screen, cx, cy, r, a1, a2, size, strokecolor)
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

// Text draws text beginning at (x,y)
// using percent-based coordinates and measures
func (c *Canvas) Text(x, y, size float32, s string, textcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy := dimen(x, y, cw, ch)
	size = pct(size, cw)
	btext(c.Screen, float64(cx), float64(cy), float64(size), s, textcolor)
}

// CText draws text centered at (x,y)
// using percent-based coordinates and measures
func (c *Canvas) CText(x, y, size float32, s string, textcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy := dimen(x, y, cw, ch)
	size = pct(size, cw)
	ctext(c.Screen, float64(cx), float64(cy), float64(size), s, textcolor)
}

// EText draws text with end point at (x,y)
// using percent-based coordinates and measures
func (c *Canvas) EText(x, y, size float32, s string, textcolor color.NRGBA) {
	cw, ch := float32(c.Width), float32(c.Height)
	cx, cy := dimen(x, y, cw, ch)
	size = pct(size, cw)
	etext(c.Screen, float64(cx), float64(cy), float64(size), s, textcolor)
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

// MapRange maps a value between low1 and high1, return the corresponding value between low2 and high2
func MapRange(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
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
