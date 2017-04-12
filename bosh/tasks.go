package bosh

import (
	"strings"

	"github.com/jhunt/bsh/query"
)

type TasksFilter struct {
	States     []string
	Deployment string
	ContextID  string
	Limit      int
	Verbose    int
}

func (tf TasksFilter) String() string {
	q := query.New()
	if len(tf.States) != 0 {
		q.Set("state", strings.Join(tf.States, ","))
	}
	q.Maybe("deployment", tf.Deployment)
	q.Maybe("context_id", tf.ContextID)
	if tf.Limit > 0 {
		q.Set("limit", tf.Limit)
	}
	if tf.Verbose != 0 {
		q.Set("verbose", tf.Verbose)
	}

	return q.String()
}

func (t Target) GetTasks(filter TasksFilter) ([]Task, error) {
	var l []Task
	return l, t.GetJSON("/tasks"+filter.String(), &l)
}
