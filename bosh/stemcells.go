package bosh

func (t Target) GetStemcells() ([]Stemcell, error) {
	var l []Stemcell
	return l, t.GetJSON("/stemcells", &l)
}
