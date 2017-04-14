package main

import (
	fmt "github.com/jhunt/go-ansi"
	"io/ioutil"
	"os"

	"github.com/jhunt/bsh/bosh"
	"github.com/jhunt/bsh/query"
)

func runDeploy(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}
	if opt.Deployment == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	manifest, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	diff, err := t.Diff(opt.Deployment, manifest, opt.Deploy.Redact)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	bosh.FormatDiff(os.Stdout, diff)
	fmt.Printf("\n")
	fmt.Printf("@M{%s} > @B{%s}\n", t.Name, opt.Deployment)
	if !confirm(fmt.Sprintf("Deploy these changes? [@G{yes}/@R{no}] ")) {
		fmt.Printf("@R{aborting...}\n")
		os.Exit(OopsCancelled)
	}

	q := query.New()
	q.Set("context", `{"cloud_config_id":null,"runtime_config_id":null}`)
	q.Bool("recreate", opt.Deploy.Recreate)
	r, err := t.PostYAML(fmt.Sprintf("/deployments%s", q), manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	watch(t, r, okfail("deploy"))
}
