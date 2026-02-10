package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [OPTIONS]

Clone environment variable files.

Options:
  -input, -i			Input file path (default: ".env")
  -output, -o			Output file path (default: "copy")
  -keep, -k			Keep original values (default: false)
  -exclude, -x <PATTERN>...	Exclude environment variables
				Examples:
				 -x KEY_1 KEY_2
				 -x "DB_*" -x "KEY_*"
  -help, -h			Show this help message

Examples:
  %s -input .env -output processed.env
  %s -i input.env -o output.env

`, os.Args[0], os.Args[0], os.Args[0])
	os.Exit(1)
}

func main() {
	// default values for the flags
	input := ".env"
	output := "copy"
	keep := false

	if len(os.Args) < 2 {
		usage()
	}

	excludedVars := []string{}

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-input", "-i":
			if i+1 < len(os.Args) {
				input = os.Args[i+1]
				i++
			}
		case "-output", "-o":
			if i+1 < len(os.Args) {
				output = os.Args[i+1]
				i++
			}
		case "-keep", "-k":
			keep = true
			i++
		case "-exclude", "-x":
			for i+1 < len(os.Args) {
				next := os.Args[i+1]
				// check if the next argument is a flag or value
				if next[0] == '-' {
					break
				}
				excludedVars = append(excludedVars, next)
				i++
			}
		case "-help", "-h":
			usage()
		}
	}

	in, err := os.Open(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	out, err := os.Create(fmt.Sprintf(".env-%s", output))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer out.Close()

	err = stripwrite(in, out, keep, excludedVars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error stripping original environment file: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
