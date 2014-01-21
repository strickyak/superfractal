/*
	ifs_png writes a normal (level-1) fractal to stdout as a PNG.

	Give parameters with a string to -p (separate matrix items with " " and
	ifs member functions with ",").   Or name a builtin IFS with -p.
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

	green := canvas.RGB(0, 255, 0)
	c := canvas.NewCanvas(*width, *height)
	ifs := superfractal.ParseIfsParams(*params)
	ifs.ChaosGame(c, *num, green)
	c.WritePng(os.Stdout)
}
