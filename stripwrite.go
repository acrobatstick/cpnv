package main

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
)

type offset struct {
	// starting offset of the line
	startLine int
	// variable value start offset
	start int
	// variable value end offset
	end int
}

func (o *offset) decrease(n int) {
	o.startLine -= n
	o.start -= n
	o.end -= n
}

func stripwrite(r io.Reader, w io.Writer, keep bool, excludedVars []string) error {
	// the starting offset position of a variable
	offsets := make([]offset, 0)
	excludedOffsets := make([]offset, 0)

	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	// current reading position on the buffer
	var cursor int
	// amount to decrease for the next variable offset
	decreaseNext := 0
	lineNumber := 1

	// to track the offset of a comment or a multi-line comments
	commentStart := -1

match:
	for cursor < len(buf) {
		var line strings.Builder
		startLine := cursor
		for cursor < len(buf) && buf[cursor] != '\n' {
			line.WriteByte(buf[cursor])
			cursor++
		}
		end := cursor

		// consume newline
		if cursor < len(buf) && buf[cursor] == '\n' {
			cursor++
			lineNumber++
		}
		linestr := line.String()
		if len(linestr) == 0 {
			commentStart = -1
			continue
		}

		if linestr[0] == '#' {
			if commentStart == -1 {
				commentStart = startLine
			}
			continue
		}

		key, _, found := strings.Cut(linestr, "=")
		if !found {
			return fmt.Errorf("environment variable is not valid at line %d\n", lineNumber)
		}

		start := startLine + len(key) + 1 // + 1 to consume the equal sign

		// check for excluded variable patterns if they match with
		// the scanned keys
		for _, pat := range excludedVars {
			matched, err := path.Match(pat, string(key))
			if err != nil {
				return fmt.Errorf("invalid glob pattern %q: %w", pat, err)
			}
			if matched {
				excludedStart := startLine
				if commentStart != -1 {
					excludedStart = commentStart
				}
				excludedOffset := offset{excludedStart, excludedStart, cursor}
				excludedOffsets = append(excludedOffsets, excludedOffset)
				decreaseNext += cursor - excludedStart

				// since we already found the matching variable. we want
				// to start over to find the new comment for the next variable
				commentStart = -1

				continue match
			}
		}

		offset := offset{startLine, start, end}
		offset.decrease(decreaseNext)
		offsets = append(offsets, offset)
		commentStart = -1
	}

	// remove excluded variables (entire lines including comments)
	for i := len(excludedOffsets) - 1; i >= 0; i-- {
		o := excludedOffsets[i]
		buf = append(buf[:o.start], buf[o.end:]...)
	}

	// remove the environment value and update the offsets
	if !keep {
		for i := len(offsets) - 1; i >= 0; i-- {
			o := offsets[i]
			buf = append(buf[:o.start], buf[o.end:]...)
		}
	}

	_, err = w.Write(bytes.TrimLeft(buf, "\n"))
	return err
}
