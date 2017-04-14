package bosh

import (
	"fmt"

	"github.com/jhunt/bsh/query"
)

func (t Target) GetReleases() ([]Release, error) {
	var l []Release
	return l, t.GetJSON("/releases", &l)
}

func (t Target) GetRelease(name, version string) (Release, error) {
	var r Release
	q := query.New()
	q.Maybe("version", version)
	return r, t.GetJSON(fmt.Sprintf("/releases/%s%s", name, q), &r)
}
