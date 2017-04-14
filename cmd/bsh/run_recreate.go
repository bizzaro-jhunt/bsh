package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/bosh"
	"github.com/jhunt/bsh/query"
)

func runRecreate(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	if opt.Deployment == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

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

	q := query.New()
	q.Set("state", "recreate")
	q.Bool("skip_drain", opt.Recreate.SkipDrain)
	q.Bool("force", opt.Recreate.Force)
	q.Bool("fix", opt.Start.Fix)
	q.Bool("dry_run", opt.Recreate.DryRun)
	if opt.Recreate.Canaries != 0 {
		q.Set("canaries", fmt.Sprintf("%d", opt.Recreate.Canaries))
	}
	if opt.Recreate.MaxInFlight != 0 {
		q.Set("max_in_flight", fmt.Sprintf("%d", opt.Recreate.MaxInFlight))
	}

	name := opt.Deployment
	ok := true
	for _, p := range l {
		fmt.Printf("recreating @C{%s} @G{%s}/@G{%s}...\n", name, p[0], p[1])
		res, err := t.PutYAMLish(fmt.Sprintf("/deployments/%s/jobs/%s/%s%s", name, p[0], p[1], q))
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			ok = false
			continue
		}

		var task bosh.Task
		err = t.InterpretJSON(res, &task)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsJSONFailed)
		}

		follow(t, task.ID, Done{
			Good: fmt.Sprintf("%s/%s recreated successfully", p[0], p[1]),
			Bad:  fmt.Sprintf("%s/%s failed to recreate", p[0], p[1]),
		}, true)
	}

	if !ok {
		os.Exit(OopsCommunicationFailed)
	}
}
