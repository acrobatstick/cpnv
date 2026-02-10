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

func TestStripValueNoSpaceBetween(t *testing.T) {
	input := `VAR_1=FOO
VAR_2=BAR`

	expected := `VAR_1=
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
}

func TestStripWriteWithExclude(t *testing.T) {
	input := `# Variable 1
VAR_1=FOO

# Variable 2
# Variable comment
VAR_2=BAR

LAST_KEY=GOOD BYE FRIENDS :'(
# Butt cheeks
# Booty hole
BUTT=CHEEKS`
	expected := `LAST_KEY=
# Butt cheeks
# Booty hole
BUTT=`
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
