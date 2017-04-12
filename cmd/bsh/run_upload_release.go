package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func runUploadRelease(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadConfiguration)
	}

	body, n, err := upload(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	fmt.Printf("uploading release @C{%s}...\n", args[0])
	res, err := t.PostFile("/releases", body, n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	watch(t, res, okfail("release upload"))
}
