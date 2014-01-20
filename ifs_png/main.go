/*
	ifs_png writes a normal (level-1) fractal to stdout as a PNG.

	Give parameters with a string to -p (separate matrix items with "," and
	functions with " ").

	Or plot a builtin IFS, naming it with -b.
*/
package main

import (
	. "fmt"
	"flag"
	"github.com/strickyak/canvas"
	"github.com/strickyak/superfractal"
	"os"
)

var (
	num  = flag.Int("n", 1000000, "number points to plot")
	width  = flag.Int("w", 1000, "width in pixels")
	height = flag.Int("h", 1000, "width in pixels")
	params = flag.String("p", "", "IFS parameters: a,b,c,d,e,f ...")
	builtin = flag.String("b", "", "Use an IFS by name from the library.")
)

func main() {
	flag.Parse()
	if *builtin != "" {
		var ok bool
		*params, ok = superfractal.Builtin[*builtin]
		if !ok {
			panic(Errorf("Cannot find builtin: %q", *builtin))
		}
	}
	green := canvas.RGB(0, 255, 0)
	c := canvas.NewCanvas(*width, *height)
	ifs := superfractal.ParseIfsParams(*params)
	ifs.ChaosGame(c, *num, green)
	c.WritePng(os.Stdout)
}
