package workflow

import "errors"

var (
	// Validation errors
	ErrEC2CommonFieldsRequired        = errors.New("ec2CommonFields is required for EC2 deployment type")
	ErrEC2ProjectsRequired            = errors.New("ec2Projects is required for EC2 deployment type")
	ErrKubernetesCommonFieldsRequired = errors.New("kubernetesCommonFields is required for Kubernetes deployment type")
	ErrKubernetesProjectsRequired     = errors.New("kubernetesProjects is required for Kubernetes deployment type")
	ErrInvalidWorkflowName            = errors.New("workflow name must contain only alphanumeric characters, hyphens, and underscores")
	ErrInvalidDeploymentType          = errors.New("deployment type must be either 'ec2' or 'kubernetes'")

	// GitHub API errors
	ErrWorkflowAlreadyExists   = errors.New("workflow file already exists in repository")
	ErrGitHubAPIFailed         = errors.New("failed to communicate with GitHub API")
	ErrInsufficientPermissions = errors.New("insufficient permissions to create workflow file")

	// Template errors
	ErrTemplateGenerationFailed = errors.New("failed to generate workflow template")
	ErrInvalidYAMLGenerated     = errors.New("generated YAML is invalid")
)
