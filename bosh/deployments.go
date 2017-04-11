package bosh

import (
	"fmt"
)

type Release struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Stemcell struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Deployment struct {
	Name        string     `json:"name"`
	Releases    []Release  `json:"releases"`
	Stemcells   []Stemcell `json:"stemcells"`
	CloudConfig string     `json:"cloud_config"`
}

func (t Target) GetDeployments() ([]Deployment, error) {
	var l []Deployment

	r, err := t.Get("/deployments")
	if err != nil {
		return l, err
	}

	if r.StatusCode != 200 {
		return l, fmt.Errorf("BOSH API returned %s", r.Status)
	}

	if err = t.InterpretJSON(r, &l); err != nil {
		return l, err
	}
	return l, nil
}
