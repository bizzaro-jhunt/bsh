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

	follow(t, int(id), okfail("task"))
}
