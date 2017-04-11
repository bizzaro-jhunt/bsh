package table

import (
	"io"
	"strings"

	fmt "github.com/jhunt/go-ansi"
)

type Table struct {
	Headers []string
	Widths  []int
	Rows    [][]string
	Prefix  string
}

func NewTable(headers ...string) Table {
	w := make([]int, len(headers))
	for i := range headers {
		w[i] = len(headers[i])
	}

	return Table{
		Headers: headers,
		Widths:  w,
		Rows:    make([][]string, 0),
	}
}

func (t *Table) Row(data ...interface{}) {
	row := make([]string, len(data))
	for i := range data {
		row[i] = fmt.Sprintf("%v", data[i])
		if len(row[i]) > t.Widths[i] {
			t.Widths[i] = len(row[i])
		}
	}

	t.Rows = append(t.Rows, row)
}

func (t Table) Print(out io.Writer) {
	fmt.Fprintf(out, t.Prefix)
	for i := range t.Headers {
		if i != 0 {
			fmt.Fprintf(out, "  ")
		}
		fmt.Fprintf(out, "@M{%-*s}", t.Widths[i], t.Headers[i])
	}
	fmt.Fprintf(out, "\n")

	fmt.Fprintf(out, t.Prefix)
	for i := range t.Widths {
		if i != 0 {
			fmt.Fprintf(out, "  ")
		}
		fmt.Fprintf(out, "%s", strings.Repeat("=", t.Widths[i]))
	}
	fmt.Fprintf(out, "\n")

	for i := range t.Rows {
		fmt.Fprintf(out, t.Prefix)
		for j := range t.Rows[i] {
			if j != 0 {
				fmt.Fprintf(out, "  ")
			}
			fmt.Fprintf(out, "%-*s", t.Widths[j], t.Rows[i][j])
		}
		fmt.Fprintf(out, "\n")
	}
}
