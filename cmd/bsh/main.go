package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"
	cli "github.com/jhunt/go-cli"

	"github.com/jhunt/bsh/bosh"
	"github.com/jhunt/bsh/table"
)

const (
	OopsBadOptions int = iota
	OopsNotImplemented
	OopsBadConfiguration
	OopsCommunicationFailed
	OopsSaveConfigFailed
	OopsJSONFailed
)

type Opt struct {
	Help    bool `cli:"-h, --help"`
	Version bool `cli:"-v, --version"`

	URL      string `cli:"--director, --url"`
	Username string `cli:"-u, --username"`
	Password string `cli:"-p, --password"`
	CaCert   string `cli:"--ca-cert"`
	Insecure bool   `cli:"-k, --insecure, --no-insecure"`

	Config     string `cli:"-c, --config"`
	BOSHTarget string `cli:"-t, --target"`

	AsJSON bool `cli:"--json"`

	Task struct {
	} `cli:"task"`

	Tasks struct {
		States     []string `cli:"-s, --state"`
		Deployment string   `cli:"-d, --deployment"`
		ContextID  string   `cli:"-C, --context, --context-id"`
		Limit      int      `cli:"-l, --limit"`
	} `cli:"tasks"`

	Deployments struct {
	} `cli:"deployments"`

	Login struct {
	} `cli:"login"`

	Status struct {
	} `cli:"status"`

	VMs struct {
		Vitals     bool   `cli:"-V, --vitals"`
		Deployment string `cli:"-d, --deployment"`
	} `cli:"vms"`

	Targets struct {
	} `cli:"targets"`

	Target struct {
	} `cli:"target"`
}

func main() {
	var opt Opt
	opt.Config = fmt.Sprintf("%s/%s", os.Getenv("HOME"), ".boshrc")

	opt.Tasks.Limit = 30

	/* make sure ~/.boshrc exists... */
	if err := bosh.DefaultConfig(opt.Config); err != nil {
		fmt.Fprintf(os.Stderr, "%s: @Y{%s}\n", opt.Config, err)
	}

	command, args, err := cli.Parse(&opt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsBadOptions)
	}

	if command == "" {
		fmt.Fprintf(os.Stderr, "@R{a command is required...}\n")
		os.Exit(OopsBadOptions)
	}

	switch command {
	case "tasks":
		_, t := targeting(opt.Config)
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
		os.Exit(0)

	case "deployments":
		_, t := targeting(opt.Config)
		deployments, err := t.GetDeployments()
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsCommunicationFailed)
		}

		if opt.AsJSON {
			jsonify(deployments)
			os.Exit(0)
		}

		tbl := table.NewTable("Name", "Release(s)", "Stemcell(s)", "Cloud-Config")
		for _, d := range deployments {
			n := len(d.Releases)
			if len(d.Stemcells) > n {
				n = len(d.Stemcells)
			}

			for i := 0; i < n; i++ {
				var rel, stem string

				if i < len(d.Releases) {
					rel = fmt.Sprintf("%s/%s", d.Releases[i].Name, d.Releases[i].Version)
				}
				if i < len(d.Stemcells) {
					stem = fmt.Sprintf("%s/%s", d.Stemcells[i].Name, d.Stemcells[i].Version)
				}
				if i == 0 {
					tbl.Row(d.Name, rel, stem, d.CloudConfig)
				} else {
					tbl.Row("", rel, stem, "")
				}
			}
			tbl.Row("", "", "", "")
		}
		tbl.Print(os.Stdout)
		os.Exit(0)

	case "login":
		cfg, t := targeting(opt.Config)
		user := opt.Username
		if user == "" {
			user = prompt("Username: ", false)
		}

		pass := opt.Password
		if pass == "" {
			pass = prompt("Password: ", true)
		}

		t.Username = user
		t.Password = pass
		err = cfg.Save(opt.Config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsSaveConfigFailed)
		}
		os.Exit(0)

	case "status":
		_, t := targeting(opt.Config)
		info, err := t.GetInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsCommunicationFailed)
		}

		fmt.Printf("@G{Config}\n")
		fmt.Printf("  %-10s @Y{%s}\n", "", opt.Config)
		fmt.Printf("\n")
		fmt.Printf("@G{Director}\n")
		fmt.Printf("  %-10s @Y{%s}\n", "Name", info.Name)
		fmt.Printf("  %-10s @Y{%s}\n", "URL", t.URL)
		fmt.Printf("  %-10s @Y{%s}\n", "Version", info.Version)
		fmt.Printf("  %-10s @Y{%s}\n", "User", t.Username)
		fmt.Printf("  %-10s @Y{%s}\n", "UUID", info.UUID)
		fmt.Printf("  %-10s @Y{%s}\n", "CPI", info.CPI)
		os.Exit(0)

	case "vms":
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
		os.Exit(0)

	case "targets":
		cfg := readConfigFrom(opt.Config)
		tbl := table.NewTable("Alias", "Name", "URL", "UUID")
		for _, t := range cfg.Targets {
			tbl.Row(t.Alias, t.Name, t.URL, t.UUID)
		}
		tbl.Print(os.Stdout)
		os.Exit(0)

	case "target":
		cfg := readConfigFrom(opt.Config)
		if len(args) == 0 {
			if cfg.Current == "" {
				fmt.Printf("no default BOSH target has been set\n")
				os.Exit(0)
			}
			t, err := cfg.CurrentTarget()
			if err != nil {
				fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
				os.Exit(OopsBadConfiguration)
			}
			fmt.Printf("currently targeting @C{%s} at @G{%s}\n", t.Alias, t.URL)
			os.Exit(0)
		}

		if len(args) == 1 {
			t, ok := cfg.Targets[args[0]]
			if !ok {
				fmt.Fprintf(os.Stderr, "no such target @C{%s} in %s\n", args[0], opt.Config)
				os.Exit(OopsBadConfiguration)
			}
			cfg.Current = t.Alias
			err = cfg.Save(opt.Config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
				os.Exit(OopsSaveConfigFailed)
			}
			os.Exit(0)
		}

		if len(args) == 2 {
			t := bosh.Target{
				URL:      args[0],
				Alias:    args[1],
				Insecure: opt.Insecure,
			}

			err = t.Sync()
			if err != nil {
				fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
				os.Exit(OopsCommunicationFailed)
			}

			cfg.Targets[t.Alias] = &t
			cfg.Current = t.Alias
			err = cfg.Save(opt.Config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
				os.Exit(OopsSaveConfigFailed)
			}
			os.Exit(0)
		}

		fmt.Fprintf(os.Stderr, "@Y{incorrect usage}\n")
		os.Exit(OopsBadOptions)

	default:
		fmt.Fprintf(os.Stderr, "%s - @*{not yet implemented...}\n", command)
		os.Exit(OopsNotImplemented)
	}
}
