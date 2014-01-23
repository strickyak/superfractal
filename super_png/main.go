/*
	Draw a superfractal using the IFSs named by -p='ifs1;ifs2;ifs3;...'
	outputting PNGs whose filenames are suffixed from -base.

	See: Barnsley, Michael, John E. Hutchinson, and Orjan Stenflo. "V-variable fractals and superfractals."

	GOMAXPROCS=4 go run super_png/main.go -p='s;t' -n=20

	GOMAXPROCS=4 go run super_png/main.go -np=30 -p='g;s;t' -n=32 -h=3000 -w=3000
*/
package main

import (
	"flag"
	. "fmt"
	"os"

	"github.com/strickyak/canvas"
	"github.com/strickyak/superfractal"
)

var (
	base      = flag.String("base", "_", "basename for output images")
	numPanels = flag.Int("np", 5, "number of panels")
	num       = flag.Int("n", 1000000, "number if iterations")
	width     = flag.Int("w", 1000, "width in pixels")
	height    = flag.Int("h", 1000, "width in pixels")
	params    = flag.String("p", "", "IFS parameters, as lists of matrices, or name of a Builtin IFS.")
	list      = flag.Bool("l", false, "List short names for builtin IFSs.")
)

func main() {
	flag.Parse()

	// Special -l list command, listing short names for builtin IFSs.
	if *list {
		for k, v := range superfractal.Builtin {
			Printf("%s\t%s\n", k, v)
		}
		return
	}

	// Parse IFSs & construct initial Triptych.
	ifss := superfractal.ParseListOfIfsParams(*params)
	src := superfractal.NewTriptych(*numPanels, *width, *height)

	// Fill initial Triptych with different colors.
	r, g, b := byte(255), byte(0), byte(0)
	for i := 0; i < *numPanels; i++ {
		src.Panels[i].Fill(0, 0, *width, *height, canvas.RGB(r, g, b))

		// Rotate & dim the colors a bit, on each iteration.
		// TODO: something better than this.
		r, g, b = byte(0.96*float64(b)), byte(0.96*float64(r)), byte(0.96*float64(g))
	}

	// Loop for *num steps, creating successive superfractal approximations.
	for i := 0; i < *num; i++ {
		targ := superfractal.NewTriptych(*numPanels, *width, *height)
		targ.SuperStep(src, ifss)

		for j := 0; j < *numPanels; j++ {
			filename := Sprintf("%s.%03d.%03d.png", *base, i, j)
			canvas.Say(filename)
			f, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			targ.Panels[j].WritePng(f)
			f.Close()
		}
		src = targ
	}
}
