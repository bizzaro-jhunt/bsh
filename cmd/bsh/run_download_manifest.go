package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func runDownloadManifest(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)

	name := opt.Download.Deployment
	if name == "" {
		if len(args) == 1 {
			name = args[0]
		}
	}
	if name == "" {
		fmt.Fprintf(os.Stderr, "@R{!!! usage}\n")
		os.Exit(OopsCommunicationFailed)
	}

	deployment, err := t.GetDeployment(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	if opt.AsJSON {
		jsonify(yamlr(deployment.Manifest))
		os.Exit(0)
	}

	downloadto(os.Stdout, deployment.Manifest, opt.Download.Output, opt.Download.Force)
}
