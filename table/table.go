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

func (t *Table) row(data ...interface{}) {
	row := make([]string, len(data))
	for i := range data {
		row[i] = fmt.Sprintf("%v", data[i])
		if len(row[i]) > t.Widths[i] {
			t.Widths[i] = len(row[i])
		}
	}

	t.Rows = append(t.Rows, row)
}

func (t *Table) Row(raw ...interface{}) {
	h := 1

	for i, v := range raw {
		if l, ok := v.([]interface{}); ok {
			if len(l) > h {
				h = len(l)
			}
		} else if l, ok := v.([]string); ok {
			if len(l) > h {
				h = len(l)
			}
		} else {
			raw[i] = []interface{}{v}
		}
	}

	data := make([][]interface{}, h)
	for y := 0; y < h; y++ {
		data[y] = make([]interface{}, len(raw))
		for x := range raw {
			data[y][x] = ""
		}
	}
	for x, v := range raw {
		if l, ok := v.([]interface{}); ok {
			for y, w := range l {
				data[y][x] = fmt.Sprintf("%v", w)
			}
		} else if l, ok := v.([]string); ok {
			for y, w := range l {
				data[y][x] = fmt.Sprintf("%v", w)
			}
		}
	}

	for _, row := range data {
		t.row(row...)
	}
}

func (t *Table) Spacer() {
	filler := make([]interface{}, len(t.Headers))
	for i := range filler {
		filler[i] = ""
	}
	t.row(filler...)
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
