package bosh

func (t Target) GetReleases() ([]Release, error) {
	var l []Release
	return l, t.GetJSON("/releases", &l)
}
