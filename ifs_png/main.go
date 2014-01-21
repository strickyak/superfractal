/*
	ifs_png writes a normal (level-1) fractal to stdout as a PNG.

	Give parameters with a string to -p (separate matrix items with " " and
	ifs member functions with ",").   Or name a builtin IFS with -p.

	-s=0 is chaos game. -s=1 is mapping whole pictures (can use multiple cores).

	# Deterministic demo:
	GOMAXPROCS=4 time go run main.go -p t -s 1 -n 20 > _1.png

	# Chaos demo:
	time go run main.go -p t -s 0 -n 5000000 > _0.png

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
	num    = flag.Int("n", 1000000, "number points to plot")
	width  = flag.Int("w", 1000, "width in pixels")
	height = flag.Int("h", 1000, "width in pixels")
	params = flag.String("p", "", "IFS parameters, as lists of matrices, or name of a Builtin IFS.")
	list   = flag.Bool("l", false, "List the library.")
	style  = flag.Int("s", 0, "Style: 0=chaos 1=whole")
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

	c := canvas.NewCanvas(*width, *height)
	ifs := superfractal.ParseIfsParams(*params)
	switch *style {
	case 1: // Whole.
		c.Fill(0, 0, *width, *height, canvas.Blue)
		// c.WritePng(os.Stdout)
		// return
		c = ifs.WholeGame(c, *num)
	case 0: // Chaos.
		ifs.ChaosGame(c, *num, canvas.Green)
	}
	c.WritePng(os.Stdout)
}
