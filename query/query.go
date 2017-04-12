package query

import (
	"fmt"
	"net/url"
)

type Query map[string]string

func New() Query {
	return make(map[string]string)
}

func (q Query) Set(key string, value interface{}) {
	q[key] = url.QueryEscape(fmt.Sprintf("%v", value))
}

func (q Query) Bool(key string, value bool) {
	if value {
		q.Set(key, "true")
	} else {
		q.Set(key, "false")
	}
}

func (q Query) Maybe(key, value string) {
	if value != "" {
		q.Set(key, value)
	}
}

func (q Query) String() string {
	if len(q) == 0 {
		return ""
	}
	s := ""
	for k, v := range q {
		if s != "" {
			s += "&"
		}
		s += k + "=" + v
	}
	return "?" + s
}
