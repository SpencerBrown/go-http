// FILEPATH: /Users/admin/proj/go-http/util/util_test.go

package util

import (
	"strings"
	"testing"
)

func TestIndent(t *testing.T) {
	// Test case 1: Indent with positive indentLevel
	input1 := "Hello\nWorld\n"
	expected1 := "  Hello\n  World\n"
	output1 := Indent(input1, 1)
	if output1 != expected1 {
		t.Errorf("Indent(%q, %d) = %q, expected %q", input1, 2, output1, expected1)
	}

	// Test case 2: Indent with zero indentLevel
	input2 := "Hello\nWorld\n"
	expected2 := "Hello\nWorld\n"
	output2 := Indent(input2, 0)
	if output2 != expected2 {
		t.Errorf("Indent(%q, %d) = %q, expected %q", input2, 0, output2, expected2)
	}

	// Test case 3: Indent with negative indentLevel
	input3 := "Hello\nWorld\n"
	expected3 := "Hello\nWorld\n"
	output3 := Indent(input3, -2)
	if output3 != expected3 {
		t.Errorf("Indent(%q, %d) = %q, expected %q", input3, -2, output3, expected3)
	}

	// Test case 4: Indent with indentLevel > 20
	input4 := "Hello\nWorld\n"
	indent40 := strings.Repeat("  ", 20)
	expected4 := indent40 + "Hello\n" + indent40 + "World\n"
	output4 := Indent(input4, 25)
	if output4 != expected4 {
		t.Errorf("Indent(%q, %d) = %q, expected %q", input4, 25, output4, expected4)
	}

	// Test case 5: Indent with empty string
	input5 := ""
	expected5 := ""
	output5 := Indent(input5, 4)
	if output5 != expected5 {
		t.Errorf("Indent(%q, %d) = %q, expected %q", input5, 4, output5, expected5)
	}
}
