package github

import (
	"context"
	"fmt"
	"net/http"
)

// OrganizationClient handles GitHub organization operations
type OrganizationClient struct {
	*Client
}

// NewOrganizationClient creates a new organization client
func NewOrganizationClient() *OrganizationClient {
	return &OrganizationClient{
		Client: NewClient(),
	}
}

// GetUserOrganizations retrieves all organizations for the authenticated user
func (oc *OrganizationClient) GetUserOrganizations(ctx context.Context, token string) ([]Organization, error) {
	path := "/user/orgs"
	resp, err := oc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var orgs []Organization
	if err := resp.UnmarshalJSON(&orgs); err != nil {
		return nil, err
	}

	return orgs, nil
}

// GetOrganizationRepositories retrieves all repositories for an organization
func (oc *OrganizationClient) GetOrganizationRepositories(ctx context.Context, token, orgName string) ([]Repository, error) {
	path := fmt.Sprintf("/orgs/%s/repos", orgName)
	resp, err := oc.doRequest(ctx, token, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	var repos []Repository
	if err := resp.UnmarshalJSON(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

// GetUserRepositories retrieves all repositories accessible to the user from their organizations
func (oc *OrganizationClient) GetUserRepositories(ctx context.Context, token string) (map[string][]Repository, error) {
	// First get all organizations
	orgs, err := oc.GetUserOrganizations(ctx, token)
	if err != nil {
		return nil, err
	}

	// Then get repositories for each organization
	repositoriesByOrg := make(map[string][]Repository)
	for _, org := range orgs {
		repos, err := oc.GetOrganizationRepositories(ctx, token, org.Login)
		if err != nil {
			// Log error but continue with other organizations
			continue
		}
		repositoriesByOrg[org.Login] = repos
	}

	return repositoriesByOrg, nil
}
