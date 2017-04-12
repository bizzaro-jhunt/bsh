package bosh

import (
	fmt "github.com/jhunt/go-ansi"
	"io"
)

type ManifestDiff struct {
	Diff [][]string `json:"diff"`
}

func (t Target) Diff(deployment string, manifest []byte, redact bool) ([][]string, error) {
	qs := "?redact=false"
	if redact {
		qs = "?redact=true"
	}
	r, err := t.PostYAML(fmt.Sprintf("/deployments/%s/diff%s", deployment, qs), manifest)
	if err != nil {
		return nil, err
	}

	var diff ManifestDiff
	err = t.InterpretJSON(r, &diff)
	if err != nil {
		return nil, err
	}

	return diff.Diff, nil
}

func FormatDiff(out io.Writer, diff [][]string) {
	for _, delta := range diff {
		if len(delta) != 2 {
			continue
		}
		switch delta[1] {
		case "added":
			fmt.Fprintf(out, "@G{%s}\n", delta[0])
		case "removed":
			fmt.Fprintf(out, "@R{%s}\n", delta[0])
		default:
			fmt.Fprintf(out, "%s\n", delta[0])
		}
	}
}
