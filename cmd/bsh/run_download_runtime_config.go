package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/bosh"
)

func runDownloadRuntimeConfig(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	var l []bosh.RuntimeConfig
	err := t.GetJSON("/runtime_configs?limit=1", &l)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	config := ""
	if len(l) == 1 {
		config = l[0].Properties
	}

	if opt.AsJSON {
		jsonify(yamlr(config))
		os.Exit(0)
	}

	downloadto(os.Stdout, config, opt.Download.Output, opt.Download.Force)
}
