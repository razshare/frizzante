package assemblies

import (
	"regexp"
	"strings"
)

var AssemblyLinePattern = regexp.MustCompile(`^(\s+)(0x[0-9a-f]+)(\s+)([0-9a-f]+)(\s+)(.+)$`)

// FormatAssemblyWithAlignment aligns the instruction column (last column) in assembly output.
// It preserves Go source code lines (non-assembly lines) without modification.
func FormatAssemblyWithAlignment(content string) string {
	lines := strings.Split(content, "\n")
	maxHexWidth := 0
	for _, line := range lines {
		matches := AssemblyLinePattern.FindStringSubmatch(line)
		if matches != nil {
			// matches[4] is the hex bytes group
			hexBytes := matches[4]
			if len(hexBytes) > maxHexWidth {
				maxHexWidth = len(hexBytes)
			}
		}
	}
	var builder strings.Builder
	for i, line := range lines {
		matches := AssemblyLinePattern.FindStringSubmatch(line)
		if matches != nil {
			leadingSpace := matches[1]
			address := matches[2]
			hexBytes := matches[4]
			instruction := matches[6]
			paddedHexBytes := hexBytes + strings.Repeat(" ", maxHexWidth-len(hexBytes))
			builder.WriteString(leadingSpace)
			builder.WriteString(address)
			builder.WriteString("      ")
			builder.WriteString(paddedHexBytes)
			builder.WriteString("      ")
			builder.WriteString(instruction)
		} else {
			// Non-assembly line (Go source code reference), keep as-is
			builder.WriteString(line)
		}
		if i < len(lines)-1 || line != "" {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}
