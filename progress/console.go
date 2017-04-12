package progress

import (
	fmt "github.com/jhunt/go-ansi"
	"io"
)

func Console(out io.Writer, min, max int, before, after string, progress rune) func(bool, int64, int64, bool) {
	if after == "" {
		after = "\n"
	}

	var saved, width int64
	extra := len(before) + len(after)

	if n, err := TerminalWidth(); err != nil || n < min {
		width = int64(min - extra)
	} else if n > max {
		width = int64(max - extra)
	} else {
		width = int64(n - extra)
	}

	return func(start bool, n int64, total int64, end bool) {
		if start {
			fmt.Fprintf(out, before)
			return
		}
		have := width * saved / total
		want := width * n / total
		for want > have {
			fmt.Fprintf(out, "%c", progress)
			have++
		}
		saved = n

		if end {
			fmt.Fprintf(out, after)
		}
	}
}
