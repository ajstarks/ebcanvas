package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/ajstarks/ebcanvas"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	screenWidth  int = 1000
	screenHeight int = 1000
	earth        image.Image
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	work(screen)
}

func work(screen *ebiten.Image) {
	canvas := new(ebcanvas.Canvas)
	canvas.Screen = screen
	canvas.Width = screenWidth
	canvas.Height = screenHeight

	black := color.NRGBA{0, 0, 0, 255}
	white := color.NRGBA{255, 255, 255, 255}
	gray := color.NRGBA{200, 200, 200, 255}
	red := color.NRGBA{255, 0, 0, 255}
	green := color.NRGBA{0, 128, 0, 255}
	magenta := color.NRGBA{255, 0, 255, 255}
	yellow := color.NRGBA{255, 255, 0, 255}
	orange := color.NRGBA{255, 165, 0, 255}
	maroon := color.NRGBA{128, 0, 0, 255}

	bgcolor := white
	txcolor := black
	descolor := maroon
	if bgcolor == black {
		descolor = gray
	}

	const (
		top       float32 = 95.0
		fx        float32 = 5.0
		textsize  float32 = 2.25
		dotsize   float32 = 0.5
		vspace    float32 = 9.0
		hspace    float32 = 10.0
		halfhs    float32 = hspace / 2
		halfvs    float32 = vspace / 2
		objx      float32 = 85.0
		linewidth float32 = dotsize / 2
		pi                = 3.14159265359
	)

	// Begin...
	canvas.Background(bgcolor)
	canvas.Grid(0, 0, 100, 100, 0.1, 5, gray)

	// Title
	canvas.Text(fx, top, textsize*1.5, "Ebiten Canvas API", txcolor)
	yp := top - halfvs

	// API labels
	funcnames := []string{
		"{R,C,E}Text(x, y, {angle}, size float32, s string, c color.NRGBA)",
		"TextWrap(x, y, w, size float32, s string, c color.NRGBA)",
		"Circle(x, y, r float32, c color.NRGBA)",
		"Rect(x, y, w, h float32, c color.NRGBA)",
		"Square(x, y, w float32, c color.NRGBA)",
		"{Stroked}Arc(cx, cy, r, a1, a2, {,sw} float32, c color.NRGBA)",
		"{Stroked}Curve(bx,by, cx,cy, ex,ey {,sw} float32, c color.NRGBA)",
		`Line(x1,y1, x2,y2, sw float32, c color.NRGBA)`,
		"{Stroked}Polygon(x,y {,sw} float32, c color.NRGBA)",
		"{Corner}Image(x,y, scale float32,img image.Image)",
	}
	funcdesc := []string{
		"Rorated, Centered, End, Start aligned text at (x,y)",
		"Wrap text at (x,y) to width w",
		"Circle at (x,y), radius r",
		"Rectangle centered at (x,y), dimensions (w,h)",
		"Square centered at (x,y), size w",
		"Stroked or filled arc centered at (x,y), radius r, between angles a1,a2 (deg)",
		"Stroked or filled Quadradic Bézier curve; begin: (bx,by), control: (cx,cy), end: (ex,ey)",
		"Line between (x1,y1) and (x2,y2). Horizontal or vertical: (x1,y1), length",
		"Stroked or filled Polygon using points in (x,y)",
		"Image anchored at the top left corner or center at (x,y)",
	}
	for i, ly := 0, yp; i < len(funcnames); i++ {
		canvas.Text(fx, ly, textsize, funcnames[i], txcolor)
		canvas.Text(fx, ly-(textsize*1.2), textsize*0.75, funcdesc[i], descolor)
		ly -= vspace
	}

	// Text
	message := "hello"
	wmessage := "This is text wrapped at a specified width"
	labelx := objx - halfhs
	canvas.CText(labelx, yp, textsize, message, txcolor)
	canvas.Circle(labelx, yp, dotsize, red)
	labelx += hspace
	canvas.EText(labelx, yp, textsize, message, txcolor)
	canvas.Circle(labelx, yp, dotsize, red)
	labelx += hspace * 0.4
	canvas.Text(labelx, yp, textsize, message, txcolor)
	canvas.Circle(objx, yp+halfvs, dotsize, red)
	canvas.RText(objx, yp+halfvs, 45, textsize, message, txcolor)
	yp -= vspace
	canvas.TextWrap(objx-hspace, yp+(halfvs/3), 20, textsize*0.8, wmessage, txcolor)

	// Circle
	yp -= vspace
	canvas.Circle(objx, yp, halfhs, red)
	canvas.Circle(objx, yp, dotsize, white)

	// Rect
	yp -= vspace
	canvas.Rect(objx, yp, 20, halfvs, green)
	canvas.Circle(objx, yp, dotsize, red)

	// Square
	yp -= vspace
	canvas.Square(objx, yp, halfhs, green)
	canvas.Circle(objx, yp, dotsize, red)

	// Arc
	yp -= vspace
	canvas.Arc(objx, yp, halfhs, 0, 180, yellow)
	canvas.StrokedArc(objx, yp, halfhs, 0, 180, linewidth, red)
	canvas.Circle(objx, yp, dotsize, red)

	// Curve
	yp -= vspace
	curvex := []float32{objx - halfhs, objx + halfhs, objx + halfhs}
	curvey := []float32{yp, yp + halfvs, yp}
	canvas.Curve(curvex[0], curvey[0], curvex[1], curvey[1], curvex[2], curvey[2], orange)
	canvas.StrokedCurve(curvex[0], curvey[0], curvex[1], curvey[1], curvex[2], curvey[2], linewidth, red)
	for i := 0; i < len(curvex); i++ {
		canvas.Circle(curvex[i], curvey[i], dotsize, red)
	}

	// Line
	yp -= vspace
	canvas.Line(objx, yp, objx+hspace, yp+halfvs, linewidth, txcolor)
	canvas.Circle(objx, yp, dotsize, red)
	canvas.Circle(objx+hspace, yp+halfvs, dotsize, red)
	canvas.HLine(objx-hspace, yp, halfhs, linewidth, txcolor)
	canvas.VLine(objx-5, yp, halfvs, linewidth, txcolor)

	canvas.Circle(objx-hspace, yp, dotsize, red)
	canvas.Circle(objx-5, yp, dotsize, red)

	// Polygon
	yp -= vspace
	px := []float32{objx, objx - hspace, objx, objx + hspace}
	py := []float32{yp + halfvs, yp + 2, yp, yp + 2}
	for i := 0; i < len(px); i++ {
		canvas.Circle(px[i], py[i], dotsize, red)
	}
	canvas.StrokedPolygon(px, py, linewidth*2, txcolor)
	canvas.Polygon(px, py, magenta)

	// Image
	yp -= vspace
	canvas.CornerImage(objx-hspace, yp+5, 20, earth)
	canvas.Image(objx+hspace, yp, 20, earth)
	canvas.Circle(objx+hspace, yp, dotsize, red)
	canvas.Circle(objx-hspace, yp+5, dotsize, red)
}

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebiten Canvas API")
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
	if err = ebiten.RunGame(&Game{}); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
}
