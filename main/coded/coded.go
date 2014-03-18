/*
	Draw a (member of a) superfractal using the IFSs named by -p='ifs1;ifs2;ifs3;...'
	outputting PNGs whose filenames are suffixed from -base.

	See: Barnsley, Michael, John E. Hutchinson, and Orjan Stenflo. "V-variable fractals and superfractals."
	But we're doing it differently.

	go build main/coded/coded.go && for s in 1 2 3 4 5 6 ; do time ./coded  -o _e$s.png -p 'W;X;Y' -n 100000 -s $s -d 8 -fuzz=8 ; done

	go build main/coded/coded.go && for s in 1 2 3 4 5 6 ; do time ./coded  -o _e$s.png -p 'g;f' -n 5000000 -s $s -d 18 -fuzz=1 -w=5000 -h=5000 ; done

	 go build main/coded/coded.go && for s in 1 2 3 4 5 6 ; do time ./coded  -o _e$s.png -p 'f;f;g' -f 4 -n 1000000 -s $s -d 14 -fuzz=1 -w=2000 -h=2000 ; done


*/
package main

import (
	"flag"
	. "fmt"
	"math/rand"
	"os"

	"github.com/strickyak/canvas"
	"github.com/strickyak/superfractal"
)

var (
	filename = flag.String("o", "_", "basename for output images")
	seed     = flag.Int("s", 1, "seed for code tree")
	// fanout   = flag.Int("f", 3, "depth of code tree")
	depth    = flag.Int("d", 12, "depth of code tree")
	num      = flag.Int("n", 1000, "number if iterations")
	width    = flag.Int("w", 1000, "width in pixels")
	height   = flag.Int("h", 1000, "width in pixels")
	params   = flag.String("p", "", "IFS parameters, as lists of matrices, or name of a Builtin IFS.")
	list     = flag.Bool("l", false, "Just list short names for builtin IFSs.")
	fuzz     = flag.Int("fuzz", 1, "Fuzzy Point Size")
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
	tree := &superfractal.CodeTree{
		Depth:  *depth,
		IFSs:   ifss,
	}
	tree.Init(*seed)

	can := canvas.NewCanvas(*width, *height)
	tmp1 := make([]byte, *depth)
	tmp2 := make([]byte, *depth)
	for i := 0; i < *num; i++ {
		x, y := tree.RandomPoint(tmp1, tmp2)

		if *fuzz > 1 {
			xf := rand.Intn(*fuzz)
			yf := rand.Intn(*fuzz)
			can.FSet(x+float64(xf)/float64(*width), y+float64(yf)/float64(*height), canvas.White)
		} else {
			can.FSet(x, y, canvas.White)
		}
		/*
			for xf := 0; xf < *fuzz; xf++ {
			for yf := 0; yf < *fuzz; yf++ {
				can.FSet(x + float64(xf)/float64(*width), y + float64(yf)/float64(*height), canvas.White)
			}
			}
		*/
	}

	canvas.Say(*filename)
	fd, err := os.Create(*filename)
	if err != nil {
		panic(err)
	}
	can.WritePng(fd)
	fd.Close()
}
