package figlet

import (
	"strings"
	"testing"
	"testing/fstest"
)

const testFLF = `flf2a$ 5 4 6 -1 1 0 0 0
Test font
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
  ##  @
 #  # @
 #### @
 #  # @
 #  # @@
 ###  @
 #  # @
 ###  @
 #  # @
 ###  @@
  ##  @
 #  # @
 #    @
 #  # @
  ##  @@
 ###  @
 #  # @
 #  # @
 #  # @
 ###  @@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
 $@
 $@
 $@
 $@
 $@@
`

func TestReadFont(t *testing.T) {
	font, err := ReadFont(strings.NewReader(testFLF))
	if err != nil {
		t.Fatal(err)
	}
	if font.Height != 5 {
		t.Errorf("expected height 5, got %d", font.Height)
	}
	if font.HardBlank != '$' {
		t.Errorf("expected hardblank '$', got %c", font.HardBlank)
	}
	glyph, ok := font.Glyphs['A']
	if !ok {
		t.Fatal("glyph for 'A' not found")
	}
	if len(glyph.Lines) != 5 {
		t.Errorf("expected 5 lines for 'A', got %d", len(glyph.Lines))
	}
}

func TestRender(t *testing.T) {
	font, err := ReadFont(strings.NewReader(testFLF))
	if err != nil {
		t.Fatal(err)
	}
	result := font.Render("ABCD")
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}
	// First line should contain " ## " from A and " ### " from B, etc.
	if !strings.Contains(lines[0], "##") {
		t.Errorf("expected ## in first line, got %q", lines[0])
	}
}

func TestLoadFont(t *testing.T) {
	// Test with the standard figlet font if available
	font, err := LoadFont("/usr/share/figlet/standard.flf")
	if err != nil {
		t.Skip("standard.flf not available:", err)
	}
	result := font.Render("Hi")
	if result == "" {
		t.Error("expected non-empty result")
	}
	t.Log("\n" + result)
}

func TestLoadFontFS(t *testing.T) {
	font, err := LoadFontFS(fstest.MapFS{
		"fonts/test.flf": &fstest.MapFile{Data: []byte(testFLF)},
	}, "fonts/test.flf")
	if err != nil {
		t.Fatal(err)
	}
	if font.Height != 5 {
		t.Errorf("expected height 5, got %d", font.Height)
	}
	if _, ok := font.Glyphs['A']; !ok {
		t.Fatal("glyph for 'A' not found")
	}
}
