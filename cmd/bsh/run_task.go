package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
	"strconv"
)

func runTask(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "@R{usage...}\n")
		os.Exit(OopsBadOptions)
	}

	id, err := strconv.ParseUint(args[0], 10, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsBadOptions)
	}

	if opt.Task.Event {
		follow(t, int(id), okfail("task"), false)
		fmt.Printf("\n")
	}

	task, _ := t.GetTask(int(id))
	if !opt.Task.Event && !opt.Task.Debug && !opt.Task.CPI && !opt.Task.Result {
		follow(t, int(id), okfail("task"), false)
		s, _ := t.GetTaskOutput(task, "result")
		fmt.Printf("\n%s\n\n", s)
	}

	if opt.Task.Debug {
		fmt.Printf("\n@Y{DEBUG OUTPUT}\n============\n")
		s, err := t.GetTaskOutput(task, "debug")
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{usage...}\n")
		}
		fmt.Printf("%s\n\n", s)
	}

	if opt.Task.CPI {
		fmt.Printf("\n@Y{CPI LOG}\n======\n")
		s, err := t.GetTaskOutput(task, "cpi")
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{usage...}\n")
		}
		fmt.Printf("%s\n\n", s)
	}

	if opt.Task.Result {
		fmt.Printf("\n@Y{RESULT}\n======\n")
		s, err := t.GetTaskOutput(task, "result")
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{usage...}\n")
		}
		fmt.Printf("%s\n\n", s)
	}

	if !task.Succeeded() {
		os.Exit(OopsTaskFailed)
	}
}
