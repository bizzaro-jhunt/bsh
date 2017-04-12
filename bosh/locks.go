package bosh

type Lock struct {
	Type     string `json:"type"`
	Resource string `json:"resource"`
	Timeout  string `json:"timeout"`
}

func (t Target) GetLocks() ([]Lock, error) {
	var l []Lock
	return l, t.GetJSON("/locks", &l)
}
