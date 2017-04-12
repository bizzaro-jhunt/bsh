package bosh

type DeployParams struct {
	Redact   bool
	Recreate bool
}

func (p DeployParams) DiffQuery() string {
	if p.Redact {
		return "?redact=true"
	} else {
		return "?redact=false"
	}
}
