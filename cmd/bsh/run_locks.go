package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/table"
)

func runLocks(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	locks, err := t.GetLocks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	if opt.AsJSON {
		jsonify(locks)
		os.Exit(0)
	}

	tbl := table.NewTable("Type", "Resource", "Timeout")
	for _, lock := range locks {
		tbl.Row(lock.Type, lock.Resource, lock.Timeout)
	}
	tbl.Print(os.Stdout)
}
