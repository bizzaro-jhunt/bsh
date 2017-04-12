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

	Type    string            `json:"type"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`

	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
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

		if ev.Type == "deprecation" {
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

		} else if ev.State == "failed" {
			if last.Stage == ev.Stage && last.Task == ev.Task {
				fmt.Fprintf(out, ". @R{FAILED}\n      !!! %s", ev.Data["error"])
				ev.State = "failx"
			} else {
				if last.State != ev.State {
					fmt.Fprintf(out, "\n")
				}
				fmt.Fprintf(out, "\n   @R{FAILED %s} > @Y{%s}.\n      !!! %s",
					strings.ToLower(ev.Stage), strings.ToLower(ev.Task),
					ev.Data["error"])
			}

		} else if ev.Error.Code > 0 {
			if last.State != ev.State {
				fmt.Fprintf(out, "\n")
			}
			fmt.Fprintf(out, "\n   @R{OOPS: %s} (error @Y{%d})", ev.Error.Message, ev.Error.Code)

		} else {
			panic(sc.Text())
		}

		last = ev
	}
	if n > 0 {
		fmt.Fprintf(out, "\n\n")
	}
	return sc.Err()
}
