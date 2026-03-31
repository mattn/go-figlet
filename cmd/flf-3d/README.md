# flf-3d

Generates 3D shadow FIGlet fonts from existing FLF files. Each filled pixel becomes `█` with a `░` shadow at the lower-left.

## Usage

```bash
flf-3d [-double] input.flf > output.flf
```

## Flags

- `-double`
  Render each source pixel at double width. Default is single-width; pass `-double` to generate the wider variant.

## Example

```bash
flf-3d efont-b12_b.flf > 3d-efont-b12.flf
go-figlet -font 3d-efont-b12.flf "figlet"
```

```
  ███    ██          ███           ██
 ██░██  ░░          ░░██          ░██
░██░░   ███   ████   ░██   ████  ██████
████   ░░██  ██░░██  ░██  ██░░██ ░░██░
░██     ░██  ██ ░██  ░██  █████   ░██
░██     ░██  ██ ░██  ░██  ██░░    ░██
░██     ████ ░█████  ████ ░████   ░░███
░░     ░░░░  ░░░░██ ░░░░  ░░░░     ░░░
             ████
             ░░░░
```
