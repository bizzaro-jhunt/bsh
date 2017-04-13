package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/query"
)

func runDeleteDeployment(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)

	name := opt.Delete.Deployment.Deployment
	if name == "" {
		if len(args) == 1 {
			name = args[0]
		}
	}
	if name == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	q := query.New()
	q.Bool("force", opt.Delete.Force)

	res, err := t.Delete(fmt.Sprintf("/deployments/%s%s", name, q))
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	watch(t, res, Done{
		Good: "deployment deleted successfully",
		Bad:  "deployment deletion failed",
	})
}
