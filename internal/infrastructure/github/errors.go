package github

import "errors"

var (
	// API errors
	ErrAPIFailed    = errors.New("github api request failed")
	ErrUnauthorized = errors.New("unauthorized: invalid or expired token")
	ErrForbidden    = errors.New("forbidden: insufficient permissions")
	ErrNotFound     = errors.New("not found: resource does not exist or no access")

	// Repository errors
	ErrRepositoryNotFound = errors.New("repository not found")
	ErrBranchNotFound     = errors.New("branch not found")

	// Workflow errors
	ErrWorkflowNotFound      = errors.New("workflow not found")
	ErrWorkflowAlreadyExists = errors.New("workflow already exists")

	// Organization errors
	ErrOrganizationNotFound = errors.New("organization not found")
)
