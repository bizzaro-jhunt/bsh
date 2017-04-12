package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/table"
)

func runErrands(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	if opt.Errands.Deployment == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	errands, err := t.GetErrandsFor(opt.Errands.Deployment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	if opt.AsJSON {
		jsonify(errands)
		os.Exit(0)
	}

	tbl := table.NewTable("Name")
	for _, errand := range errands {
		tbl.Row(errand.Name)
	}
	tbl.Print(os.Stdout)
}
