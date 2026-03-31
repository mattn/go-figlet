package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"

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
	dripRatio := flag.Float64("drip-ratio", 0.12, "maximum drip length as a fraction of glyph height")
	dripDensity := flag.Float64("drip-density", 0.5, "fraction of exposed lower edges that get drips")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: flf-bloody [-double] [-drip-ratio=0.12] [-drip-density=0.5] input.flf > output.flf\n")
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

	cfg := bloodyConfig{
		dripRatio:   clampFloat(*dripRatio, 0, 1),
		dripDensity: clampFloat(*dripDensity, 0, 1),
	}

	maxDrip := maxDripLength(font.Height, cfg)
	newHeight := font.Height + maxDrip

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

	maxWidth := font.MaxLen * pw
	fmt.Fprintf(out, "flf2a$ %d %d %d -1 2 0 16384 %d\n",
		newHeight, font.Baseline+1, maxWidth, len(extras))
	fmt.Fprintf(out, "Bloody font generated from bitmap glyphs.\n\n")

	for _, r := range requiredChars {
		g, ok := font.Glyphs[r]
		if ok {
			writeBloodyGlyph(out, g, newHeight, font.HardBlank, pw, cfg)
		} else {
			writeBlank(out, 2, newHeight)
		}
	}

	for _, r := range extras {
		fmt.Fprintf(out, "%d\n", r)
		writeBloodyGlyph(out, font.Glyphs[r], newHeight, font.HardBlank, pw, cfg)
	}
}

type bloodyConfig struct {
	dripRatio   float64
	dripDensity float64
}

func writeBloodyGlyph(w *bufio.Writer, g *figlet.Glyph, newHeight int, hardBlank rune, pw int, cfg bloodyConfig) {
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
				dripLen := dripLength(x, y, origW, origH, cfg)
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
		fmt.Fprint(w, strings.TrimRight(string(grid[y]), " "))
		fmt.Fprint(w, "$@")
	}
	fmt.Fprint(w, "@\n")
}

func dripLength(x, y, width, height int, cfg bloodyConfig) int {
	maxLen := maxDripLength(height, cfg)
	if maxLen == 0 {
		return 0
	}

	seed := (x*17 + y*31 + width*7 + height*13) % 100
	if float64(seed)/100 >= cfg.dripDensity {
		return 0
	}

	weightSeed := (x*11 + y*23 + width*5 + height*19) % 100
	weight := 0.35 + (float64(weightSeed) / 100 * 0.65)
	length := int(math.Round(float64(maxLen) * weight))
	if length < 1 {
		length = 1
	}
	if length > maxLen {
		length = maxLen
	}
	return length
}

func maxDripLength(height int, cfg bloodyConfig) int {
	if cfg.dripRatio <= 0 || height <= 0 {
		return 0
	}
	length := int(math.Round(float64(height) * cfg.dripRatio))
	if length < 1 {
		length = 1
	}
	return length
}

func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
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
