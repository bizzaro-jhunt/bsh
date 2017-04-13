package main

import (
	fmt "github.com/jhunt/go-ansi"
	"io"
	"os"
)

func runCurl(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadConfiguration)
	}

	r, err := t.Get(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	io.Copy(os.Stdout, r.Body)
}
