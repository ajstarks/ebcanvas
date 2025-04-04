package ebcanvas

import (
	"image/color"
	"math"
	"strconv"
	"strings"
)

// colornames maps SVG color names to RGB triples.
var colornames = map[string]color.NRGBA{
	"aliceblue":            {240, 248, 255, 255},
	"antiquewhite":         {250, 235, 215, 255},
	"aqua":                 {0, 255, 255, 255},
	"aquamarine":           {127, 255, 212, 255},
	"azure":                {240, 255, 255, 255},
	"beige":                {245, 245, 220, 255},
	"bisque":               {255, 228, 196, 255},
	"black":                {0, 0, 0, 255},
	"blanchedalmond":       {255, 235, 205, 255},
	"blue":                 {0, 0, 255, 255},
	"blueviolet":           {138, 43, 226, 255},
	"brown":                {165, 42, 42, 255},
	"burlywood":            {222, 184, 135, 255},
	"cadetblue":            {95, 158, 160, 255},
	"chartreuse":           {127, 255, 0, 255},
	"chocolate":            {210, 105, 30, 255},
	"coral":                {255, 127, 80, 255},
	"cornflowerblue":       {100, 149, 237, 255},
	"cornsilk":             {255, 248, 220, 255},
	"crimson":              {220, 20, 60, 255},
	"cyan":                 {0, 255, 255, 255},
	"darkblue":             {0, 0, 139, 255},
	"darkcyan":             {0, 139, 139, 255},
	"darkgoldenrod":        {184, 134, 11, 255},
	"darkgray":             {169, 169, 169, 255},
	"darkgreen":            {0, 100, 0, 255},
	"darkgrey":             {169, 169, 169, 255},
	"darkkhaki":            {189, 183, 107, 255},
	"darkmagenta":          {139, 0, 139, 255},
	"darkolivegreen":       {85, 107, 47, 255},
	"darkorange":           {255, 140, 0, 255},
	"darkorchid":           {153, 50, 204, 255},
	"darkred":              {139, 0, 0, 255},
	"darksalmon":           {233, 150, 122, 255},
	"darkseagreen":         {143, 188, 143, 255},
	"darkslateblue":        {72, 61, 139, 255},
	"darkslategray":        {47, 79, 79, 255},
	"darkslategrey":        {47, 79, 79, 255},
	"darkturquoise":        {0, 206, 209, 255},
	"darkviolet":           {148, 0, 211, 255},
	"deeppink":             {255, 20, 147, 255},
	"deepskyblue":          {0, 191, 255, 255},
	"dimgray":              {105, 105, 105, 255},
	"dimgrey":              {105, 105, 105, 255},
	"dodgerblue":           {30, 144, 255, 255},
	"firebrick":            {178, 34, 34, 255},
	"floralwhite":          {255, 250, 240, 255},
	"forestgreen":          {34, 139, 34, 255},
	"fuchsia":              {255, 0, 255, 255},
	"gainsboro":            {220, 220, 220, 255},
	"ghostwhite":           {248, 248, 255, 255},
	"gold":                 {255, 215, 0, 255},
	"goldenrod":            {218, 165, 32, 255},
	"gray":                 {128, 128, 128, 255},
	"green":                {0, 128, 0, 255},
	"greenyellow":          {173, 255, 47, 255},
	"grey":                 {128, 128, 128, 255},
	"honeydew":             {240, 255, 240, 255},
	"hotpink":              {255, 105, 180, 255},
	"indianred":            {205, 92, 92, 255},
	"indigo":               {75, 0, 130, 255},
	"ivory":                {255, 255, 240, 255},
	"khaki":                {240, 230, 140, 255},
	"lavender":             {230, 230, 250, 255},
	"lavenderblush":        {255, 240, 245, 255},
	"lawngreen":            {124, 252, 0, 255},
	"lemonchiffon":         {255, 250, 205, 255},
	"lightblue":            {173, 216, 230, 255},
	"lightcoral":           {240, 128, 128, 255},
	"lightcyan":            {224, 255, 255, 255},
	"lightgoldenrodyellow": {250, 250, 210, 255},
	"lightgray":            {211, 211, 211, 255},
	"lightgreen":           {144, 238, 144, 255},
	"lightgrey":            {211, 211, 211, 255},
	"lightpink":            {255, 182, 193, 255},
	"lightsalmon":          {255, 160, 122, 255},
	"lightseagreen":        {32, 178, 170, 255},
	"lightskyblue":         {135, 206, 250, 255},
	"lightslategray":       {119, 136, 153, 255},
	"lightslategrey":       {119, 136, 153, 255},
	"lightsteelblue":       {176, 196, 222, 255},
	"lightyellow":          {255, 255, 224, 255},
	"lime":                 {0, 255, 0, 255},
	"limegreen":            {50, 205, 50, 255},
	"linen":                {250, 240, 230, 255},
	"magenta":              {255, 0, 255, 255},
	"maroon":               {128, 0, 0, 255},
	"mediumaquamarine":     {102, 205, 170, 255},
	"mediumblue":           {0, 0, 205, 255},
	"mediumorchid":         {186, 85, 211, 255},
	"mediumpurple":         {147, 112, 219, 255},
	"mediumseagreen":       {60, 179, 113, 255},
	"mediumslateblue":      {123, 104, 238, 255},
	"mediumspringgreen":    {0, 250, 154, 255},
	"mediumturquoise":      {72, 209, 204, 255},
	"mediumvioletred":      {199, 21, 133, 255},
	"midnightblue":         {25, 25, 112, 255},
	"mintcream":            {245, 255, 250, 255},
	"mistyrose":            {255, 228, 225, 255},
	"moccasin":             {255, 228, 181, 255},
	"navajowhite":          {255, 222, 173, 255},
	"navy":                 {0, 0, 128, 255},
	"oldlace":              {253, 245, 230, 255},
	"olive":                {128, 128, 0, 255},
	"olivedrab":            {107, 142, 35, 255},
	"orange":               {255, 165, 0, 255},
	"orangered":            {255, 69, 0, 255},
	"orchid":               {218, 112, 214, 255},
	"palegoldenrod":        {238, 232, 170, 255},
	"palegreen":            {152, 251, 152, 255},
	"paleturquoise":        {175, 238, 238, 255},
	"palevioletred":        {219, 112, 147, 255},
	"papayawhip":           {255, 239, 213, 255},
	"peachpuff":            {255, 218, 185, 255},
	"peru":                 {205, 133, 63, 255},
	"pink":                 {255, 192, 203, 255},
	"plum":                 {221, 160, 221, 255},
	"powderblue":           {176, 224, 230, 255},
	"purple":               {128, 0, 128, 255},
	"red":                  {255, 0, 0, 255},
	"rosybrown":            {188, 143, 143, 255},
	"royalblue":            {65, 105, 225, 255},
	"saddlebrown":          {139, 69, 19, 255},
	"salmon":               {250, 128, 114, 255},
	"sandybrown":           {244, 164, 96, 255},
	"seagreen":             {46, 139, 87, 255},
	"seashell":             {255, 245, 238, 255},
	"sienna":               {160, 82, 45, 255},
	"silver":               {192, 192, 192, 255},
	"skyblue":              {135, 206, 235, 255},
	"slateblue":            {106, 90, 205, 255},
	"slategray":            {112, 128, 144, 255},
	"slategrey":            {112, 128, 144, 255},
	"snow":                 {255, 250, 250, 255},
	"springgreen":          {0, 255, 127, 255},
	"steelblue":            {70, 130, 180, 255},
	"tan":                  {210, 180, 140, 255},
	"teal":                 {0, 128, 128, 255},
	"thistle":              {216, 191, 216, 255},
	"tomato":               {255, 99, 71, 255},
	"turquoise":            {64, 224, 208, 255},
	"violet":               {238, 130, 238, 255},
	"wheat":                {245, 222, 179, 255},
	"white":                {255, 255, 255, 255},
	"whitesmoke":           {245, 245, 245, 255},
	"yellow":               {255, 255, 0, 255},
	"yellowgreen":          {154, 205, 50, 255},
}

