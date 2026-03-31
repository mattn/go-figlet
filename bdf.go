package figlet

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type bdfChar struct {
	encoding rune
	width    int
	bitmap   string
	printed  bool
}

// BDF2FLF converts a BDF font from r and writes FIGlet font data to w.
func BDF2FLF(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var fontName, copyright string
	var width, height, baseline int
	var numChars int
	var isJIS0208 bool

	// Parse BDF headers
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch strings.ToUpper(fields[0]) {
		case "FONTBOUNDINGBOX":
			if len(fields) < 5 {
				return fmt.Errorf("bdf2flf: invalid FONTBOUNDINGBOX")
			}
			width, _ = strconv.Atoi(fields[1])
			height, _ = strconv.Atoi(fields[2])
			yoff, _ := strconv.Atoi(fields[4])
			baseline = height + yoff
		case "FONT":
			if len(fields) > 1 {
				fontName = strings.Join(fields[1:], " ")
				if strings.Contains(fontName, "JISX0208") {
					isJIS0208 = true
				}
			}
		case "COPYRIGHT":
			if idx := strings.Index(line, " "); idx >= 0 {
				copyright = strings.TrimSpace(line[idx+1:])
			}
		case "CHARS":
			if len(fields) > 1 {
				numChars, _ = strconv.Atoi(fields[1])
			}
			goto readChars
		}
	}

readChars:
	if numChars <= 0 || width <= 0 || height <= 0 {
		return fmt.Errorf("bdf2flf: could not extract BDF data")
	}
	if width > 32 || height > 32 {
		return fmt.Errorf("bdf2flf: font too large (max 32x32)")
	}

	chars := make([]bdfChar, 0, numChars)
	var cur bdfChar
	inBitmap := false

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		keyword := strings.ToUpper(fields[0])
		switch {
		case keyword == "ENCODING" && len(fields) > 1:
			enc, _ := strconv.Atoi(fields[1])
			cur.encoding = rune(enc)
		case keyword == "BBX" && len(fields) > 1:
			cur.width, _ = strconv.Atoi(fields[1])
		case keyword == "BITMAP":
			inBitmap = true
			cur.bitmap = ""
		case keyword == "ENDCHAR":
			inBitmap = false
			chars = append(chars, cur)
			cur = bdfChar{}
		case inBitmap:
			cur.bitmap += strings.TrimSpace(line)
		}
	}

	// Convert JIS X 0208 encodings to Unicode
	if isJIS0208 {
		decoder := japanese.EUCJP.NewDecoder()
		for i := range chars {
			jis := uint(chars[i].encoding)
			euc := []byte{byte((jis >> 8) | 0x80), byte((jis & 0xFF) | 0x80)}
			utf8bytes, err := io.ReadAll(transform.NewReader(
				strings.NewReader(string(euc)), decoder))
			if err == nil && len(utf8bytes) > 0 {
				runes := []rune(string(utf8bytes))
				if len(runes) == 1 {
					chars[i].encoding = runes[0]
				}
			}
			decoder.Reset()
		}
	}

	// Build required set for fast lookup
	requiredSet := make(map[rune]bool, len(requiredChars))
	for _, r := range requiredChars {
		requiredSet[r] = true
	}

	// Count code-tagged (non-required) characters
	codetags := 0
	for _, ch := range chars {
		if !requiredSet[ch.encoding] {
			codetags++
		}
	}

	// Write FLF header
	comments := 3
	if copyright != "" {
		comments = 4
	}
	fmt.Fprintf(w, "flf2a$ %d %d %d -1 %d 0 16384 %d\n",
		height, baseline, width+2, comments, codetags)
	fmt.Fprintf(w, "FIGlet font derived from %s\n", fontName)
	if copyright != "" {
		fmt.Fprintf(w, "which has copyright: %s\n", copyright)
	}
	fmt.Fprintf(w, "Converted to flf with go-bdf2flf.\n\n")

	// Write required characters first
	for _, req := range requiredChars {
		found := false
		for i := range chars {
			if chars[i].encoding == req {
				writeCharacter(w, chars[i].bitmap, chars[i].width, height)
				chars[i].printed = true
				found = true
				break
			}
		}
		if !found {
			writeBlankCharacter(w, width, height)
		}
	}

	// Write additional characters with code tags
	for i := range chars {
		if !chars[i].printed {
			fmt.Fprintf(w, "%d\n", chars[i].encoding)
			writeCharacter(w, chars[i].bitmap, chars[i].width, height)
		}
	}

	return nil
}

func writeCharacter(w io.Writer, bitmap string, width, height int) {
	hexChars := (width + 7) / 8 * 2
	pos := 0
	for y := 0; y < height; y++ {
		if y > 0 {
			fmt.Fprint(w, "\n")
		}
		// Extract hex string for this row
		hex := ""
		for x := 0; x < hexChars && pos < len(bitmap); x++ {
			hex += string(bitmap[pos])
			pos++
		}
		// Pad if needed
		for len(hex) < hexChars {
			hex += "0"
		}

		val, _ := strconv.ParseUint(hex, 16, 64)
		bit := uint64(1) << (uint(hexChars)*4 - 1)
		for x := 0; x < width; x++ {
			if val&bit != 0 {
				fmt.Fprint(w, "#")
			} else {
				fmt.Fprint(w, " ")
			}
			if bit == 1 {
				break
			}
			bit >>= 1
		}
		fmt.Fprint(w, "$@")
	}
	fmt.Fprint(w, "@\n")
}

func writeBlankCharacter(w io.Writer, width, height int) {
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

