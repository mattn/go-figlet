# flf-shadow

Generates an ANSI shadow-style FIGlet font from an existing FLF file. In double-width mode it uses box-drawing edges around solid blocks, and in single-width mode it keeps the same look while preserving more visible foreground width.

## Usage

```bash
flf-shadow [-double] input.flf > output.flf
```

## Flags

- `-double`
  Render each source pixel at double width. Default is single-width; pass `-double` to generate the wider variant.

## Example

```bash
flf-shadow efont-b12_b.flf > shadow-efont-b12.flf
go-figlet -font shadow-efont-b12.flf "figlet"
```

```
  ███╗    ██╗           ███╗            ██╗
 ██╔██╗   ╚═╝           ╚██║            ██║
 ██║╚═╝  ███╗   ████╗    ██║   ████╗  ██████╗
████╗    ╚██║  ██╔═██╗   ██║  ██╔═██╗ ╚═██╔═╝
╚██╔╝     ██║  ██║ ██║   ██║  █████╔╝   ██║
 ██║      ██║  ██║ ██║   ██║  ██╔══╝    ██║
 ██║     ████╗ ╚█████║  ████╗ ╚████╗    ╚███╗
 ╚═╝     ╚═══╝  ╚══██║  ╚═══╝  ╚═══╝     ╚══╝
                ████╔╝
                ╚═══╝
```
