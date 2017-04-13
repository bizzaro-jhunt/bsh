package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/bosh"
	"github.com/jhunt/bsh/table"
)

func runTasks(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	if !opt.Tasks.All && len(opt.Tasks.States) == 0 {
		opt.Tasks.States = append(opt.Tasks.States, "running")
	}
	tasks, err := t.GetTasks(bosh.TasksFilter{
		States:     opt.Tasks.States,
		Deployment: opt.Tasks.Deployment,
		ContextID:  opt.Tasks.ContextID,
		Limit:      opt.Tasks.Limit,
		Verbose:    2,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	tbl := table.NewTable("ID", "State", "Started", "Last Activity", "User", "Deployment", "Description", "Result")
	for _, task := range tasks {
		result := "(none)"
		if task.Result != nil {
			result = *task.Result
		}
		tbl.Row(task.ID, task.State,
			tstamp(task.StartedAt), tstamp(task.Timestamp),
			task.User, task.Deployment, task.Description, result)
	}
	tbl.Print(os.Stdout)
}
