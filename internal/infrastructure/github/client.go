package github

import (
	"context"
	"fmt"
	"net/http"

	pkghttp "github.com/vmaurya-21/Calance-Workflow/internal/pkg/http"
)

// Client handles GitHub API interactions
type Client struct {
	httpClient *pkghttp.Client
	baseURL    string
}

// NewClient creates a new GitHub API client
func NewClient() *Client {
	return &Client{
		httpClient: pkghttp.NewClient(),
		baseURL:    "https://api.github.com",
	}
}

// doRequest performs a GitHub API request with standard headers
func (c *Client) doRequest(ctx context.Context, token, method, path string, body interface{}) (*pkghttp.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	headers := map[string]string{
		"Authorization":        "Bearer " + token,
		"Accept":               "application/vnd.github+json",
		"X-GitHub-Api-Version": "2022-11-28",
	}

	resp, err := c.httpClient.Do(ctx, pkghttp.Request{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	})

	if err != nil {
		return nil, fmt.Errorf("github api request failed: %w", err)
	}

	return resp, nil
}

// checkResponse checks if the response is successful
func checkResponse(resp *pkghttp.Response) error {
	if resp.IsSuccess() {
		return nil
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return fmt.Errorf("%w: %s", ErrAPIFailed, resp.GetErrorMessage())
	}
}
