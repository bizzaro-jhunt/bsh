package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func runDeleteVM(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	res, err := t.Delete(fmt.Sprintf("/vms/%s", args[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsBadOptions)
	}

	watch(t, res, Done{
		Good: "vm deleted successfully",
		Bad:  "vm deletion failed",
	})
}
