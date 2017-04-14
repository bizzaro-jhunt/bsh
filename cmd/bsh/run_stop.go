package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/bosh"
	"github.com/jhunt/bsh/query"
)

func runStop(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	if opt.Stop.Soft && opt.Stop.Hard {
		fmt.Fprintf(os.Stderr, "@R{!!! cannot specify both --soft and --hard}\n")
		os.Exit(OopsBadOptions)
	}

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
	if opt.Stop.Hard {
		q.Set("state", "detached")
	} else {
		q.Set("state", "stopped")
	}
	q.Bool("skip_drain", opt.Stop.SkipDrain)
	q.Bool("force", opt.Stop.Force)
	q.Bool("dry_run", opt.Stop.DryRun)
	if opt.Stop.Canaries != 0 {
		q.Set("canaries", fmt.Sprintf("%d", opt.Stop.Canaries))
	}
	if opt.Stop.MaxInFlight != 0 {
		q.Set("max_in_flight", fmt.Sprintf("%d", opt.Stop.MaxInFlight))
	}

	name := opt.Deployment
	ok := true
	for _, p := range l {
		if opt.Stop.Hard {
			fmt.Printf("stopping @C{%s} @G{%s}/@G{%s} and deleting its VM...\n", name, p[0], p[1])
		} else {
			fmt.Printf("stopping @C{%s} @G{%s}/@G{%s}...\n", name, p[0], p[1])
		}
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
			Good: fmt.Sprintf("%s/%s stopped successfully", p[0], p[1]),
			Bad:  fmt.Sprintf("%s/%s failed to stop", p[0], p[1]),
		}, true)
	}

	if !ok {
		os.Exit(OopsCommunicationFailed)
	}
}
