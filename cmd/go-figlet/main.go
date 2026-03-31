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
	fontPath := flag.String("font", "", "path to .flf font file")
	colorStr := flag.String("color", "", "colors: preset (rainbow, warm, cool), names, or hex (semicolon-separated)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-font font.flf] [-color colors] text...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nColor options:\n")
		fmt.Fprintf(os.Stderr, "  presets: rainbow, warm, cool\n")
		fmt.Fprintf(os.Stderr, "  names:   red;green;blue\n")
		fmt.Fprintf(os.Stderr, "  hex:     FF0000;00FF00;0000FF\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *fontPath == "" {
		candidates := []string{
			"/usr/share/figlet/standard.flf",
			"/usr/local/share/figlet/standard.flf",
			"/usr/share/figlet/fonts/standard.flf",
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				*fontPath = c
				break
			}
		}
		if *fontPath == "" {
			fmt.Fprintln(os.Stderr, "error: no font file specified and no default font found")
			fmt.Fprintln(os.Stderr, "use -f option to specify a .flf font file")
			os.Exit(1)
		}
	}

	var colors []figlet.Color
	if *colorStr != "" {
		var err error
		colors, err = figlet.ParseColors(*colorStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	font, err := figlet.LoadFont(*fontPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	printFn := func(text string) {
		if len(colors) > 0 {
			font.PrintWithColor(text, colors)
		} else {
			font.Print(text)
		}
	}

	if flag.NArg() > 0 {
		printFn(strings.Join(flag.Args(), " "))
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			printFn(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
}
