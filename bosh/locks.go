package bosh

import (
	"fmt"
)

type Lock struct {
	Type     string `json:"type"`
	Resource string `json:"resource"`
	Timeout  string `json:"timeout"`
}

func (t Target) GetLocks() ([]Lock, error) {
	var l []Lock

	r, err := t.Get("/locks")
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
