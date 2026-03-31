package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/mattn/go-figlet"
)

func TestWriteBloodyGlyphSinglePixel(t *testing.T) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	cfg := bloodyConfig{dripRatio: 1, dripDensity: 1}
	writeBloodyGlyph(w, &figlet.Glyph{
		Lines: []string{"#"},
		Width: 1,
	}, 3, '$', 1, cfg)

	if err := w.Flush(); err != nil {
		t.Fatalf("flush failed: %v", err)
	}

	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}

	if got, want := lines[0], "▓$@"; got != want {
		t.Fatalf("first line = %q, want %q", got, want)
	}
	if got, want := lines[1], "▒$@"; got != want {
		t.Fatalf("second line = %q, want %q", got, want)
	}
	if got, want := lines[2], "$@@"; got != want {
		t.Fatalf("third line = %q, want %q", got, want)
	}
}

func TestDripLengthScalesWithHeight(t *testing.T) {
	cfg := bloodyConfig{dripRatio: 0.12, dripDensity: 1}

	if got, want := maxDripLength(10, cfg), 1; got != want {
		t.Fatalf("maxDripLength(10) = %d, want %d", got, want)
	}
	if got, want := maxDripLength(24, cfg), 3; got != want {
		t.Fatalf("maxDripLength(24) = %d, want %d", got, want)
	}
}

func TestDripLengthIsDeterministic(t *testing.T) {
	cfg := bloodyConfig{dripRatio: 0.2, dripDensity: 1}
	if got, want := dripLength(0, 0, 1, 1, cfg), 1; got != want {
		t.Fatalf("dripLength(0,0,1,1) = %d, want %d", got, want)
	}
	if got, want := dripLength(3, 4, 10, 12, cfg), dripLength(3, 4, 10, 12, cfg); got != want {
		t.Fatalf("dripLength should be deterministic")
	}
}
