package bosh

import (
	"fmt"
	"io"
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

func (t Task) Completed() bool {
	return !(t.State == "queued" || t.State == "processing" || t.State == "cancelling")
}

func (t Task) Succeeded() bool {
	return t.State == "done"
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

func (t Target) Follow(out io.Writer, id int) error {
	poller := make(chan error)
	tracer := make(chan error)
	rd, wr := io.Pipe()
	go func() {
		offset := 0
		for {
			/* strategy: keep poking BOSH until the task is no longer running */
			task, err := t.GetTask(id)
			if err != nil {
				poller <- err
				return
			}

			/* go get our output */
			output, err := t.getTaskOutput(task, "event")
			if err != nil {
				poller <- err
				return
			}

			if len(output) > offset {
				wr.Write([]byte(output[offset:]))
				offset = len(output)
			}

			if task.Completed() {
				poller <- nil
				return
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	go func() {
		tracer <- TraceEvents(out, rd)
	}()

	err := <-poller
	wr.Close()
	if err != nil {
		<-tracer
		return err
	}

	return <-tracer
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
