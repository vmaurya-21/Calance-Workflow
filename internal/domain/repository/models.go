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

// Package represents a GitHub package
type Package struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	PackageType    string `json:"package_type"`
	Visibility     string `json:"visibility"`
	URL            string `json:"url"`
	HTMLURL        string `json:"html_url"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	OwnerLogin     string `json:"owner_login"`
	RepositoryName string `json:"repository_name,omitempty"`
}
