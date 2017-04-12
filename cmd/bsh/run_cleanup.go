package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func runCleanup(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	task, err := t.Cleanup(opt.Cleanup.All)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	err = t.Follow(os.Stdout, task.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsTaskFailed)
	}
	fmt.Printf("@G{cleanup complete.}\n")
}
