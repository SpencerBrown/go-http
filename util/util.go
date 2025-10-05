package util

import "strings"

// Indent takes a multi-line string, where the lines are delimited by newline characters,
// and indents it by the given number of levels. Each level is two spaces.
// This is useful for pretty-printing nested information. 
// A negative level is treated as zero, and a level >20 is treated as 20.
// It returns the indented string.
// If the input string is empty, it returns an empty string.
// If the input string has a trailing newline, the output will have a trailing newline.
func Indent(s string, indentLevel int) string {
	if indentLevel < 0 {
		indentLevel = 0
	}
	if indentLevel > 20 { // 20 levels is the maximum
		indentLevel = 20
	}
	lines := strings.Split(s, "\n")
	var out strings.Builder
	var indent = strings.Repeat("  ", indentLevel)
	for i, line := range lines {
		if !(i == len(lines)-1 && line == "") {
			// special case: last line ignored if empty, because
			// we don't want to add an extra newline at the end
			out.WriteString(indent)
			out.WriteString(line)
			out.WriteByte('\n')
		}
	}
	return out.String()
}
