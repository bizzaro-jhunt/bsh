package bosh

import (
	"fmt"
)

type CleanupParams struct {
	Config struct {
		RemoveAll bool `json:"remove_all"`
	} `json:"config"`
}

func (t Target) Cleanup(all bool) (Task, error) {
	p := CleanupParams{}
	p.Config.RemoveAll = all

	r, err := t.Post("/cleanup", p)
	if err != nil {
		return Task{}, err
	}

	if r.StatusCode != 200 {
		return Task{}, fmt.Errorf("BOSH API returned %s", r.Status)
	}

	var task Task
	err = t.InterpretJSON(r, &task)
	if err != nil {
		return Task{}, err
	}

	return task, err
}
