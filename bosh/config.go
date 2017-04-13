package bosh

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Current string             `yaml:"current"`
	Targets map[string]*Target `yaml:"targets"`
}

func DefaultConfig(path string) error {
	var c Config
	b, err := ioutil.ReadFile(path)
	if err != nil && os.IsNotExist(err) {
		b, err = yaml.Marshal(&c)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, b, 0666)
		return err
	}
	return nil
}

func ReadConfig(path string) (Config, error) {
	var c Config

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (c Config) Target(name string) (*Target, error) {
	t, ok := c.Targets[name]
	if !ok {
		return nil, fmt.Errorf("BOSH target '%s' is not defined.")
	}

	return t, nil
}

func (c Config) CurrentTarget() (*Target, error) {
	if c.Current == "" {
		return nil, fmt.Errorf("no BOSH target is currently selected")
	}

	t, ok := c.Targets[c.Current]
	if !ok {
		return nil, fmt.Errorf("BOSH target '%s' is set as the current target, but is not defined.")
	}

	return t, nil
}

func (c Config) Save(path string) error {
	b, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0666)
	if err != nil {
		return err
	}

	return nil
}
