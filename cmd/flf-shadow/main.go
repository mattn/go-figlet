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
	double := flag.Bool("double", true, "double width (2 chars per pixel)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: flf-shadow [-double=false] input.flf > output.flf\n")
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

	pixelWidth := 2
	if !*double {
		pixelWidth = 1
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

	maxWidth := font.MaxLen*pixelWidth + 2
	if pixelWidth == 1 {
		maxWidth = font.MaxLen + 2
	}
	fmt.Fprintf(out, "flf2a$ %d %d %d -1 2 0 16384 %d\n",
		newHeight, font.Baseline+1, maxWidth, len(extras))
	fmt.Fprintf(out, "ANSI Shadow font generated from efont.\n\n")

	for _, r := range requiredChars {
		g, ok := font.Glyphs[r]
		if ok {
			writeGlyph(out, g, newHeight, font.HardBlank, pixelWidth)
		} else {
			writeBlank(out, 2, newHeight)
		}
	}
	for _, r := range extras {
		fmt.Fprintf(out, "%d\n", r)
		writeGlyph(out, font.Glyphs[r], newHeight, font.HardBlank, pixelWidth)
	}
}

func writeGlyph(w *bufio.Writer, g *figlet.Glyph, newHeight int, hardBlank rune, pw int) {
	origH := len(g.Lines)
	origW := g.Width

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

	f := func(y, x int) bool {
		if y < 0 || y >= origH || x < 0 || x >= origW {
			return false
		}
		return orig[y][x]
	}

	outW := origW * pw
	if pw == 1 {
		outW = origW + 1
	}

	grid := make([][]rune, newHeight)
	for i := range grid {
		grid[i] = make([]rune, outW)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	for y := 0; y < newHeight; y++ {
		for x := 0; x < origW; x++ {
			cur := f(y, x)
			above := f(y-1, x)
			right := f(y, x+1)
			aboveRight := f(y-1, x+1)
			left := f(y, x-1)
			aboveLeft := f(y-1, x-1)

			if pw == 2 {
				L := x * 2
				R := x*2 + 1
				if cur {
					grid[y][L] = '█'
					if right {
						grid[y][R] = '█'
					} else if above && aboveRight {
						grid[y][R] = '╔'
					} else if above && !aboveRight {
						grid[y][R] = '║'
					} else {
						grid[y][R] = '╗'
					}
				} else {
					if left && above && aboveLeft {
						grid[y][L] = '═'
					} else if left && above && !aboveLeft {
						grid[y][L] = '╚'
					} else if above && !aboveLeft {
						grid[y][L] = '╚'
					} else if above {
						grid[y][L] = '═'
					}
					if above && right {
						grid[y][R] = '═'
					} else if above && !aboveRight {
						grid[y][R] = '╝'
					} else if above {
						grid[y][R] = '═'
					}
				}
				continue
			}

			L := x
			R := x + 1
			if cur {
				grid[y][L] = '█'
				if right {
					if grid[y][R] == ' ' {
						grid[y][R] = '█'
					}
				} else if above && aboveRight {
					if grid[y][R] == ' ' {
						grid[y][R] = '╔'
					}
				} else if above && !aboveRight {
					if grid[y][R] == ' ' {
						grid[y][R] = '║'
					}
				} else {
					if grid[y][R] == ' ' {
						grid[y][R] = '╗'
					}
				}
			} else {
				if left && above && aboveLeft {
					setShadowCell(grid[y], L, '═')
				} else if left && above && !aboveLeft {
					setShadowCell(grid[y], L, '╚')
				} else if above && !aboveLeft {
					setShadowCell(grid[y], L, '╚')
				} else if above {
					setShadowCell(grid[y], L, '═')
				}
				if above && right {
					setShadowCell(grid[y], R, '═')
				} else if above && !aboveRight {
					setShadowCell(grid[y], R, '╝')
				} else if above {
					setShadowCell(grid[y], R, '═')
				}
			}
		}
	}

	for y := 0; y < newHeight; y++ {
		if y > 0 {
			fmt.Fprint(w, "\n")
		}
		fmt.Fprint(w, strings.TrimRight(string(grid[y]), " "))
		fmt.Fprint(w, "$@")
	}
	fmt.Fprint(w, "@\n")
}

func setShadowCell(row []rune, x int, r rune) {
	if x < 0 || x >= len(row) {
		return
	}
	if row[x] == ' ' {
		row[x] = r
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
