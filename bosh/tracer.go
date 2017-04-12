package bosh

import (
	"bufio"
	"encoding/json"
	fmt "github.com/jhunt/go-ansi"
	"io"
	"strings"
)

type Event struct {
	Timestamp int    `json:"timestamp"`
	Stage     string `json:"stage"`
	Task      string `json:"task"`
	Index     int    `json:"index"`
	Total     int    `json:"total"`
	State     string `json:"state"`
	Progress  int    `json:"progress"`
}

func TraceEvents(out io.Writer, in io.Reader) error {
	var last Event
	n := 0

	sc := bufio.NewScanner(in)
	for sc.Scan() {
		var ev Event
		err := json.Unmarshal(sc.Bytes(), &ev)
		if err != nil {
			fmt.Printf("@R{!!! %s}\n", err)
			continue
		}

		n++
		if ev.State == "started" {
			if last.Stage != ev.Stage || last.Task != ev.Task || last.State != ev.State {
				fmt.Fprintf(out, "\n")
			}
			fmt.Fprintf(out, "  @Y{Started} %s > @G{%s}",
				strings.ToLower(ev.Stage), ev.Task)

		} else if ev.State == "finished" {
			if last.Stage == ev.Stage && last.Task == ev.Task {
				fmt.Fprintf(out, ". @B{Done} (...)")
				ev.State = "finx"
			} else {
				if last.State != ev.State {
					fmt.Fprintf(out, "\n")
				}
				fmt.Fprintf(out, "\n     @B{Done} %s > @G{%s}. (...)",
					strings.ToLower(ev.Stage), strings.ToLower(ev.Task))
			}
		}
		last = ev
	}
	if n > 0 {
		fmt.Fprintf(out, "\n\n")
	}
	return sc.Err()
}
