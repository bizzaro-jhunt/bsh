package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func runCleanup(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	task, err := t.Cleanup(opt.Cleanup.All)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}
	follow(t, task.ID, okfail("cleanup"), true)
}
