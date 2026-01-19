package repository

// Branch represents a repository branch
type Branch struct {
	Name      string `json:"name"`
	CommitSHA string `json:"commit_sha"`
	Protected bool   `json:"protected"`
}

// TagReference represents a tag reference
type TagReference struct {
	Ref       string `json:"ref"`
	ObjectSHA string `json:"object_sha"`
	URL       string `json:"url"`
}
