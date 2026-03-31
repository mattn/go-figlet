package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/mattn/go-figlet"
)

func TestWriteGlyphPreservesRightMargin(t *testing.T) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	writeGlyph(w, &figlet.Glyph{
		Lines: []string{"#"},
		Width: 1,
	}, 2, '$', 1)

	if err := w.Flush(); err != nil {
		t.Fatalf("flush failed: %v", err)
	}

	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}

	if got, want := lines[0], "╗ $@"; got != want {
		t.Fatalf("first line = %q, want %q", got, want)
	}
	if got, want := lines[1], "╚ $@@"; got != want {
		t.Fatalf("second line = %q, want %q", got, want)
	}
}
