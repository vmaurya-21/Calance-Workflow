package workflow

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/vmaurya-21/Calance-Workflow/internal/infrastructure/github"
	"github.com/vmaurya-21/Calance-Workflow/internal/infrastructure/template"
	"github.com/vmaurya-21/Calance-Workflow/internal/pkg/logger"
)

// Service handles workflow business logic
type Service struct {
	githubClient *github.WorkflowClient
	ec2Template  *template.EC2Generator
	k8sTemplate  *template.KubernetesGenerator
}

// NewService creates a new workflow service
func NewService() *Service {
	return &Service{
		githubClient: github.NewWorkflowClient(),
		ec2Template:  template.NewEC2Generator(),
		k8sTemplate:  template.NewKubernetesGenerator(),
	}
}

// GenerateWorkflow generates workflow YAML based on request
func (s *Service) GenerateWorkflow(req *Request) (string, error) {
	if !isValidWorkflowName(req.WorkflowName) {
		return "", ErrInvalidWorkflowName
	}

	if err := req.Validate(); err != nil {
		return "", err
	}

	var yamlContent string
	var err error

	switch req.DeploymentType {
	case DeploymentTypeEC2:
		yamlContent, err = s.ec2Template.Generate(req)
	case DeploymentTypeKubernetes:
		yamlContent, err = s.k8sTemplate.Generate(req)
	default:
		return "", ErrInvalidDeploymentType
	}

	if err != nil {
		logger.Error().Err(err).Str("deployment_type", string(req.DeploymentType)).Msg("Failed to generate workflow")
		return "", fmt.Errorf("%w: %v", ErrTemplateGenerationFailed, err)
	}

	return yamlContent, nil
}

// CreateWorkflow creates a workflow in GitHub repository
func (s *Service) CreateWorkflow(ctx context.Context, token, owner, repo, workflowName, content string) (*Response, error) {
	if err := s.githubClient.VerifyRepository(ctx, token, owner, repo); err != nil {
		return nil, err
	}

	defaultBranch, err := s.githubClient.GetDefaultBranch(ctx, token, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get default branch: %w", err)
	}

	branchName := fmt.Sprintf("workflow/%s-%d", workflowName, time.Now().Unix())

	baseSHA, err := s.githubClient.GetBranchSHA(ctx, token, owner, repo, defaultBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to get base branch SHA: %w", err)
	}

	if err := s.githubClient.CreateBranch(ctx, token, owner, repo, branchName, baseSHA); err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	filePath := fmt.Sprintf(".github/workflows/%s.yml", workflowName)
	message := fmt.Sprintf("Add workflow: %s", workflowName)

	if err := s.githubClient.CreateFile(ctx, token, owner, repo, filePath, content, message, branchName); err != nil {
		return nil, err
	}

	prTitle := fmt.Sprintf("Add workflow: %s", workflowName)
	prBody := fmt.Sprintf("This PR adds the GitHub Actions workflow for `%s`.\n\nGenerated automatically by Calance Workflow Manager.", workflowName)
	prURL, prNumber, err := s.githubClient.CreatePullRequest(ctx, token, owner, repo, branchName, defaultBranch, prTitle, prBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	return &Response{
		Owner:        owner,
		Repository:   repo,
		WorkflowName: workflowName,
		FilePath:     filePath,
		FileURL:      prURL,
		Message:      fmt.Sprintf("Pull request #%d created for workflow '%s'", prNumber, workflowName),
	}, nil
}

// GetWorkflows retrieves all workflows from a repository
func (s *Service) GetWorkflows(ctx context.Context, token, owner, repo string) ([]File, error) {
	contents, err := s.githubClient.GetWorkflowFiles(ctx, token, owner, repo)
	if err != nil {
		return nil, err
	}

	files := make([]File, len(contents))
	for i, c := range contents {
		files[i] = File{
			Name:        c.Name,
			Path:        c.Path,
			SHA:         c.SHA,
			Size:        c.Size,
			URL:         c.HTMLURL,
			DownloadURL: c.DownloadURL,
		}
	}

	return files, nil
}

func isValidWorkflowName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched && len(name) > 0 && len(name) <= 255
}

