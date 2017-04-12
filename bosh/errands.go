package bosh

import (
	"fmt"
)

type Errand struct {
	Name string `json:"name"`
}

func (t Target) GetErrandsFor(deployment string) ([]Errand, error) {
	var l []Errand
	return l, t.GetJSON(fmt.Sprintf("/deployments/%s/errands", deployment), &l)
}
