package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/bosh"
	"github.com/jhunt/bsh/query"
)

func runUploadStemcell(opt Opt, command string, args []string) {
	_, t := targeting(opt.Config)
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadConfiguration)
	}

	if opt.Upload.Name != "" && opt.Upload.Version != "" {
		var l []bosh.Stemcell
		err := t.GetJSON("/stemcells", &l)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsCommunicationFailed)
		}

		for _, sc := range l {
			if sc.Name == opt.Upload.Name && sc.Version == opt.Upload.Version {
				fmt.Printf("stemcell @C{%s}/@B{%s} already exists; skipping...\n",
					opt.Upload.Name, opt.Upload.Version)
				os.Exit(0)
			}
		}
	}

	q := query.New()
	q.Bool("fix", opt.Upload.Fix)
	q.Bool("skip_if_exists", opt.Upload.SkipIfExists)
	q.Maybe("sha1", opt.Upload.SHA1)

	body, n, err := upload(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	fmt.Printf("uploading stemcell @C{%s}...\n", args[0])
	res, err := t.PostFile(fmt.Sprintf("/stemcells%s", q), body, n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	watch(t, res, okfail("stemcell upload"))
}
