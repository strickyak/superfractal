package superfractal

import (
	// . "fmt"
	"math/rand"
	// "strings"
	// "unicode"

	// "github.com/strickyak/canvas"
)

type CodeTree struct {
	// Set on creation:
	Depth  int
	Fanout int
	IFSs   []*IFS

	// Computed by Init:
	Codes [][]byte
}

func (o *CodeTree) Init(seed int) {
	src := rand.NewSource(int64(seed))
	r := rand.New(src)
	o.Codes = make([][]byte, o.Depth)
	pow := 1
	for i := 0; i < o.Depth; i++ {
		o.Codes[i] = make([]byte, pow)
		for j, _ := range o.Codes[i] {
			o.Codes[i][j] = byte(r.Int31())
		}
		pow *= o.Fanout
	}
}

func (o *CodeTree) LookupPath(path []byte, out []byte) {
	out[0] = o.Codes[0][0]
	off := 0
	for i, p := range path[1:] {
		off += int(p)
		out[i+1] = o.Codes[i+1][off]
		off *= o.Fanout
	}
}

func (o *CodeTree) TransformPath(path []byte, choices []byte) (float64, float64) {
	n := len(o.IFSs)
	var x, y float64
	/*
		for i := len(path) - 1; i >= 0; i-- {
			ifs := o.IFSs[int(choices[i]) % n]
			fn := ifs.Choices[path[i]]
			x, y = fn.Apply(x, y)
		}
	*/
	for i := len(path) - 1; i > 0; i-- {
		ifs := o.IFSs[int(choices[i-1])%n]
		fn := ifs.Choices[path[i]]
		x, y = fn.Apply(x, y)
	}

	return x, y
}

func (o *CodeTree) RandomPoint(tmp1, tmp2 []byte) (float64, float64) {
	// Choose the path.
	n := len(tmp1)
	for i := 0; i < n; i++ {
		tmp1[i] = byte(rand.Intn(o.Fanout))
	}

	o.LookupPath(tmp1, tmp2)

	return o.TransformPath(tmp1, tmp2)
}
