package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/table"
)

func runDeployments(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	deployments, err := t.GetDeployments()
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	if opt.AsJSON {
		jsonify(deployments)
		os.Exit(0)
	}

	tbl := table.NewTable("Name", "Release(s)", "Stemcell(s)", "Cloud-Config")
	for _, d := range deployments {
		n := len(d.Releases)
		if len(d.Stemcells) > n {
			n = len(d.Stemcells)
		}

		for i := 0; i < n; i++ {
			var rel, stem string

			if i < len(d.Releases) {
				rel = fmt.Sprintf("%s/%s", d.Releases[i].Name, d.Releases[i].Version)
			}
			if i < len(d.Stemcells) {
				stem = fmt.Sprintf("%s/%s", d.Stemcells[i].Name, d.Stemcells[i].Version)
			}
			if i == 0 {
				tbl.Row(d.Name, rel, stem, d.CloudConfig)
			} else {
				tbl.Row("", rel, stem, "")
			}
		}
		tbl.Row("", "", "", "")
	}
	tbl.Print(os.Stdout)
}
