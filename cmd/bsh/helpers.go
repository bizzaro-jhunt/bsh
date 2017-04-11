package main

import (
	"encoding/json"
	fmt "github.com/jhunt/go-ansi"
	"os"
	"time"

	"github.com/jhunt/bsh/bosh"
)

func readConfigFrom(path string) bosh.Config {
	cfg, err := bosh.ReadConfig(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read bsh configuration from %s:\n@R{!!! %s}\n", path, err)
		os.Exit(OopsBadConfiguration)
	}
	return cfg
}

func targeting(path string) (bosh.Config, *bosh.Target) {
	cfg := readConfigFrom(path)
	t, err := cfg.CurrentTarget()
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsBadConfiguration)
	}
	return cfg, t
}

func tstamp(ts int) string {
	t := time.Unix(int64(ts), 0)
	return t.Format("2006-01-02 15:04:05-0700 MST")
}

func or(iff, els string) string {
	if iff != "" {
		return iff
	}
	return els
}

func jsonify(x interface{}) {
	b, err := json.Marshal(x)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsJSONFailed)
	}
	fmt.Printf("%s\n", string(b))
}
