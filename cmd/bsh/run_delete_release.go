package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/query"
)

func runDeleteRelease(opt Opt, command string, args []string) {
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
		q.Set("version", p[1])
		q.Bool("force", opt.Delete.Release.Force)

		fmt.Printf("@R{deleting} release @B{%s}/@M{%s}...\n", p[0], p[1])
		res, err := t.Delete(fmt.Sprintf("/releases/%s%s", p[0], q))
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsCommunicationFailed)
		}

		watch(t, res, Done{
			Good: "release deleted successfully",
			Bad:  "release deletion failed",
		})
	}
}
