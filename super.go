package superfractal

import (
	"github.com/strickyak/canvas"
	"math/rand"
	"strings"
)

// Triptych is an assembly of n panels (dispite the name, n does not have to be 3).
type Triptych struct {
	Panels []*canvas.Canvas
}

func NewTriptych(numPanels, width, height int) *Triptych {
	z := &Triptych{
		Panels: make([]*canvas.Canvas, numPanels),
	}
	for i := 0; i < numPanels; i++ {
		z.Panels[i] = canvas.NewCanvas(width, height)
	}
	return z
}

// Execute one step building a superfractal, reading images from src,
// and writing on the receiver, which should be a new blank Triptych.
func (o *Triptych) SuperStep(src *Triptych, ifss []*IFS) {
	wait := make(chan int)
	for i := 0; i < len(o.Panels); i++ {
		dest := o.Panels[i]

		// Pick an IFS at random.
		ifs := ifss[rand.Intn(len(ifss))]

		// For each map in the ifs,
		for j := 0; j < len(ifs.Choices); j++ {
			affine := ifs.Choices[j]

			// pick a source, and map it onto the dest.
			src := src.Panels[rand.Intn(len(src.Panels))]
			go affine.MapImageTo(src, dest, i*1000+j, wait)
		}
		for j := 0; j < len(ifs.Choices); j++ {
			k := <-wait
			canvas.Say("Done", k)
		}
	}
}

// ParseListOfIfsParams splits the string on ";", and the pieces are parsed with ParseIfsParams.
func ParseListOfIfsParams(p string) []*IFS {
	parts := strings.Split(p, ";")
	z := make([]*IFS, len(parts))
	for i, part := range parts {
		z[i] = ParseIfsParams(part)
	}
	return z
}