// GetWorkflowContent retrieves the content of a workflow file
func (s *Service) GetWorkflowContent(ctx context.Context, token, owner, repo, filePath string) (*FileContentResponse, error) {
	// Validate that the file path is a workflow file
	if !strings.HasPrefix(filePath, ".github/workflows/") {
		return nil, fmt.Errorf("invalid workflow file path: must be in .github/workflows/")
	}

	if !strings.HasSuffix(filePath, ".yml") && !strings.HasSuffix(filePath, ".yaml") {
		return nil, fmt.Errorf("invalid workflow file: must be a .yml or .yaml file")
	}

	// Fetch file content from GitHub
	content, sha, err := s.githubClient.GetFileContent(ctx, token, owner, repo, filePath)
	if err != nil {
		logger.Error().Err(err).
			Str("owner", owner).
			Str("repo", repo).
			Str("file_path", filePath).
			Msg("Failed to fetch workflow file content")
		return nil, fmt.Errorf("failed to fetch file content: %w", err)
	}

	// Extract file name from path
	parts := strings.Split(filePath, "/")
	fileName := parts[len(parts)-1]

	return &FileContentResponse{
		Name:    fileName,
		Path:    filePath,
		SHA:     sha,
		Size:    len(content),
		Content: content,
	}, nil
}

// UpdateWorkflow updates an existing workflow file and creates a PR
func (s *Service) UpdateWorkflow(ctx context.Context, token string, req *UpdateWorkflowRequest) (*Response, error) {
	// Validate that the file path is a workflow file
	if !strings.HasPrefix(req.FilePath, ".github/workflows/") {
		return nil, fmt.Errorf("invalid workflow file path: must be in .github/workflows/")
	}

	if !strings.HasSuffix(req.FilePath, ".yml") && !strings.HasSuffix(req.FilePath, ".yaml") {
		return nil, fmt.Errorf("invalid workflow file: must be a .yml or .yaml file")
	}

	// Verify repository exists
	if err := s.githubClient.VerifyRepository(ctx, token, req.Owner, req.Repository); err != nil {
		return nil, err
	}

	// Get default branch
	defaultBranch, err := s.githubClient.GetDefaultBranch(ctx, token, req.Owner, req.Repository)
	if err != nil {
		return nil, fmt.Errorf("failed to get default branch: %w", err)
	}

	// Extract workflow name from file path
	parts := strings.Split(req.FilePath, "/")
	fileName := parts[len(parts)-1]
	workflowName := strings.TrimSuffix(fileName, ".yml")
	workflowName = strings.TrimSuffix(workflowName, ".yaml")

	// Create a new branch for the update
	branchName := fmt.Sprintf("update-workflow/%s-%d", workflowName, time.Now().Unix())

	baseSHA, err := s.githubClient.GetBranchSHA(ctx, token, req.Owner, req.Repository, defaultBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to get base branch SHA: %w", err)
	}

	if err := s.githubClient.CreateBranch(ctx, token, req.Owner, req.Repository, branchName, baseSHA); err != nil {
		return nil, fmt.Errorf("failed to create branch: %w", err)
	}

	// Set default commit message if not provided
	message := req.CommitMessage
	if message == "" {
		message = fmt.Sprintf("Update workflow: %s", workflowName)
	}

	// Update the file on the new branch
	if err := s.githubClient.UpdateFile(ctx, token, req.Owner, req.Repository, req.FilePath, req.Content, message, branchName, req.SHA); err != nil {
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	// Create pull request
	prTitle := fmt.Sprintf("Update workflow: %s", workflowName)
	prBody := fmt.Sprintf("This PR updates the GitHub Actions workflow `%s`.\n\nUpdated automatically by Calance Workflow Manager.", workflowName)
	prURL, prNumber, err := s.githubClient.CreatePullRequest(ctx, token, req.Owner, req.Repository, branchName, defaultBranch, prTitle, prBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	return &Response{
		Owner:        req.Owner,
		Repository:   req.Repository,
		WorkflowName: workflowName,
		FilePath:     req.FilePath,
		FileURL:      prURL,
		Message:      fmt.Sprintf("Pull request #%d created for workflow '%s' update", prNumber, workflowName),
	}, nil
}
