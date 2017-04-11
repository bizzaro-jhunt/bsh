package bosh

import (
	"fmt"
	"io/ioutil"
	"time"
)

type Task struct {
	ID          int     `json:"id"`
	State       string  `json:"state"`
	Description string  `json:"description"`
	StartedAt   int     `json:"started_at"` // undocumented
	Timestamp   int     `json:"timestamp"`
	Result      *string `json:"result"`
	User        string  `json:"user"`
	ContextID   *string `json:"context_id"`
	Deployment  string  `json:"deployment"` // undocumented
}

func (t Target) GetTask(id int) (Task, error) {
	var task Task
	r, err := t.Get(fmt.Sprintf("/tasks/%d", id))
	if err != nil {
		return task, err
	}

	if r.StatusCode == 204 {
		return task, nil
	}
	if r.StatusCode != 200 {
		return task, fmt.Errorf("BOSH API returned %s", r.Status)
	}

	if err = t.InterpretJSON(r, &task); err != nil {
		return task, err
	}

	return task, nil
}

func (t Target) WaitTask(id int, sleep time.Duration) (Task, error) {
	for {
		task, err := t.GetTask(id)
		if err != nil {
			return task, err
		}
		if task.Result != nil {
			return task, nil
		}

		time.Sleep(sleep)
	}
}

func (t Target) getTaskOutput(task Task, what string) (string, error) {
	r, err := t.Get(fmt.Sprintf("/tasks/%d/output?type=%s", task.ID, what))
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (t Target) GetTaskDebugOutput(task Task) (string, error) {
	return t.getTaskOutput(task, "debug")
}

func (t Target) GetTaskResult(task Task) (string, error) {
	return t.getTaskOutput(task, "result")
}
