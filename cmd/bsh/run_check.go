package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/bosh"
)

func runCheck(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)

	if opt.Check.Deployment == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	r, err := t.Post(fmt.Sprintf("/deployments/%s/scans", opt.Check.Deployment), nil)
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
