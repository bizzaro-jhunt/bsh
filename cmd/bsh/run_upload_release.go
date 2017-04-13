package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/query"
)

func runUploadRelease(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadConfiguration)
	}

	if opt.Upload.Name != "" && opt.Upload.Version != "" {
		q := query.New()
		q.Set("version", opt.Upload.Version)
		res, err := t.Get(fmt.Sprintf("/releases%s", q))
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsCommunicationFailed)
		}

		if res.StatusCode == 200 {
			fmt.Printf("release @C{%s}/@B{%s} already exists; skipping...\n",
				opt.Upload.Name, opt.Upload.Version)
			os.Exit(0)
		}
	}

	q := query.New()
	q.Bool("fix", opt.Upload.Fix)
	q.Bool("rebase", opt.Upload.Release.Rebase)
	q.Bool("skip_if_exists", opt.Upload.SkipIfExists)
	q.Maybe("sha1", opt.Upload.SHA1)

	body, n, err := upload(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	fmt.Printf("uploading release @C{%s}...\n", args[0])
	res, err := t.PostFile(fmt.Sprintf("/releases%s", q), body, n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsCommunicationFailed)
	}

	watch(t, res, okfail("release upload"))
}
