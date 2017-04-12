package bosh

func (t Target) GetDeployments() ([]Deployment, error) {
	var l []Deployment
	return l, t.GetJSON("/deployments", &l)
}
