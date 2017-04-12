package bosh

import (
	"fmt"
	"strings"
)

type TasksFilter struct {
	States     []string
	Deployment string
	ContextID  string
	Limit      int
	Verbose    int
}

func (tf TasksFilter) String() string {
	l := make([]string, 0)
	if len(tf.States) != 0 {
		l = append(l, fmt.Sprintf("state=%s", strings.Join(tf.States, ",")))
	}
	if tf.Deployment != "" {
		l = append(l, fmt.Sprintf("deployment=%s", tf.Deployment))
	}
	if tf.ContextID != "" {
		l = append(l, fmt.Sprintf("context_id=%s", tf.ContextID))
	}
	if tf.Limit > 0 {
		l = append(l, fmt.Sprintf("limit=%d", tf.Limit))
	}
	if tf.Verbose != 0 {
		l = append(l, fmt.Sprintf("verbose=%d", tf.Verbose))
	}

	if len(l) != 0 {
		return "?" + strings.Join(l, "&")
	}
	return ""
}

func (t Target) GetTasks(filter TasksFilter) ([]Task, error) {
	var l []Task
	return l, t.GetJSON("/tasks" + filter.String(), &l)
}
