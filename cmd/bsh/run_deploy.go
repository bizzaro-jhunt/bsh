package main

import (
	fmt "github.com/jhunt/go-ansi"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/jhunt/bsh/bosh"
)

func runDeploy(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}
	if opt.Deploy.Deployment == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	manifest, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	diff, err := t.Diff(opt.Deploy.Deployment, manifest, opt.Deploy.Redact)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	bosh.FormatDiff(os.Stdout, diff)
	fmt.Printf("\n")
	fmt.Printf("@M{%s} > @B{%s}\n", t.Name, opt.Deploy.Deployment)
	if !confirm(fmt.Sprintf("Deploy these changes? [@G{yes}/@R{no}] ")) {
		fmt.Printf("@R{aborting...}\n")
		os.Exit(OopsCancelled)
	}

	qs := "?context=" + url.QueryEscape(`{"cloud_config_id":null,"runtime_config_id":null}`)
	if opt.Deploy.Recreate {
		qs += "&recreate=true"
	}
	r, err := t.PostYAML("/deployments"+qs, manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	var task bosh.Task
	err = t.InterpretJSON(r, &task)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	err = t.Follow(os.Stdout, task.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	task, err = t.GetTask(task.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	if !task.Succeeded() {
		fmt.Fprintf(os.Stderr, "@R{task failed.}\n")
		os.Exit(OopsTaskFailed)
	}
}
