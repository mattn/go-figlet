package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-figlet"
)

func main() {
	double := flag.Bool("double", false, "double width (2 chars per pixel)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: flf-3d [-double] input.flf > output.flf\n")
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

	newHeight := font.Height + 1

	requiredChars := []rune{
		32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
		64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
		80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
		96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109,
		110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
		123, 124, 125, 126, 196, 214, 220, 228, 246, 252, 223,
	}
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

	maxWidth := font.MaxLen*pw + 2
	fmt.Fprintf(out, "flf2a$ %d %d %d -1 2 0 16384 %d\n",
		newHeight, font.Baseline+1, maxWidth, len(extras))
	fmt.Fprintf(out, "3D shadow font generated from efont.\n\n")

	for _, r := range requiredChars {
		g, ok := font.Glyphs[r]
		if ok {
			write3DGlyph(out, g, newHeight, font.HardBlank, pw)
		} else {
			writeBlank(out, 2, newHeight)
		}
	}

	for _, r := range extras {
		g := font.Glyphs[r]
		fmt.Fprintf(out, "%d\n", r)
		write3DGlyph(out, g, newHeight, font.HardBlank, pw)
	}
}

func write3DGlyph(w *bufio.Writer, g *figlet.Glyph, newHeight int, hardBlank rune, pw int) {
	origH := len(g.Lines)
	origW := g.Width
	outW := origW * pw

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

	grid := make([][]rune, newHeight)
	for i := range grid {
		grid[i] = make([]rune, outW)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Place foreground
	for y := 0; y < origH; y++ {
		for x := 0; x < origW; x++ {
			if orig[y][x] {
				for d := 0; d < pw; d++ {
					grid[y][x*pw+d] = '█'
				}
			}
		}
	}

	// Place shadow: for each filled pixel at (y, x), put ░ at (y+1) shifted left by 1
	for y := 0; y < origH; y++ {
		for x := 0; x < origW; x++ {
			if !orig[y][x] {
				continue
			}
			sy := y + 1
			if sy >= newHeight {
				continue
			}
			for d := 0; d < pw; d++ {
				sx := x*pw + d - 1
				if sx >= 0 && grid[sy][sx] == ' ' {
					grid[sy][sx] = '░'
				}
			}
		}
	}

	for y := 0; y < newHeight; y++ {
		if y > 0 {
			fmt.Fprint(w, "\n")
		}
		line := strings.TrimRight(string(grid[y]), " ")
		fmt.Fprint(w, line)
		fmt.Fprint(w, "$@")
	}
	fmt.Fprint(w, "@\n")
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
