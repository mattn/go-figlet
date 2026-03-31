package figlet

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Required characters that must appear at the start of a FIGfont, in order.
// ASCII 32-126 plus German characters.
var requiredChars = []rune{
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
	48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
	64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
	80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
	96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109,
	110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122,
	123, 124, 125, 126, 196, 214, 220, 228, 246, 252, 223,
}

// Glyph represents a single FIGlet character glyph.
type Glyph struct {
	Lines []string
	Width int
}

// Font represents a loaded FIGlet font.
type Font struct {
	Height    int
	Baseline  int
	MaxLen    int
	HardBlank rune
	Glyphs   map[rune]*Glyph
}

// LoadFont reads a FIGlet font (.flf) file from the given path.
func LoadFont(path string) (*Font, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadFont(f)
}

// ReadFont reads a FIGlet font from an io.Reader.
func ReadFont(r io.Reader) (*Font, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	// Read header line
	if !scanner.Scan() {
		return nil, fmt.Errorf("figlet: empty font file")
	}
	header := scanner.Text()

	font, commentLines, codetagCount, err := parseHeader(header)
	if err != nil {
		return nil, err
	}

	// Skip comment lines
	for i := 0; i < commentLines; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("figlet: unexpected end of file in comments")
		}
	}

	// Read required characters
	for _, ch := range requiredChars {
		glyph, err := readGlyph(scanner, font.Height)
		if err != nil {
			break // EOF or error: use what we have
		}
		font.Glyphs[ch] = glyph
	}

	// Read code-tagged characters
	for i := 0; i < codetagCount; i++ {
		if !scanner.Scan() {
			break
		}
		tag := strings.TrimSpace(scanner.Text())
		if tag == "" {
			break
		}

		code, err := parseCodeTag(tag)
		if err != nil {
			continue
		}

		glyph, err := readGlyph(scanner, font.Height)
		if err != nil {
			break
		}
		font.Glyphs[rune(code)] = glyph
	}

	return font, nil
}

func parseHeader(header string) (*Font, int, int, error) {
	// Header format: flf2a<hardblank> height baseline maxlen old_layout comment_lines [print_dir full_layout codetag_count]
	if !strings.HasPrefix(header, "flf2a") {
		return nil, 0, 0, fmt.Errorf("figlet: not a FIGlet font (missing flf2a header)")
	}

	hardBlank := rune(header[5])
	parts := strings.Fields(header[6:])
	if len(parts) < 5 {
		return nil, 0, 0, fmt.Errorf("figlet: invalid header: not enough fields")
	}

	height, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, 0, 0, fmt.Errorf("figlet: invalid height: %w", err)
	}
	baseline, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, 0, 0, fmt.Errorf("figlet: invalid baseline: %w", err)
	}
	maxLen, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, 0, 0, fmt.Errorf("figlet: invalid maxlen: %w", err)
	}
	commentLines, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, 0, 0, fmt.Errorf("figlet: invalid comment lines: %w", err)
	}

	codetagCount := 0
	if len(parts) > 7 {
		codetagCount, _ = strconv.Atoi(parts[7])
	}

	font := &Font{
		Height:    height,
		Baseline:  baseline,
		MaxLen:    maxLen,
		HardBlank: hardBlank,
		Glyphs:    make(map[rune]*Glyph),
	}
	return font, commentLines, codetagCount, nil
}

func parseCodeTag(tag string) (int64, error) {
	tag = strings.Fields(tag)[0]
	if strings.HasPrefix(tag, "0x") || strings.HasPrefix(tag, "0X") {
		return strconv.ParseInt(tag[2:], 16, 64)
	}
	if strings.HasPrefix(tag, "0") && len(tag) > 1 {
		return strconv.ParseInt(tag[1:], 8, 64)
	}
	return strconv.ParseInt(tag, 10, 64)
}

func readGlyph(scanner *bufio.Scanner, height int) (*Glyph, error) {
	lines := make([]string, 0, height)
	width := 0
	for i := 0; i < height; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected end of glyph data")
		}
		line := scanner.Text()
		// Remove trailing @ markers
		// Last line ends with @@, others end with @
		line = strings.TrimRight(line, "@")
		if w := len([]rune(line)); w > width {
			width = w
		}
		lines = append(lines, line)
	}
	return &Glyph{Lines: lines, Width: width}, nil
}

// Render renders the given text using the font and returns the result as a string.
func (f *Font) Render(text string) string {
	var buf strings.Builder
	f.FprintTo(&buf, text)
	return buf.String()
}

// RenderWithColor renders the given text with colors cycling per character.
func (f *Font) RenderWithColor(text string, colors []Color) string {
	var buf strings.Builder
	f.fprintInternal(&buf, text, colors)
	return buf.String()
}

// FprintTo renders the given text using the font and writes the result to w.
func (f *Font) FprintTo(w io.Writer, text string) {
	f.fprintInternal(w, text, nil)
}

// FprintToWithColor renders the given text with colors and writes the result to w.
func (f *Font) FprintToWithColor(w io.Writer, text string, colors []Color) {
	f.fprintInternal(w, text, colors)
}

func (f *Font) fprintInternal(w io.Writer, text string, colors []Color) {
	runes := []rune(text)
	for row := 0; row < f.Height; row++ {
		for i, ch := range runes {
			glyph, ok := f.Glyphs[ch]
			if !ok {
				glyph, ok = f.Glyphs[' ']
				if !ok {
					continue
				}
			}
			line := ""
			if row < len(glyph.Lines) {
				line = glyph.Lines[row]
			}
			line = strings.ReplaceAll(line, string(f.HardBlank), " ")
			if len(colors) > 0 {
				c := colors[i%len(colors)]
				fmt.Fprint(w, c.prefix(), line, c.suffix())
			} else {
				fmt.Fprint(w, line)
			}
		}
		fmt.Fprintln(w)
	}
}

// Print renders the given text to stdout.
func (f *Font) Print(text string) {
	f.FprintTo(os.Stdout, text)
}

// PrintWithColor renders the given text to stdout with colors.
func (f *Font) PrintWithColor(text string, colors []Color) {
	f.FprintToWithColor(os.Stdout, text, colors)
}
