package main

import (
	"encoding/json"
	fmt "github.com/jhunt/go-ansi"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

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

func targeting(opt Opt) (bosh.Config, *bosh.Target) {
	cfg := readConfigFrom(opt.Config)

	var t *bosh.Target
	var err error

	if opt.BOSHTarget == "" {
		t, err = cfg.CurrentTarget()
	} else {
		t, err = cfg.Target(opt.BOSHTarget)
	}
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

func deinterface(v interface{}) interface{} {
	switch v.(type) {
	case map[interface{}]interface{}:
		m := v.(map[interface{}]interface{})
		o := make(map[string]interface{})
		for k := range m {
			o[fmt.Sprintf("%v", k)] = deinterface(m[k])
		}
		return o

	case []interface{}:
		l := v.([]interface{})
		o := make([]interface{}, len(l))
		for i := range l {
			o[i] = deinterface(l[i])
		}
		return o

	default:
		return v
	}
}

func yamlr(s string) interface{} {
	var src map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(s), &src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
		os.Exit(OopsJSONFailed)
	}

	return deinterface(src)
}

func downloadto(out io.Writer, contents, path string, force bool) {
	if path != "" {
		flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		if !force {
			flags |= os.O_EXCL
		}
		f, err := os.OpenFile(path, flags, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@R{!!! %s}\n", err)
			os.Exit(OopsDownloadFailed)
		}
		out = f
	}

	fmt.Fprintf(out, contents)
}

func kbsize(kb string) string {
	k, err := strconv.ParseFloat(kb, 32)
	if err != nil {
		return kb
	}

	if k == 0 {
		return "   0 "
	}

	unit := "K"
	if k > 1024 {
		k /= 1024.0
		unit = "M"

		if k > 1024 {
			k /= 1024.0
			unit = "G"

			if k > 1024 {
				k /= 1024.0
				unit = "T"
			}
		}
	}

	return fmt.Sprintf("% 4d%s", int(k), unit)
}

func percent(s string) string {
	return fmt.Sprintf("%s%%", s)
}

func clocked(t uint64) string {
	var d, h, m, s uint64

	s = t % 60
	t /= 60

	m = t % 60
	t /= 60

	h = t % 24
	t /= 24

	d = t

	if d > 0 {
		return fmt.Sprintf("%dd %d:%02dm", d, h, m)
	}
	if h > 0 || m > 0 {
		return fmt.Sprintf("%d:%02dm", h, m)
	}
	return fmt.Sprintf("%d:%02ds", m, s)
}
