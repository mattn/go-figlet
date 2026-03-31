# go-figlet

FIGlet implementation in Go. Render text using FIGlet fonts (.flf), with color support.

Also includes `go-bdf2flf`, a tool to convert BDF bitmap fonts to FIGlet font format.

## Installation

```
go install github.com/mattn/go-figlet/cmd/go-figlet@latest
go install github.com/mattn/go-figlet/cmd/go-bdf2flf@latest
```

## Usage

### go-figlet

```
go-figlet [-font font.flf] [-color colors] text...
```

Arguments can also be read from stdin:

```
echo "Hello" | go-figlet
```

#### Color options

```
# Presets: rainbow, warm, cool
go-figlet -color rainbow "Hello"

# Named colors (semicolon-separated)
go-figlet -color "red;green;blue" "Hello"

# Hex colors
go-figlet -color "FF0000;00FF00;0000FF" "Hello"
```

Colors cycle per character.

### go-bdf2flf

Convert BDF font to FIGlet font:

```
go-bdf2flf < input.bdf > output.flf
```

Supports JIS X 0208 encoded BDF fonts (auto-detected and converted to Unicode).

## Library

```go
package main

import (
	"embed"
	"os"

	"github.com/mattn/go-figlet"
)

//go:embed fonts/*.flf
var fonts embed.FS

func main() {
	font, _ := figlet.LoadFontFS(fonts, "fonts/standard.flf")

	// Simple rendering
	font.Print("Hello")

	// With colors
	font.PrintWithColor("Hello", []figlet.Color{
		figlet.ColorRed,
		figlet.ColorGreen,
		figlet.ColorBlue,
	})

	// True color
	font.PrintWithColor("Hello", []figlet.Color{
		figlet.NewTrueColor(255, 100, 0),
		figlet.NewTrueColor(0, 200, 255),
	})

	// Render to string
	s := font.Render("Hello")
	os.Stdout.WriteString(s)

	// Convert BDF to FLF
	bdf, _ := os.Open("input.bdf")
	flf, _ := os.Create("output.flf")
	figlet.BDF2FLF(bdf, flf)
}
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
