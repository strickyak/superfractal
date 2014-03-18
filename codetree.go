package superfractal

import (
	// "fmt"
	"math/rand"
	// "os"
	// "strings"
	// "unicode"

	// "github.com/strickyak/canvas"
)

type CodeTree struct {
	// Set on creation:
	Depth  int
	IFSs   []*IFS

	// Computed by Init:
	fanout int
	codes [][]byte
}

func (o *CodeTree) Init(seed int) {
	// fanout is the max of IFS choices lengths.
	o.fanout = len(o.IFSs[0].Choices)
	for _, ifs := range o.IFSs[1:] {
		if o.fanout < len(ifs.Choices) {
			o.fanout = len(ifs.Choices)
		}
	}

	src := rand.NewSource(int64(seed))
	r := rand.New(src)
	o.codes = make([][]byte, o.Depth)
	pow := 1
	for i := 0; i < o.Depth; i++ {
		o.codes[i] = make([]byte, pow)
		for j, _ := range o.codes[i] {
			o.codes[i][j] = byte(r.Int31())
		}
		pow *= o.fanout
	}
}

// LookupPath takes a random path and looks up which IFS to use for it, at each step.
func (o *CodeTree) LookupPath(path []byte, out []byte) {
	n := len(o.IFSs)
	out[0] = o.codes[0][0]
	off := 0
	for i, p := range path {
		if i > 0 {
			ifs := o.IFSs[int(out[i-1])%n]
			m := len(ifs.Choices)
			off += int(p) % m
			out[i] = o.codes[i][off]
			off *= o.fanout
		}
	}
}

// TransformPath uses 'path' to choose maps, and 'choices' to choose IFSs, and answers a point x,y.
func (o *CodeTree) TransformPath(path []byte, choices []byte) (float64, float64) {
	n := len(o.IFSs)
	var x, y float64
	for i := len(path) - 1; i > 0; i-- {
		ifs := o.IFSs[int(choices[i-1])%n]
		m := len(ifs.Choices)
		fn := ifs.Choices[int(path[i]) % m]
		x, y = fn.Apply(x, y)
	}

	return x, y
}

// RandomPoint picks a random path, looks up IFSs for it, and returns a transformed point x,y.
// It reuses the provided byte buffers for the path and choices.
func (o *CodeTree) RandomPoint(tmp1, tmp2 []byte) (float64, float64) {
	// Choose the path.
	n := len(tmp1)
	for i := 0; i < n; i++ {
		tmp1[i] = byte(rand.Intn(256))
	}

	o.LookupPath(tmp1, tmp2)

	return o.TransformPath(tmp1, tmp2)
}
