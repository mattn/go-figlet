# flf-bloody

Generates a rough "bloody" FIGlet font from an existing FLF file. The top edge uses `▄`, the body uses `█` and `▓`, and exposed lower edges can grow deterministic drips with `▓` and `▒`.

## Usage

```bash
flf-bloody input.flf > output.flf
```

`-drip-ratio` scales the maximum drip length relative to glyph height, so larger fonts naturally get longer drips. `-drip-density` controls how often exposed lower edges produce drips.

## Example

```bash
flf-bloody -drip-ratio 0.16 -drip-density 0.65 efont-b12_b.flf > bloody-efont-b12.flf
go-figlet -font bloody-efont-b12.flf "blood"
```
