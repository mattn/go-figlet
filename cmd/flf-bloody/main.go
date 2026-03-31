package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/mattn/go-figlet"
)

var requiredChars = []rune{
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
	48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
	64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
	80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
	96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109,
	110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
	123, 124, 125, 126, 196, 214, 220, 228, 246, 252, 223,
}

func main() {
	double := flag.Bool("double", false, "double width (2 chars per pixel)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: flf-bloody [-double] input.flf > output.flf\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	font, err := figlet.LoadFont(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	pw := 1
	if *double {
		pw = 2
	}

	newHeight := font.Height + 2

	requiredSet := map[rune]bool{}
	for _, r := range requiredChars {
		requiredSet[r] = true
	}

	var extras []rune
	for r := range font.Glyphs {
		if !requiredSet[r] {
			extras = append(extras, r)
		}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	maxWidth := font.MaxLen*pw + 1
	fmt.Fprintf(out, "flf2a$ %d %d %d -1 2 0 16384 %d\n",
		newHeight, font.Baseline+1, maxWidth, len(extras))
	fmt.Fprintf(out, "Bloody font generated from bitmap glyphs.\n\n")

	for _, r := range requiredChars {
		g, ok := font.Glyphs[r]
		if ok {
			writeBloodyGlyph(out, g, newHeight, font.HardBlank, pw)
		} else {
			writeBlank(out, 2, newHeight)
		}
	}

	for _, r := range extras {
		fmt.Fprintf(out, "%d\n", r)
		writeBloodyGlyph(out, font.Glyphs[r], newHeight, font.HardBlank, pw)
	}
}

func writeBloodyGlyph(w *bufio.Writer, g *figlet.Glyph, newHeight int, hardBlank rune, pw int) {
	origH := len(g.Lines)
	origW := g.Width
	outW := origW*pw + 1

	orig := make([][]bool, origH)
	for i, line := range g.Lines {
		runes := []rune(line)
		row := make([]bool, origW)
		for j := 0; j < origW && j < len(runes); j++ {
			if runes[j] != ' ' && runes[j] != hardBlank {
				row[j] = true
			}
		}
		orig[i] = row
	}

	filled := func(y, x int) bool {
		if y < 0 || y >= origH || x < 0 || x >= origW {
			return false
		}
		return orig[y][x]
	}

	grid := make([][]rune, newHeight)
	for i := range grid {
		grid[i] = make([]rune, outW)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for y := 0; y < origH; y++ {
		for x := 0; x < origW; x++ {
			if !filled(y, x) {
				continue
			}

			exposedTop := !filled(y-1, x)
			exposedBottom := !filled(y+1, x)
			exposedLeft := !filled(y, x-1)
			exposedRight := !filled(y, x+1)

			ch := rune('█')
			switch {
			case exposedTop && !exposedBottom:
				ch = '▄'
			case exposedBottom && (exposedLeft || exposedRight):
				ch = '▓'
			case exposedLeft || exposedRight:
				ch = '▓'
			}

			for dx := 0; dx < pw; dx++ {
				grid[y][x*pw+dx] = ch
			}

			if exposedBottom {
				dripLen := dripLength(x, y, origW, origH)
				for d := 1; d <= dripLen && y+d < newHeight; d++ {
					dripChar := rune('▓')
					if d == dripLen {
						dripChar = '▒'
					}
					for dx := 0; dx < pw; dx++ {
						sx := x*pw + dx
						if grid[y+d][sx] == ' ' {
							grid[y+d][sx] = dripChar
						}
					}
				}
			}
		}
	}

	for y := 0; y < newHeight; y++ {
		if y > 0 {
			fmt.Fprint(w, "\n")
		}
		fmt.Fprint(w, string(grid[y]))
		fmt.Fprint(w, "$@")
	}
	fmt.Fprint(w, "@\n")
}

func dripLength(x, y, width, height int) int {
	seed := (x*17 + y*31 + width*7 + height*13) % 8
	switch seed {
	case 0, 1:
		return 2
	case 2, 3, 4:
		return 1
	default:
		return 0
	}
}

func writeBlank(w *bufio.Writer, width, height int) {
	for y := 0; y < height; y++ {
		if y > 0 {
			fmt.Fprint(w, "\n")
		}
		for x := 0; x < width; x++ {
			fmt.Fprint(w, " ")
		}
		fmt.Fprint(w, "$@")
	}
	fmt.Fprint(w, "@\n")
}
