/*
	GOMAXPROCS=4 go run super_png/main.go -p='s;t' -n=20 -h=3000 -w=3000

	GOMAXPROCS=4 go run super_png/main.go -np=20 -p='g;s;t' -n=32 -h=3000 -w=3000  

	GOMAXPROCS=4 go run super_png/main.go -np=32 -p='g;s;t' -n=32 -h=4000 -w=4000  
*/
package main

import (
	"flag"
	. "fmt"
	"github.com/strickyak/canvas"
	"github.com/strickyak/superfractal"
	"os"
)

var (
	base      = flag.String("base", "_", "basename for output images")
	numPanels = flag.Int("np", 5, "number of panels")
	num       = flag.Int("n", 1000000, "number if iterations")
	width     = flag.Int("w", 1000, "width in pixels")
	height    = flag.Int("h", 1000, "width in pixels")
	params    = flag.String("p", "", "IFS parameters, as lists of matrices, or name of a Builtin IFS.")
	list      = flag.Bool("l", false, "List the library.")
)

func main() {
	flag.Parse()

	// Special -l list command.
	if *list {
		for k, v := range superfractal.Builtin {
			Printf("%s\t%s\n", k, v)
		}
		return
	}

	ifss := superfractal.ParseListOfIfsParams(*params)

	src := superfractal.NewTriptych(*numPanels, *width, *height)
	colors := make([]canvas.Color, *numPanels)
	r, g, b := byte(255), byte(0), byte(0)
	for i := 0; i < *numPanels; i++ {
		colors[i] = canvas.RGB(r, g, b)
		r, g, b = byte(0.95*float64(b)), byte(0.95*float64(r)), byte(0.95*float64(g))
	}
	src.Fill(colors)
	for i := 0; i < *num; i++ {
		targ := superfractal.NewTriptych(*numPanels, *width, *height)
		targ.SuperStep(src, ifss)

		for j := 0; j < *numPanels; j++ {
			f, err := os.Create(Sprintf("%s.%03d.%03d.png", *base, i, j))
			if err != nil {
				panic(err)
			}
			targ.Panels[j].WritePng(f)
			f.Close()
		}
		src = targ
	}

	// c.WritePng(os.Stdout)
}
