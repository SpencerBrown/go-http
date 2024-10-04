package util

import "strings"

func Indent(s string, indentSpaces int) string {
	lines := strings.Split(s, "\n")
	var out, indent strings.Builder
	for i := 0; i < indentSpaces; i++ {
		indent.WriteString("  ")
	}
	for i, line := range lines {
		if len(line) > 0 {
			out.WriteString(indent.String())
			out.WriteString(line)
			if i < len(lines)-1 {
				out.WriteByte('\n')
			}
		}
	}
	return out.String()
}
