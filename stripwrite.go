package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func stripwrite(r io.Reader, w io.Writer, keep bool) error {
	scanner := bufio.NewScanner(r)

	var line int
	for scanner.Scan() {
		line++
		b := scanner.Bytes()

		// check if line only consists of whitespace or an empty line
		if len(bytes.TrimSpace(b)) == 0 {
			w.Write([]byte("\n\n"))
			continue
		}

		// write comments directly
		if b[0] == '#' {
			w.Write(b)
			w.Write([]byte("\n"))
			continue
		}
		key, after, found := bytes.Cut(b, []byte("="))
		if !found {
			return fmt.Errorf("environment variable is not valid at line %d", line)
		}

		w.Write(key)
		w.Write([]byte("="))

		if keep {
			w.Write(after)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
