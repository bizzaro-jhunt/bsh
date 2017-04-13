package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/table"
)

func runStemcells(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	stemcells, err := t.GetStemcells()
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	if opt.AsJSON {
		jsonify(stemcells)
		os.Exit(0)
	}

	tbl := table.NewTable("Name", "Version(s)", "OS", "CID")
	for _, stem := range stemcells {
		deployed := " "
		if len(stem.Deployments) > 0 {
			deployed = "*"
		}
		tbl.Row(stem.Name, stem.OS,
			fmt.Sprintf("%s%s", stem.Version, deployed), stem.CID)
	}
	tbl.Print(os.Stdout)
}
