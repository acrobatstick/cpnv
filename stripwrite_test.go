package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestStripWrite(t *testing.T) {
	input := `# Variable 1
VAR_1=FOO

# Variable 2
# Variable comment
VAR_2=BAR`

	expected := `# Variable 1
VAR_1=

# Variable 2
# Variable comment
VAR_2=`
	var output bytes.Buffer
	reader := strings.NewReader(input)

	if err := stripwrite(reader, &output, false, []string{}); err != nil {
		t.Fatal(err)
	}

	got := output.String()
	if got != expected {
		t.Fatalf("stripWrite() output mismatch\ngot:\n%s\nexpected:\n%s\n", got, expected)
	}

	// test with keep original value
	output = bytes.Buffer{}
	reader = strings.NewReader(input)
	if err := stripwrite(reader, &output, true, []string{}); err != nil {
		t.Fatal(err)
	}

	got = output.String()
	if got != input {
		t.Fatalf("stripWrite() output mismatch\ngot:\n%s\nexpected:\n%s\n", got, expected)
	}
}

func TestStripWriteWithExclude(t *testing.T) {
	input := `# Variable 1
VAR_1=FOO

# Variable 2
# Variable comment
VAR_2=BAR

LAST_KEY=NOW IM ALL ALONE`
	expected := `

LAST_KEY=`
	var output bytes.Buffer
	reader := strings.NewReader(input)

	if err := stripwrite(reader, &output, false, []string{"VAR_*"}); err != nil {
		t.Fatal(err)
	}

	got := output.String()
	if got != expected {
		t.Fatalf("stripWrite() output mismatch\ngot:\n%s\nexpected:\n%s\n", got, expected)
	}
}
