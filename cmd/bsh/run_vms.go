package main

import (
	"fmt"
	"os"

	"github.com/jhunt/bsh/bosh"
	"github.com/jhunt/bsh/table"
)

func runVMs(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	var deployments []string

	if len(args) > 0 {
		deployments = append(deployments, args...)
	} else if opt.VMs.Deployment != "" {
		deployments = append(deployments, opt.VMs.Deployment)
	}
	if len(deployments) == 0 {
		deploys, err := t.GetDeployments()
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsCommunicationFailed)
		}
		for _, d := range deploys {
			deployments = append(deployments, d.Name)
		}
	}

	if opt.AsJSON {
		m := make(map[string][]bosh.VM)
		for _, deployment := range deployments {
			vms, err := t.GetVMsFor(deployment)
			if err != nil {
				fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
				os.Exit(OopsCommunicationFailed)
			}
			m[deployment] = vms
		}
		jsonify(m)
		os.Exit(0)
	}

	for n, deployment := range deployments {
		vms, err := t.GetVMsFor(deployment)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsCommunicationFailed)
		}

		tbl := table.NewTable("Name", "State", "UUID", "AZ", "Type", "IPs")
		tbl.Prefix = "   "
		for _, vm := range vms {
			for i := range vm.IPs {
				if i == 0 {
					tbl.Row(vm.FullName(), vm.JobState, vm.CID, or(vm.AZ, "-"), vm.Type, vm.IPs[i])
				} else {

					tbl.Row("", "", "", "", "", vm.IPs[i])
				}
			}
		}
		if n != 0 {
			fmt.Printf("\n\n\n")
		}
		fmt.Printf("@G{%s} vms:\n\n", deployment)
		tbl.Print(os.Stdout)
	}
}
