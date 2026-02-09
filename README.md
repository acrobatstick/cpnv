# cpnv

`cpnv` is a command-line tool to easily duplicate your .env file

## Installation

```bash
go install github.com/acrobatstick/cpnv@latest
```

## Usage

```bash
Usage: cpnv [OPTIONS]

Copy environment variable file.

Options:
  -input, -i    Input file path (default: ".env")
  -output, -o   Output file path (default: "copy")
  -keep, -k     Keep original values (default: false)
  -help, -h     Show this help message

Examples:
  cpnv -input .env -output processed.env
  cpnv -i input.env -o output.env
```
