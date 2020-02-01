package goldga

import (
	"strings"

	"github.com/andreyvit/diff"
	"github.com/logrusorgru/aurora"
)

func diffString(color bool, snapshot, received string) string {
	lines := []string{
		"- Snapshot",
		"+ Received",
		"",
	}
	lines = append(lines, diff.LineDiffAsLines(snapshot, received)...)

	if color {
		for i, line := range lines {
			if len(line) == 0 {
				continue
			}

			switch line[0] {
			case '+':
				lines[i] = aurora.BrightGreen(line).String()
			case '-':
				lines[i] = aurora.BrightRed(line).String()
			default:
				lines[i] = aurora.BrightBlack(line).String()
			}
		}
	}

	return strings.Join(lines, "\n")
}
