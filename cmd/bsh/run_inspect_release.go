package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/bsh/table"
)

func runInspectRelease(opt Opt, command string, args []string) {
	_, t := targeting(opt)

	l := make([][]string, 0)
	for len(args) > 0 {
		rel, ver, rest, err := thingversion(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
			os.Exit(OopsBadOptions)
		}

		l = append(l, []string{rel, ver})
		args = rest
	}

	if len(l) == 0 {
		fmt.Fprintf(os.Stderr, "@R{!!! usage...}\n")
		os.Exit(OopsBadOptions)
	}

	ok := true
	for _, p := range l {
		release, err := t.GetRelease(p[0], p[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s/%s: %s}\n", p[0], p[1], err)
			ok = false
			continue
		}

		tbl := table.NewTable("Job", "Fingerprint", "Blobstore ID", "SHA1")
		for _, job := range release.Jobs {
			tbl.Row(job.Name, job.Fingerprint, job.BlobstoreID, job.SHA1)
		}

		fmt.Printf("@G{%s}/@G{%s} jobs:\n", p[0], p[1])
		tbl.Print(os.Stdout)
		fmt.Printf("\n\n")

		tbl = table.NewTable("Package", "Fingerprint", "Compiled For", "Blobstore ID", "SHA1")
		for _, pkg := range release.Packages {
			tbl.Row(pkg.Name, pkg.Fingerprint, "(source)", pkg.BlobstoreID, pkg.SHA1)
			for _, cp := range pkg.CompiledPackages {
				tbl.Row("", "", cp.Stemcell, cp.BlobstoreID, cp.SHA1)
			}
			tbl.Row("", "", "", "", "")
		}

		fmt.Printf("@G{%s}/@G{%s} packages:\n", p[0], p[1])
		tbl.Print(os.Stdout)
		fmt.Printf("\n\n")
	}

	if !ok {
		os.Exit(OopsCommunicationFailed)
	}
}
