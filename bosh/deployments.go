package bosh

import (
	"fmt"
)

func (t Target) GetDeployments() ([]Deployment, error) {
	var l []Deployment
	return l, t.GetJSON("/deployments", &l)
}

func (t Target) GetDeployment(name string) (Deployment, error) {
	var d Deployment
	return d, t.GetJSON(fmt.Sprintf("/deployments/%s", name), &d)
}
