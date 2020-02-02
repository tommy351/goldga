package goldga

import (
	"os"
	"strings"

	"github.com/andreyvit/diff"
	"github.com/logrusorgru/aurora"
	isatty "github.com/mattn/go-isatty"
)

// nolint: gochecknoglobals
var (
	DefaultDiffer Differ = &ColorDiffer{}

	colorSupported = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
)

type Differ interface {
	Diff(snapshot, received []byte) []byte
}

var _ Differ = (*ColorDiffer)(nil)

type ColorDiffer struct{}

func (ColorDiffer) Diff(snapshot, received []byte) []byte {
	lines := []string{
		"- Snapshot",
		"+ Received",
		"",
	}
	lines = append(lines, diff.LineDiffAsLines(string(snapshot), string(received))...)

	if colorSupported {
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

	return []byte(strings.Join(lines, "\n"))
}
