package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/bosh"
	"github.com/jhunt/bsh/query"
)

func runStart(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	if opt.Start.Soft && opt.Start.Hard {
		fmt.Fprintf(os.Stderr, "@R{!!! cannot specify both --soft and --hard}\n")
		os.Exit(OopsBadOptions)
	}

	if opt.Start.Deployment == "" {
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
	q.Set("state", "started")
	q.Bool("skip_drain", opt.Start.SkipDrain)
	q.Bool("force", opt.Start.Force)
	q.Bool("fix", opt.Start.Fix)
	q.Bool("dry_run", opt.Start.DryRun)
	if opt.Start.Canaries != 0 {
		q.Set("canaries", fmt.Sprintf("%d", opt.Start.Canaries))
	}
	if opt.Start.MaxInFlight != 0 {
		q.Set("max_in_flight", fmt.Sprintf("%d", opt.Start.MaxInFlight))
	}

	name := opt.Start.Deployment
	ok := true
	for _, p := range l {
		fmt.Printf("starting @C{%s} @G{%s}/@G{%s}...\n", name, p[0], p[1])
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
			Good: fmt.Sprintf("%s/%s started successfully", p[0], p[1]),
			Bad:  fmt.Sprintf("%s/%s failed to start", p[0], p[1]),
		}, true)
	}

	if !ok {
		os.Exit(OopsCommunicationFailed)
	}
}
