package bosh

type ReleaseJob struct {
	Name        string `json:"name"`
	BlobstoreID string `json:"blobstore_id"`
	SHA1        string `json:"sha1"`
	Fingerprint string `json:"fingerprint"`

	/* consumes / provides */
}

type ReleaseCompiledPackage struct {
	Stemcell    string `json:"stemcell"`
	BlobstoreID string `json:"blobstore_id"`
	SHA1        string `json:"sha1"`
}

type ReleasePackage struct {
	Name        string `json:"name"`
	BlobstoreID string `json:"blobstore_id"`
	SHA1        string `json:"sha1"`
	Fingerprint string `json:"fingerprint"`

	CompiledPackages []ReleaseCompiledPackage `json:"compiled_packages"`
}

type ReleaseVersion struct {
	Version  string   `json:"version"`
	Commit   string   `json:"commit_hash"`
	Dirty    bool     `json:"uncommitted_changes"`
	Deployed bool     `json:"currently_deployed"`
	Jobs     []string `json:"job_names"`
}

type Release struct {
	Name     string           `json:"name"`
	Version  string           `json:"version"`
	Versions []ReleaseVersion `json:"release_versions"`

	Jobs     []ReleaseJob     `json:"jobs"`
	Packages []ReleasePackage `json:"packages"`
}

type Stemcell struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	OS      string `json:"operating_system"`
	CID     string `json:"cid"`

	Deployments []struct {
		Name string `json:"name"`
	} `json:"deployments"`
}

type Deployment struct {
	Name        string     `json:"name"`
	Releases    []Release  `json:"releases"`
	Stemcells   []Stemcell `json:"stemcells"`
	CloudConfig string     `json:"cloud_config"`

	Manifest string `json:"manifest"`
}

type CloudConfig struct {
	Properties string `json:"properties"`
	CreatedAt  string `json:"created_at"`
}

type RuntimeConfig struct {
	Properties string `json:"properties"`
	CreatedAt  string `json:"created_at"`
}
