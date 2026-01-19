package organization

// Organization represents a GitHub organization
type Organization struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	AvatarURL   string `json:"avatar_url"`
	Description string `json:"description"`
}

// Repository represents a repository
type Repository struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	Private       bool   `json:"private"`
	HTMLURL       string `json:"html_url"`
	DefaultBranch string `json:"default_branch"`
}
