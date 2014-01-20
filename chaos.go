package superfractal

import (
	. "fmt"
	"math/rand"
	"strings"

	"github.com/strickyak/canvas"
)

const SKIP = 32

var Builtin = make(map[string]string)

func init() {
	// Sierpinski's Triangle.
	Builtin["s"] = "0.5,0,0,0.5,0,0 0.5,0,0,0.5,0.5,0 0.5,0,0,0.5,0,0.5"

	// Square Gasket.
	Builtin["g"] = "0.4,0,0,0.4,0,0 0.4,0,0,0.4,0.6,0 0.4,0,0,0.4,0.6,0.6 0.4,0,0,0.4,0,0.6"
}

type Affine struct {
	A, B, C, D, E, F float64
}

type IFS struct {
	Choices []*Affine
}

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

func (o IFS) Choose() *Affine {
	r := rand.Intn(len(o.Choices))
	return o.Choices[r]
}

func CheckOne(n int, err error) {
	if n != 1 || err != nil {
		panic(Errorf("Bad Number: %s", err))
	}
}

func ParseIfsParams(p string) *IFS {
	z := make([]*Affine, 0)
	pp := strings.Split(p, " ")
	for _, f := range pp {
		ff := strings.Split(f, ",")
		if len(ff) != 6 {
			panic(Errorf("bad IFS params: %q %q %q", p, f, ff))
		}
		aff := &Affine{}
		CheckOne(Sscanf(ff[0], "%f", &aff.A))
		CheckOne(Sscanf(ff[1], "%f", &aff.B))
		CheckOne(Sscanf(ff[2], "%f", &aff.C))
		CheckOne(Sscanf(ff[3], "%f", &aff.D))
		CheckOne(Sscanf(ff[4], "%f", &aff.E))
		CheckOne(Sscanf(ff[5], "%f", &aff.F))
		z = append(z, aff)
	}
	return &IFS{Choices: z}
}