// cc converts a color string to number
func cc(s string) uint8 {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return uint8(v)
}

// hc converts a hex color string to number
func hc(s string) uint8 {
	v, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		return 0
	}
	return uint8(v)
}

// ColorLookup returns a color.NRGBA corresponding to the named color or
// "rgb(r)", "rgb(r,b)", "rgb(r,g,b), "rgb(r,g,b,a)",
// "#rr",     "#rrgg",   "#rrggbb",   "#rrggbbaa" string.
// "hsv(hue,sat,value)"
// On error, return black.
func ColorLookup(s string) color.NRGBA {
	c, ok := colornames[s]
	if ok {
		return c
	}
	black := color.NRGBA{0, 0, 0, 255}
	ls := len(s)
	// rgb(...)
	if strings.HasPrefix(s, "rgb(") && strings.HasSuffix(s, ")") && ls > 5 {
		c.R, c.G, c.B, c.A = 0, 0, 0, 255
		v := strings.Split(s[4:ls-1], ",")
		switch len(v) {
		case 1:
			c.R = cc(v[0])
			return c
		case 2:
			c.R = cc(v[0])
			c.G = cc(v[1])
			return c
		case 3:
			c.R = cc(v[0])
			c.G = cc(v[1])
			c.B = cc(v[2])
			return c
		case 4:
			c.R = cc(v[0])
			c.G = cc(v[1])
			c.B = cc(v[2])
			c.A = cc(v[3])
			return c
		default:
			return black
		}
	}
	// hsv(h,s,v) or hsv(h,s,v,a); h=0-360, s, v=0-100, a=0-100
	if strings.HasPrefix(s, "hsv(") && strings.HasSuffix(s, ")") && ls > 5 {
		v := strings.Split(s[4:ls-1], ",")
		switch len(v) {
		case 3:
			hue, _ := strconv.ParseFloat(v[0], 64)
			sat, _ := strconv.ParseFloat(v[1], 64)
			value, _ := strconv.ParseFloat(v[2], 64)
			c.R, c.G, c.B = hsv2rgb(hue, sat, value)
			c.A = 255
			return c
		case 4:
			hue, _ := strconv.ParseFloat(v[0], 64)
			sat, _ := strconv.ParseFloat(v[1], 64)
			value, _ := strconv.ParseFloat(v[2], 64)
			a := cc(v[3])
			if a > 100 {
				a = 100
			}
			c.R, c.G, c.B = hsv2rgb(hue, sat, value)
			c.A = uint8((float64(a) / 100.0) * 255.0)
			return c
		default:
			return black
		}
	}
	// #rrggbb
	if strings.HasPrefix(s, "#") && (ls >= 3) {
		c.R, c.G, c.B, c.A = 0, 0, 0, 255
		switch ls {
		case 3:
			c.R = hc(s[1:3])
		case 5:
			c.R = hc(s[1:3])
			c.G = hc(s[3:5])
		case 7:
			c.R = hc(s[1:3])
			c.G = hc(s[3:5])
			c.B = hc(s[5:7])
		case 9:
			c.R = hc(s[1:3])
			c.G = hc(s[3:5])
			c.B = hc(s[5:7])
			c.A = hc(s[7:9])
		default:
			return black
		}
		return c
	}
	return black
}

// hsv2rgb converts hsv(h (0-360), s (0-100), v (0-100)) to rgb
// reference: https://en.wikipedia.org/wiki/HSL_and_HSV#HSV_to_RGB
func hsv2rgb(h, s, v float64) (uint8, uint8, uint8) {
	s /= 100
	v /= 100
	if s > 1 || v > 1 {
		return 0, 0, 0
	}
	h = math.Mod(h, 360)
	c := v * s
	section := h / 60
	x := c * (1 - math.Abs(math.Mod(section, 2)-1))

	var r, g, b float64
	switch {
	case section >= 0 && section <= 1:
		r = c
		g = x
		b = 0
	case section > 1 && section <= 2:
		r = x
		g = c
		b = 0
	case section > 2 && section <= 3:
		r = 0
		g = c
		b = x
	case section > 3 && section <= 4:
		r = 0
		g = x
		b = c
	case section > 4 && section <= 5:
		r = x
		g = 0
		b = c
	case section > 5 && section <= 6:
		r = c
		g = 0
		b = x
	default:
		return 0, 0, 0
	}
	m := v - c
	r += m
	g += m
	b += m
	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}
