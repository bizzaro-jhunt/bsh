package main

import (
	"fmt"
	"os"
)

func runLogin(opt Opt, command string, args []string) {
	cfg, t := targeting(opt.Config)
	user := opt.Username
	if user == "" {
		user = prompt("Username: ", false)
	}

	pass := opt.Password
	if pass == "" {
		pass = prompt("Password: ", true)
	}

	t.Username = user
	t.Password = pass
	err := cfg.Save(opt.Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsSaveConfigFailed)
	}
}
