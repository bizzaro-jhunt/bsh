package main

import (
	"encoding/json"
	fmt "github.com/jhunt/go-ansi"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jhunt/bsh/bosh"
	"github.com/jhunt/bsh/progress"
)

type Done struct {
	Good string
	Bad  string
}

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

func thingversion(args []string) (string, string, []string, error) {
	if len(args) < 1 {
		return "", "", nil, nil
	}

	if strings.Contains(args[0], "/") {
		l := strings.SplitN(args[0], "/", 2)
		return l[0], l[1], args[1:], nil
	}

	if len(args) >= 2 {
		if strings.Contains(args[1], "/") {
			return "", "", nil, fmt.Errorf("no version provided for '%s'", args[0])
		}
		return args[0], args[1], args[2:], nil
	}

	return "", "", nil, fmt.Errorf("no version provided for '%s'", args[0])
}

func okfail(typ string) Done {
	return Done{
		Good: fmt.Sprintf("%s finished successfully", typ),
		Bad:  fmt.Sprintf("%s failed", typ),
	}
}

func follow(t *bosh.Target, id int, done Done, bail bool) {
	fmt.Printf("bosh task @G{%d}\n", id)
	err := t.Follow(os.Stdout, id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsJSONFailed)
	}

	task, err := t.GetTask(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsJSONFailed)
	}

	if task.Succeeded() {
		fmt.Printf("@G{%s.}\n", done.Good)
	} else {
		fmt.Printf("@R{%s.}\n", done.Bad)
		if bail {
			os.Exit(OopsTaskFailed)
		}
	}
}

func watch(t *bosh.Target, res *http.Response, done Done) {
	var task bosh.Task
	err := t.InterpretJSON(res, &task)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsJSONFailed)
	}

	follow(t, task.ID, done, true)
}

func upload(path string) (io.Reader, int64, error) {
	var out progress.Reader

	file, err := os.Open(path)
	if err != nil {
		return &out, -1, err
	}

	info, err := file.Stat()
	if err != nil {
		file.Close()
		return &out, -1, err
	}

	out.Reader = file
	out.Size = int64(info.Size())
	out.Draw = progress.Console(os.Stdout, 50, 150, "uploading: ", " @G{done!}\n", '█')
	return &out, out.Size, nil
}
