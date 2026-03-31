# flf-bloody

Generates a rough "bloody" FIGlet font from an existing FLF file. The top edge uses `▄`, the body uses `█` and `▓`, and exposed lower edges can grow deterministic drips with `▓` and `▒`.

## Usage

```bash
flf-bloody [-double] [-drip-ratio=0.12] [-drip-density=0.5] input.flf > output.flf
```

## Flags

- `-double`
  Render each source pixel at double width. Default is single-width.
- `-drip-ratio`
  Scale the maximum drip length relative to glyph height. Larger fonts naturally get longer drips.
- `-drip-density`
  Control how often exposed lower edges produce drips. Use lower values for a cleaner look and higher values for heavier bleeding.

## Example

```bash
flf-bloody -drip-ratio 0.16 -drip-density 0.65 efont-b12_b.flf > bloody-efont-b12.flf
go-figlet -font bloody-efont-b12.flf "blood"
```

```
▄▄      ▓▄▄                    ▄▄
▓▓      ▒▓▓                    ▓▓
▓███▄    ▓▓   ▄██▄   ▄██▄   ▄███▓
▓▓ ▒▓▄   ▓▓  ▄▓ ▒▓▄ ▄▓ ▒▓▄ ▄▓ ▒▓▓
▓▓  ▓▓   ▓▓  ▓▓  ▓▓ ▓▓  ▓▓ ▓▓  ▓▓
▓▓  ▓▓   ▓▓  ▓▓  ▓▓ ▓▓  ▓▓ ▓▓  ▓▓
▓███▓▒  ▓██▓  ▓██▓▒  ▓██▓▒  ▓███▓
 ▒▒▒    ▒▒▒   ▒▒▒    ▒▒▒    ▒▒▒
```
