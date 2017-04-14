package main

import (
	fmt "github.com/jhunt/go-ansi"
)

func eitherIgnoreItOrDont(ignore bool, deployment string, opt Opt, args []string) error {
	_, t := targeting(opt)

	l := make([][]string, 0)
	for len(args) > 0 {
		name, id, rest, err := thingversion(args)
		if err != nil {
			return err
		}

		l = append(l, []string{name, id})
		args = rest
	}

	if len(l) == 0 {
		return fmt.Errorf("usage...")
	}

	for _, p := range l {
		if ignore {
			fmt.Printf("@R{ignoring} %s instance @C{%s}/@C{%s}...\n", deployment, p[0], p[1])
		} else {
			fmt.Printf("@G{unignoring} %s instance @C{%s}/@C{%s}...\n", deployment, p[0], p[1])
		}
		r, err := t.Put(fmt.Sprintf("/deployments/%s/instance_groups/%s/%s/ignore", deployment, p[0], p[1]),
			struct {
				Ignore bool `json:"ignore"`
			}{ignore})
		if err != nil {
			return err
		}

		if r.StatusCode != 200 {
			return fmt.Errorf("BOSH API returned %s", r.Status)
		}
	}

	return nil
}

func runIgnoreInstance(opt Opt, command string, args []string) {
	eitherIgnoreItOrDont(true, opt.Deployment, opt, args)
}

func runUnignoreInstance(opt Opt, command string, args []string) {
	eitherIgnoreItOrDont(false, opt.Deployment, opt, args)
}
