package organization

import (
	"context"

	"github.com/vmaurya-21/Calance-Workflow/internal/infrastructure/github"
)

// Service handles organization business logic
type Service struct {
	githubOrg *github.OrganizationClient
}

// NewService creates a new organization service
func NewService() *Service {
	return &Service{
		githubOrg: github.NewOrganizationClient(),
	}
}

// GetUserOrganizations retrieves user's organizations
func (s *Service) GetUserOrganizations(ctx context.Context, token string) ([]Organization, error) {
	orgs, err := s.githubOrg.GetUserOrganizations(ctx, token)
	if err != nil {
		return nil, err
	}

	result := make([]Organization, len(orgs))
	for i, o := range orgs {
		result[i] = Organization{
			ID:          o.ID,
			Login:       o.Login,
			AvatarURL:   o.AvatarURL,
			Description: o.Description,
		}
	}
	return result, nil
}

// GetOrganizationRepositories retrieves organization repositories
func (s *Service) GetOrganizationRepositories(ctx context.Context, token, orgName string) ([]Repository, error) {
	repos, err := s.githubOrg.GetOrganizationRepositories(ctx, token, orgName)
	if err != nil {
		return nil, err
	}

	result := make([]Repository, len(repos))
	for i, r := range repos {
		result[i] = Repository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.FullName,
			Description:   r.Description,
			Private:       r.Private,
			HTMLURL:       r.HTMLURL,
			DefaultBranch: r.DefaultBranch,
		}
	}
	return result, nil
}

// GetUserRepositories retrieves all user repositories by organization
func (s *Service) GetUserRepositories(ctx context.Context, token string) (map[string][]Repository, error) {
	reposByOrg, err := s.githubOrg.GetUserRepositories(ctx, token)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]Repository)
	for orgName, repos := range reposByOrg {
		orgRepos := make([]Repository, len(repos))
		for i, r := range repos {
			orgRepos[i] = Repository{
				ID:            r.ID,
				Name:          r.Name,
				FullName:      r.FullName,
				Description:   r.Description,
				Private:       r.Private,
				HTMLURL:       r.HTMLURL,
				DefaultBranch: r.DefaultBranch,
			}
		}
		result[orgName] = orgRepos
	}
	return result, nil
}
