package bosh

import (
	"net/http"
)

type Target struct {
	Alias    string `yaml:"alias"`
	UUID     string `yaml:"uuid"`
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	CaCert   string `yaml:"ca_cert"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Insecure bool   `yaml:"insecure"`

	ua http.Client
}

func (t *Target) Sync() error {
	info, err := t.GetInfo()
	if err != nil {
		return err
	}

	t.UUID = info.UUID
	t.Name = info.Name
	return nil
}
