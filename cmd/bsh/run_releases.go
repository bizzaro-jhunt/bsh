package main

import (
	"fmt"
	"os"

	"github.com/jhunt/bsh/table"
)

func runReleases(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	releases, err := t.GetReleases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	if opt.AsJSON {
		jsonify(releases)
		os.Exit(0)
	}

	var tbl table.Table
	if opt.Releases.Jobs {
		tbl = table.NewTable("Name", "Version(s)", "Commit SHA1", "Jobs")
	} else {
		tbl = table.NewTable("Name", "Version(s)", "Commit SHA1")
	}
	for _, rel := range releases {
		for i := range rel.Versions {
			name := rel.Name
			deployed := " "
			if rel.Versions[i].Deployed {
				deployed = "*"
			}
			dirty := " "
			if rel.Versions[i].Dirty {
				dirty = "+"
			}

			if opt.Releases.Jobs {
				for j, job := range rel.Versions[i].Jobs {
					if j == 0 {
						tbl.Row(name,
							fmt.Sprintf("%s%s", rel.Versions[i].Version, deployed),
							fmt.Sprintf("%s%s", rel.Versions[i].Commit, dirty),
							job)
					} else {
						tbl.Row("", "", "", job)
					}
				}
				tbl.Row("", "", "")
			} else {
				if i != 0 {
					name = ""
				}
				tbl.Row(name,
					fmt.Sprintf("%s%s", rel.Versions[i].Version, deployed),
					fmt.Sprintf("%s%s", rel.Versions[i].Commit, dirty))
			}
		}
		tbl.Row("", "", "")
	}
	tbl.Print(os.Stdout)
}
