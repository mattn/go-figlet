package figlet

import (
	"fmt"
	"strconv"
	"strings"
)

// Color represents a terminal color for rendering.
type Color interface {
	prefix() string
	suffix() string
}

const ansiReset = "\x1b[0m"

// ANSI color constants.
var (
	ColorBlack   Color = ansiColor(30)
	ColorRed     Color = ansiColor(31)
	ColorGreen   Color = ansiColor(32)
	ColorYellow  Color = ansiColor(33)
	ColorBlue    Color = ansiColor(34)
	ColorMagenta Color = ansiColor(35)
	ColorCyan    Color = ansiColor(36)
	ColorWhite   Color = ansiColor(37)
)

type ansiColor int

func (c ansiColor) prefix() string { return fmt.Sprintf("\x1b[0;%dm", int(c)) }
func (c ansiColor) suffix() string { return ansiReset }

// TrueColor represents a 24-bit RGB color.
type TrueColor struct {
	R, G, B uint8
}

func (c TrueColor) prefix() string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

func (c TrueColor) suffix() string { return ansiReset }

// NewTrueColor creates a TrueColor from RGB values.
func NewTrueColor(r, g, b uint8) Color {
	return TrueColor{R: r, G: g, B: b}
}

// NewTrueColorFromHex creates a TrueColor from a hex string like "FF8800".
func NewTrueColorFromHex(hex string) (Color, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return nil, fmt.Errorf("figlet: invalid hex color %q", hex)
	}
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("figlet: invalid hex color %q: %w", hex, err)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("figlet: invalid hex color %q: %w", hex, err)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return nil, fmt.Errorf("figlet: invalid hex color %q: %w", hex, err)
	}
	return TrueColor{R: uint8(r), G: uint8(g), B: uint8(b)}, nil
}

// colorByName returns a Color for the given name, or nil if unknown.
func colorByName(name string) Color {
	switch strings.ToLower(name) {
	case "black":
		return ColorBlack
	case "red":
		return ColorRed
	case "green":
		return ColorGreen
	case "yellow":
		return ColorYellow
	case "blue":
		return ColorBlue
	case "magenta":
		return ColorMagenta
	case "cyan":
		return ColorCyan
	case "white":
		return ColorWhite
	}
	return nil
}

// Color presets.
var presets = map[string][]Color{
	"rainbow": {
		TrueColor{255, 0, 0},
		TrueColor{255, 127, 0},
		TrueColor{255, 255, 0},
		TrueColor{0, 200, 0},
		TrueColor{0, 127, 255},
		TrueColor{75, 0, 130},
		TrueColor{148, 0, 211},
	},
	"warm": {
		TrueColor{255, 0, 0},
		TrueColor{255, 100, 0},
		TrueColor{255, 165, 0},
		TrueColor{255, 200, 0},
		TrueColor{255, 255, 0},
	},
	"cool": {
		TrueColor{0, 200, 0},
		TrueColor{0, 180, 180},
		TrueColor{0, 127, 255},
		TrueColor{75, 0, 130},
		TrueColor{148, 0, 211},
	},
}

// ColorPreset returns a preset color palette by name.
// Available presets: "rainbow", "warm", "cool".
func ColorPreset(name string) []Color {
	return presets[strings.ToLower(name)]
}

// ParseColor parses a color string, which can be a name (e.g. "red")
// or a hex code (e.g. "FF8800" or "#FF8800").
func ParseColor(s string) (Color, error) {
	if c := colorByName(s); c != nil {
		return c, nil
	}
	return NewTrueColorFromHex(s)
}

// ParseColors parses a semicolon-separated color string.
// Supports preset names (e.g. "rainbow"), color names, and hex codes.
func ParseColors(s string) ([]Color, error) {
	s = strings.TrimSpace(s)
	if p := ColorPreset(s); p != nil {
		return p, nil
	}
	var colors []Color
	for _, part := range strings.Split(s, ";") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		c, err := ParseColor(part)
		if err != nil {
			return nil, err
		}
		colors = append(colors, c)
	}
	return colors, nil
}
