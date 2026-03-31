# flf-bloody

Generates a rough "bloody" FIGlet font from an existing FLF file. The top edge uses `▄`, the body uses `█` and `▓`, and exposed lower edges can grow short deterministic drips with `▓` and `▒`.

## Usage

```bash
flf-bloody input.flf > output.flf
```

## Example

```bash
flf-bloody efont-b12_b.flf > bloody-efont-b12.flf
go-figlet -font bloody-efont-b12.flf "blood"
```
