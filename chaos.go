package superfractal

import (
	. "fmt"
	"math/rand"
	"strings"
	"unicode"

	"github.com/strickyak/canvas"
)

const SKIP = 32

var Builtin = make(map[string]string)

func init() {
	// Sierpinski's Triangle.
	Builtin["s"] = "0.5 0 0 0.5 0 0,0.5 0 0 0.5 0.5 0,0.5 0 0 0.5 0 0.5"

	// Square Gasket.
	Builtin["g"] = "0.4 0 0 0.4 0 0,0.4 0 0 0.4 0.6 0,0.4 0 0 0.4 0.6 0.6,0.4 0 0 0.4 0 0.6"

	// Wedgy Augmented Sierpinski.
	Builtin["t"] = "0.5 0 0 0.5 0 0,0.5 0 0 0.5 0.5 0,0.5 0 0 0.5 0 0.5,0.2 0 0 0.2 0.7 0.7"

	// Fern.
	Builtin["f"] = "@.12 .20 .20 .32 .46 .18,@.62 .32 .7 .22 .52 .10,@.46 .18 .461 .18 .461 0,@.26 1 1 .8 .8 .15"
	Builtin["f"] = "@.12 .20 .20 .32 .50 .18,@.62 .32 .7 .22 .52 .10,@.49 .38 .491 .38 .461 0,@.26 1.1 1 .9 .88 .10"
	Builtin["f"] = "@.12 .20 .20 .32 .50 .18,@.62 .32 .7 .22 .52 .10,@.49 .38 .491 .38 .461 0,@.26 1.18 1 .98 .88 .18"

	// Three Quarter Squares
	Builtin["W"] = "0.5 0 0 0.5 0.00 0.00,0.5 0 0 0.5 0.50 0.00,0.5 0 0 0.5 0.50 0.50"
	Builtin["X"] = "0.5 0 0 0.5 0.00 0.00,0.5 0 0 0.5 0.00 0.50,0.5 0 0 0.5 0.50 0.50"
	Builtin["Y"] = "0.5 0 0 0.5 0.00 0.00,0.5 0 0 0.5 0.50 0.00,0.5 0 0 0.5 0.00 0.50"
	Builtin["Z"] = "0.5 0 0 0.5 0.00 0.50,0.5 0 0 0.5 0.50 0.50,0.5 0 0 0.5 0.50 0.00"

}

type Affine struct {
	A, B, C, D, E, F float64
}

type IFS struct {
	Choices []*Affine
}

// Apply the Affine function once, to floats x,y, and return 2 new floats.
func (o Affine) Apply(x, y float64) (float64, float64) {
	ax := x*o.A + y*o.B + o.E
	ay := x*o.C + y*o.D + o.F
	return ax, ay
}

// ChaosGame steps SKIP+n times, plotting the last n points.
func (o IFS) ChaosGame(targ *canvas.Canvas, n int, c canvas.Color) {
	x, y := 0.5, 0.5
	for i := 0; i < SKIP; i++ {
		x, y = o.Choose().Apply(x, y)
	}
	for i := 0; i < n; i++ {
		x, y = o.Choose().Apply(x, y)
		targ.FSet(x, y, c)
	}
}

// WholeGame
func (o IFS) WholeGame(src *canvas.Canvas, n int) *canvas.Canvas {
	t1 := src.Dup()
	t2 := canvas.NewCanvas(src.Width, src.Height)
	for i := 0; i < n; i++ {
		canvas.Say("WholeStep", i)
		o.WholeStep(t1, t2)
		t1, t2 = t2, t1
		t2.Fill(0, 0, t2.Width, t2.Height, canvas.Black)
	}
	return t1
}

func (o IFS) WholeStep(src *canvas.Canvas, dest *canvas.Canvas) {
	wait := make(chan int)
	for k := 0; k < len(o.Choices); k++ {
		canvas.Say("    go Map", k)
		go o.Choices[k].MapImageTo(src, dest, k, wait)
	}
	for k := 0; k < len(o.Choices); k++ {
		w := <-wait
		canvas.Say("      Waited", w)
	}
}

// MapImageTo maps one picture through an IFS onto another.
func (o Affine) MapImageTo(src *canvas.Canvas, dest *canvas.Canvas, k int, done chan int) {
	for i := 0; i < src.Width; i++ {
		x := float64(i) / src.FWidth
		for j := 0; j < src.Height; j++ {
			r, g, b := src.Get(i, j)
			if r > 0 || g > 0 || b > 0 {
				y := float64(j) / src.FHeight
				x2, y2 := o.Apply(x, y)
				dest.FSet(x2, y2, canvas.RGB(r, g, b))
			}
		}
	}
	done <- k
}

// Choose one member of the IFS.
func (o IFS) Choose() *Affine {
	r := rand.Intn(len(o.Choices))
	return o.Choices[r]
}

// CheckOne checks result of scanf of one item.
func CheckOne(n int, err error) {
	if n != 1 || err != nil {
		panic(Errorf("Bad Number: %s", err))
	}
}

// ParseIfsParams parses list of matrices, or the name of a Builtin.
func ParseIfsParams(p string) *IFS {
	z := make([]*Affine, 0)

	if len(p) > 0 && unicode.IsLetter(rune(p[0])) {
		p2, ok := Builtin[p]
		if !ok {
			panic(Errorf("Unknown Builtin: %q", p))
		}
		p = p2
	}

	pp := strings.Split(p, ",")

	for _, f := range pp {
		ff := strings.Split(f, " ")
		if len(ff) != 6 {
			panic(Errorf("bad IFS params: %q %q %q", p, f, ff))
		}
		aff := &Affine{}
		if ff[0][0]=='@' {
			var xa, ya, xb, yb, xc, yc float64
			CheckOne(Sscanf(ff[0][1:], "%f", &xa))
			CheckOne(Sscanf(ff[1], "%f", &ya))
			CheckOne(Sscanf(ff[2], "%f", &xb))
			CheckOne(Sscanf(ff[3], "%f", &yb))
			CheckOne(Sscanf(ff[4], "%f", &xc))
			CheckOne(Sscanf(ff[5], "%f", &yc))
			aff.A = xb - xa
			aff.C = yb - ya
			aff.E = xc - aff.A
			aff.F = yc - aff.C
			aff.B = xa - aff.E
			aff.D = ya - aff.F
		} else {
			CheckOne(Sscanf(ff[0], "%f", &aff.A))
			CheckOne(Sscanf(ff[1], "%f", &aff.B))
			CheckOne(Sscanf(ff[2], "%f", &aff.C))
			CheckOne(Sscanf(ff[3], "%f", &aff.D))
			CheckOne(Sscanf(ff[4], "%f", &aff.E))
			CheckOne(Sscanf(ff[5], "%f", &aff.F))
		}
		z = append(z, aff)
	}
	return &IFS{Choices: z}
}
