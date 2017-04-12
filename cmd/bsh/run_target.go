package main

import (
	"fmt"
	"os"

	"github.com/jhunt/bsh/bosh"
)

func runTarget(opt Opt, command string, args []string) {
	cfg := readConfigFrom(opt.Config)
	if len(args) == 0 {
		if cfg.Current == "" {
			fmt.Printf("no default BOSH target has been set\n")
			os.Exit(0)
		}
		t, err := cfg.CurrentTarget()
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsBadConfiguration)
		}
		fmt.Printf("currently targeting @C{%s} at @G{%s}\n", t.Alias, t.URL)
		os.Exit(0)
	}

	if len(args) == 1 {
		t, ok := cfg.Targets[args[0]]
		if !ok {
			fmt.Fprintf(os.Stderr, "no such target @C{%s} in %s\n", args[0], opt.Config)
			os.Exit(OopsBadConfiguration)
		}
		cfg.Current = t.Alias
		err := cfg.Save(opt.Config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsSaveConfigFailed)
		}
		os.Exit(0)
	}

	if len(args) == 2 {
		t := bosh.Target{
			URL:      args[0],
			Alias:    args[1],
			Insecure: opt.Insecure,
		}

		err := t.Sync()
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsCommunicationFailed)
		}

		cfg.Targets[t.Alias] = &t
		cfg.Current = t.Alias
		err = cfg.Save(opt.Config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsSaveConfigFailed)
		}
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "@Y{incorrect usage}\n")
	os.Exit(OopsBadOptions)
}
