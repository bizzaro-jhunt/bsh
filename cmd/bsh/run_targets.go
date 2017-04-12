package main

import (
	"os"

	"github.com/jhunt/bsh/table"
)

func runTargets(opt Opt, command string, args []string) {
	cfg := readConfigFrom(opt.Config)
	tbl := table.NewTable("Alias", "Name", "URL", "UUID")
	for _, t := range cfg.Targets {
		tbl.Row(t.Alias, t.Name, t.URL, t.UUID)
	}
	tbl.Print(os.Stdout)
}
