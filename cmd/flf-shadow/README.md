# flf-shadow

Generates an ANSI shadow-style FIGlet font from an existing FLF file. In double-width mode it uses box-drawing edges around solid blocks, and in single-width mode it keeps the same look while preserving more visible foreground width.

## Usage

```
flf-shadow input.flf > output.flf
```

## Example

```bash
flf-shadow efont-b12_b.flf > shadow-efont-b12.flf
go-figlet -font shadow-efont-b12.flf "figlet"
```

```
█████▓▓        ███▓▓                 █████▓▓                   ███▓▓
███  ██▓▓      ███▓▓                 ██ ███▓▓                  ███▓▓
███  ███▓▓  █████▓▓▓▓   █████▓▓        ███▓▓    █████▓▓   █████████▓▓
███████▓▓   ███  ███▓▓ ███  ███▓▓      ███▓▓   ███  ███▓▓  ▓▓▓███▓▓
███  ███▓▓  ███  ███▓▓ ███  ███▓▓      ███▓▓   ████████▓▓    ███▓▓
███  ███▓▓  ███  ███▓▓ ███  ███▓▓      ███▓▓   ███▓▓▓▓▓▓     ███▓▓
███████▓▓   ████████▓▓  ███████▓▓    ███████▓▓  ███████▓▓    █████▓▓
▓▓▓▓▓▓      ▓▓▓▓▓▓██▓▓   ▓▓▓▓▓▓▓▓    ▓▓▓▓▓▓▓▓    ▓▓▓▓▓▓▓▓     ▓▓▓▓▓▓
            ███████▓▓
            ▓▓▓▓▓▓▓▓
```
