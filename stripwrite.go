package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path"
)

func stripwrite(r io.Reader, w io.Writer, keep bool, excludedVars []string) error {
	out := [][]byte{}
	scanner := bufio.NewScanner(r)
	var currLine int

scan:
	for scanner.Scan() {
		currLine++
		b := scanner.Bytes()

		// check if line only consists of whitespace or an empty line
		if len(bytes.TrimSpace(b)) == 0 {
			out = append(out, []byte{})
			continue
		}

		// write comments directly
		if b[0] == '#' {
			comment := append(b, '\n')
			out = append(out, comment)
			continue
		}
		key, value, found := bytes.Cut(b, []byte("="))
		if !found {
			return fmt.Errorf("environment variable is not valid at line %d\n", currLine)
		}

		// check for excluded variable patterns if they match with
		// the scanned keys
		for _, pat := range excludedVars {
			matched, err := path.Match(pat, string(key))
			if err != nil {
				return fmt.Errorf("invalid glob pattern %q: %w", pat, err)
			}

			if !matched {
				continue
			}

			for i := len(out) - 1; i >= 0; i-- {
				prev := out[i]
				// if prev output line was a variable. skip trimming out slice
				if len(prev) > 0 && prev[0] != '#' {
					break
				}
				// if the line is a comment, then move every slice index
				// after i to be i-1
				out = append(out[:i], out[i+1:]...)
			}

			continue scan
		}

		line := []byte{}
		line = append(line, key...)
		line = append(line, '=')

		if keep {
			line = append(line, value...)
		}

		out = append(out, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, line := range out {
		if len(line) == 0 {
			w.Write([]byte("\n\n"))
			continue
		}
		w.Write(line)
	}

	return nil
}
