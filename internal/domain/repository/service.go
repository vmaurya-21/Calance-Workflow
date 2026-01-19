package repository

import (
	"context"

	"github.com/vmaurya-21/Calance-Workflow/internal/infrastructure/github"
)

// Service handles repository business logic
type Service struct {
	githubRepo *github.RepositoryClient
}

// NewService creates a new repository service
func NewService() *Service {
	return &Service{
		githubRepo: github.NewRepositoryClient(),
	}
}

// GetBranches retrieves branches for a repository
func (s *Service) GetBranches(ctx context.Context, token, owner, repo string) ([]Branch, error) {
	branches, err := s.githubRepo.GetBranches(ctx, token, owner, repo)
	if err != nil {
		return nil, err
	}

	result := make([]Branch, len(branches))
	for i, b := range branches {
		result[i] = Branch{
			Name:      b.Name,
			CommitSHA: b.Commit.SHA,
			Protected: b.Protected,
		}
	}
	return result, nil
}

// GetCommits retrieves commits for a branch
func (s *Service) GetCommits(ctx context.Context, token, owner, repo, branch string, perPage int) ([]interface{}, error) {
	return s.githubRepo.GetCommits(ctx, token, owner, repo, branch, perPage)
}

// GetTags retrieves tags for a repository
func (s *Service) GetTags(ctx context.Context, token, owner, repo string) ([]interface{}, error) {
	return s.githubRepo.GetTags(ctx, token, owner, repo)
}

// CreateTag creates a new tag
func (s *Service) CreateTag(ctx context.Context, token, owner, repo, tagName, commitSHA string) (*TagReference, error) {
	ref, err := s.githubRepo.CreateTag(ctx, token, owner, repo, tagName, commitSHA)
	if err != nil {
		return nil, err
	}

	return &TagReference{
		Ref:       ref.Ref,
		ObjectSHA: ref.Object.SHA,
		URL:       ref.URL,
	}, nil
}

// GetWorkflowRuns retrieves workflow runs
func (s *Service) GetWorkflowRuns(ctx context.Context, token, owner, repo string, perPage int) ([]interface{}, error) {
	return s.githubRepo.GetWorkflowRuns(ctx, token, owner, repo, perPage)
}

// GetWorkflowRunDetail retrieves workflow run details
func (s *Service) GetWorkflowRunDetail(ctx context.Context, token, owner, repo string, runID int64) (interface{}, []interface{}, error) {
	return s.githubRepo.GetWorkflowRunDetail(ctx, token, owner, repo, runID)
}

// GetJobLogs retrieves job logs
func (s *Service) GetJobLogs(ctx context.Context, token, owner, repo string, jobID int64) (string, error) {
	return s.githubRepo.GetJobLogs(ctx, token, owner, repo, jobID)
}
