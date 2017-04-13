package main

import (
	fmt "github.com/jhunt/go-ansi"
	"io/ioutil"
	"os"

	"github.com/jhunt/bsh/bosh"
)

func runDiff(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}
	if opt.Diff.Deployment == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	manifest, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	diff, err := t.Diff(opt.Diff.Deployment, manifest, opt.Diff.Redact)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	bosh.FormatDiff(os.Stdout, diff)
	fmt.Printf("\n")
}
