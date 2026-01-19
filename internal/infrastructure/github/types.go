package github

// Repository represents a GitHub repository
type Repository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	Private       bool   `json:"private"`
	HTMLURL       string `json:"html_url"`
	DefaultBranch string `json:"default_branch"`
	Owner         Owner  `json:"owner"`
}

// Owner represents a repository or organization owner
type Owner struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
}

// Branch represents a repository branch
type Branch struct {
	Name      string `json:"name"`
	Commit    Commit `json:"commit"`
	Protected bool   `json:"protected"`
}

// Commit represents a commit
type Commit struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

// Ref represents a Git reference
type Ref struct {
	Ref    string `json:"ref"`
	NodeID string `json:"node_id"`
	URL    string `json:"url"`
	Object struct {
		SHA  string `json:"sha"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"object"`
}

// PullRequest represents a pull request
type PullRequest struct {
	ID      int64  `json:"id"`
	Number  int    `json:"number"`
	State   string `json:"state"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	Head    struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"base"`
}

// Content represents a file or directory in a repository
type Content struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	GitURL      string `json:"git_url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content,omitempty"`
	Encoding    string `json:"encoding,omitempty"`
}

// User represents a GitHub user
type User struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Company   string `json:"company"`
}

// Organization represents a GitHub organization
type Organization struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	AvatarURL   string `json:"avatar_url"`
	Description string `json:"description"`
}
