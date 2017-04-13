package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/query"
)

func runDeleteStemcell(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)

	l := make([][]string, 0)
	for len(args) > 0 {
		rel, ver, rest, err := thingversion(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
			os.Exit(OopsBadOptions)
		}

		l = append(l, []string{rel, ver})
		args = rest
	}

	if len(l) == 0 {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	for _, p := range l {
		q := query.New()
		q.Bool("force", opt.Delete.Force)

		fmt.Printf("@R{deleting} stemcell @B{%s}/@M{%s}...\n", p[0], p[1])
		res, err := t.Delete(fmt.Sprintf("/stemcells/%s/%s%s", p[0], p[1], q))
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsCommunicationFailed)
		}

		watch(t, res, Done{
			Good: "stemcell deleted successfully",
			Bad:  "stemcell deletion failed",
		})
	}
}
