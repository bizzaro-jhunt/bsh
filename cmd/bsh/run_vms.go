package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/bosh"
	"github.com/jhunt/bsh/table"
)

func runVMs(opt Opt, command string, args []string) {
	_, t := targeting(opt)
	var deployments []string

	if len(args) > 0 {
		deployments = append(deployments, args...)
	} else if opt.Deployment != "" {
		deployments = append(deployments, opt.Deployment)
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

		// Instance | State | AZ | VM Type | IPs | VM CID | Disk CID | Agent ID | Resurrection | Ignore | Uptime | Load | CPU % | CPU % | Memory Usage | Swap Usage | System | Ephemeral | Persistent |
		var tbl table.Table
		if opt.VMs.Details {
			tbl = detailedTable(vms)
		} else if opt.VMs.Vitals {
			tbl = vitalsTable(vms)
		} else if opt.VMs.Processes {
			tbl = processTable(vms)
		} else {
			tbl = summaryTable(vms)
		}

		if n != 0 {
			fmt.Printf("\n\n\n")
		}
		fmt.Printf("@G{%s} vms:\n\n", deployment)
		tbl.Print(os.Stdout)
	}
}

func summaryTable(vms []bosh.VM) table.Table {
	tbl := table.NewTable("Name", "State", "UUID", "AZ", "Type", "IPs")
	tbl.Prefix = "   "
	for _, vm := range vms {
		tbl.Row(
			vm.FullName(),
			vm.JobState,
			vm.ID,
			or(vm.AZ, "-"),
			vm.Type,
			vm.IPs)
	}
	return tbl
}

func detailedTable(vms []bosh.VM) table.Table {
	tbl := table.NewTable("Name", "State", "UUID", "AZ", "Type", "IPs", "CIDs")
	tbl.Prefix = "   "
	for _, vm := range vms {
		tbl.Row(
			vm.FullName(),
			[]string{
				vm.JobState,
				vm.ResurrectionState(),
				vm.IgnoredState(),
			},
			vm.ID,
			or(vm.AZ, "-"),
			vm.Type,
			vm.IPs,
			[]string{
				fmt.Sprintf("agent: %s", or(vm.AgentID, "(none)")),
				fmt.Sprintf("   vm: %s", or(vm.CID, "(none)")),
				fmt.Sprintf(" disk: %s", or(vm.DiskCID, "(none)")),
			},
		)
	}
	return tbl
}

func vitalsTable(vms []bosh.VM) table.Table {
	tbl := table.NewTable("Name", "State", "UUID", "AZ", "Type", "IPs", "Load", "CPU", "Memory", "Disk")
	tbl.Prefix = "   "
	for _, vm := range vms {
		tbl.Row(
			vm.FullName(),
			[]string{
				vm.JobState,
				vm.ResurrectionState(),
				vm.IgnoredState(),
			},
			vm.ID,
			or(vm.AZ, "-"),
			vm.Type,
			vm.IPs,
			[]string{
				fmt.Sprintf(" 1m: %s", vm.Vitals.Load[0]),
				fmt.Sprintf(" 5m: %s", vm.Vitals.Load[1]),
				fmt.Sprintf("15m: %s", vm.Vitals.Load[2]),
			},
			[]string{
				fmt.Sprintf(" sys %s", vm.Vitals.CPU.Sys),
				fmt.Sprintf("user %s", vm.Vitals.CPU.User),
				fmt.Sprintf("wait %s", vm.Vitals.CPU.Wait),
			},
			[]string{
				fmt.Sprintf("ram %s (%s)", kbsize(vm.Vitals.Memory.KB), percent(vm.Vitals.Memory.Percent)),
				fmt.Sprintf("swp %s (%s)", kbsize(vm.Vitals.Swap.KB), percent(vm.Vitals.Swap.Percent)),
			},
			[]string{
				fmt.Sprintf("sys %s (i:%s)", percent(vm.Vitals.Disk.System.Percent), percent(vm.Vitals.Disk.System.InodePercent)),
				fmt.Sprintf("eph %s (i:%s)", percent(vm.Vitals.Disk.Ephemeral.Percent), percent(vm.Vitals.Disk.Ephemeral.InodePercent)),
				fmt.Sprintf("prs %s (i:%s)", percent(vm.Vitals.Disk.Persistent.Percent), percent(vm.Vitals.Disk.Persistent.InodePercent)),
			},
		)
		tbl.Spacer()
	}
	return tbl
}

func processTable(vms []bosh.VM) table.Table {
	tbl := table.NewTable("Name", "State", "UUID", "AZ", "Type", "IPs", "Uptime", "CPU", "Memory")
	tbl.Prefix = "   "
	for _, vm := range vms {
		tbl.Row(
			vm.FullName(),
			[]string{
				vm.JobState,
				vm.ResurrectionState(),
				vm.IgnoredState(),
			},
			vm.ID,
			or(vm.AZ, "-"),
			vm.Type,
			vm.IPs,
			"",
			[]string{
				fmt.Sprintf(" sys %s%%", vm.Vitals.CPU.Sys),
				fmt.Sprintf("user %s%%", vm.Vitals.CPU.User),
				fmt.Sprintf("wait %s%%", vm.Vitals.CPU.Wait),
			},
			[]string{
				fmt.Sprintf("ram %s (%s)", kbsize(vm.Vitals.Memory.KB), percent(vm.Vitals.Memory.Percent)),
				fmt.Sprintf("swp %s (%s)", kbsize(vm.Vitals.Swap.KB), percent(vm.Vitals.Swap.Percent)),
			},
		)

		for _, proc := range vm.Processes {
			tbl.Row(
				fmt.Sprintf(" -- %s", proc.Name),
				proc.State,
				"",
				"",
				"",
				"",
				clocked(proc.Uptime.Seconds),
				fmt.Sprintf("     %3.1f%%", proc.CPU.Total),
				fmt.Sprintf("    %s (%3.1f%%)", kbsize(fmt.Sprintf("%d", proc.Memory.KB)), proc.Memory.Percent),
			)
		}
		tbl.Spacer()
	}
	return tbl
}
