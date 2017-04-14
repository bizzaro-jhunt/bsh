package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func runCheck(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	if opt.Deployment == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	r, err := t.Post(fmt.Sprintf("/deployments/%s/scans", opt.Deployment), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	watch(t, r, okfail("cloud check"))
}
