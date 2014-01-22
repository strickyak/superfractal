package superfractal

import (
	"github.com/strickyak/canvas"
	"math/rand"
	"strings"
)

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

func (o *Triptych) Fill(color []canvas.Color) {
	for i := 0; i < len(o.Panels); i++ {
		c := o.Panels[i]
		c.Fill(0, 0, c.Width, c.Height, color[i])
	}
}

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

func ParseListOfIfsParams(p string) []*IFS {
	parts := strings.Split(p, ";")
	z := make([]*IFS, len(parts))
	for i, part := range parts {
		z[i] = ParseIfsParams(part)
	}
	return z
}
