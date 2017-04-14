package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func runLogout(opt Opt, command string, args []string) {
	cfg, t := targeting(opt)

	t.Username = ""
	t.Password = ""
	err := cfg.Save(opt.Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsSaveConfigFailed)
	}
}
