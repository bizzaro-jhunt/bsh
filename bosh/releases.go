package bosh

import (
	"fmt"
)

func (t Target) GetReleases() ([]Release, error) {
	var l []Release

	r, err := t.Get("/releases")
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
