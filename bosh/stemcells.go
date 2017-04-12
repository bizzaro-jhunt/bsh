package bosh

import (
	"fmt"
)

func (t Target) GetStemcells() ([]Stemcell, error) {
	var l []Stemcell

	r, err := t.Get("/stemcells")
	if err != nil {
		return l, err
	}

	if r.StatusCode != 200 {
		return l, fmt.Errorf("BOSH API returned %s", r.Status)
	}

	err = t.InterpretJSON(r, &l)
	if err != nil {
		return l, err
	}

	return l, err
}
