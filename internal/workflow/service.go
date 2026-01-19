package workflow

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"text/template"

	"github.com/vmaurya-21/Calance-Workflow/internal/logger"
)

// WorkflowService handles workflow generation and GitHub API interactions
type WorkflowService struct {
	httpClient *http.Client
}

// NewWorkflowService creates a new workflow service
func NewWorkflowService() *WorkflowService {
	return &WorkflowService{
		httpClient: &http.Client{},
	}
}

// GenerateWorkflowYAML generates the complete workflow YAML based on deployment type
func (s *WorkflowService) GenerateWorkflowYAML(request *WorkflowRequest) (string, error) {
	// Validate workflow name
	if !isValidWorkflowName(request.WorkflowName) {
		return "", ErrInvalidWorkflowName
	}

	// Validate request based on deployment type
	if err := request.Validate(); err != nil {
		return "", err
	}

	// Generate YAML based on deployment type
	var yamlContent string
	var err error

	switch request.DeploymentType {
	case DeploymentTypeEC2:
		yamlContent, err = s.generateEC2Workflow(request)
	case DeploymentTypeKubernetes:
		yamlContent, err = s.generateKubernetesWorkflow(request)
	default:
		return "", ErrInvalidDeploymentType
	}

	if err != nil {
		logger.Error().
			Err(err).
			Str("deployment_type", string(request.DeploymentType)).
			Msg("Failed to generate workflow YAML")
		return "", fmt.Errorf("%w: %v", ErrTemplateGenerationFailed, err)
	}

	return yamlContent, nil
}

// generateEC2Workflow generates EC2 deployment workflow YAML
func (s *WorkflowService) generateEC2Workflow(request *WorkflowRequest) (string, error) {
	tmpl := `name: Build & Publish Image (EC2)

on:
  push:
    tags:
      - v[0-9]+.[0-9]+.[0-9]+-rc[0-9]+
      - v[0-9]+.[0-9]+.[0-9]+

jobs:
  build-and-push-dockerimages:
    strategy:
      fail-fast: false
      matrix:
        project: [{{range $i, $p := .Projects}}{{if $i}}, {{end}}{{$p.Name}}{{end}}]
    permissions:
      contents: read
      packages: write
    secrets:
      IMAGE_REGISTRY_PASSWORD: {{"{{"}} secrets.IMAGE_REGISTRY_PASSWORD {{"}}"}}

    uses: Calance-US/calance-workflows/.github/workflows/build.yml@{{.EC2CommonFields.ReleaseTag}}
    with:
      image_name: {{.Owner}}/{{.Repository}}-{{"{{"}} matrix.project {{"}}"}}
      image_registry: {{"{{"}} vars.IMAGE_REGISTRY {{"}}"}}
      image_registry_username: {{"{{"}} vars.IMAGE_REGISTRY_USERNAME {{"}}"}}
      docker_context_path: {{"{{"}} matrix.project {{"}}"}}
      dockerfile_path: ./{{"{{"}} matrix.project {{"}}"}}/Dockerfile
{{range .Projects}}      dot_env_file_testing: |
{{indent 8 .DotEnvTesting}}
{{end}}
  deploy-to-ec2:
    needs: build-and-push-dockerimages
    strategy:
      fail-fast: false
      matrix:
        project: [{{range $i, $p := .EC2Projects}}{{if $i}}, {{end}}{{$p.Name}}{{end}}]
    permissions:
      contents: read
      packages: write

    uses: Calance-US/calance-workflows/.github/workflows/deploy-ec2.yml@{{.EC2CommonFields.ReleaseTag}}
    with:
      repository_name: {{"{{"}} github.event.repository.name {{"}}"}}
      image_name: {{.Owner}}/{{.Repository}}-{{"{{"}} matrix.project {{"}}"}}
      image_registry: {{"{{"}} vars.IMAGE_REGISTRY {{"}}"}}
      version: {{"{{"}} needs.build-and-push-dockerimages.outputs.version {{"}}"}}
      cluster_environment: {{"{{"}} needs.build-and-push-dockerimages.outputs.cluster_environment {{"}}"}}
      commit_id: {{"{{"}} needs.build-and-push-dockerimages.outputs.commit_id {{"}}"}}
      aws_region: {{.EC2CommonFields.AWSRegion}}
      jenkins_jobs: {{.EC2CommonFields.JenkinsJobs}}
      workflows_release: {{.EC2CommonFields.ReleaseTag}}
      codeowners_email_ids: {{.EC2CommonFields.CodeownersEmails}}
      devops_stakeholders_email_ids: {{.EC2CommonFields.DevopsStakeholdersEmails}}
{{range .EC2Projects}}      # EC2 specific configuration for {{.Name}}
      command: {{.Command}}
      port: {{.Port}}
{{if .DockerNetwork}}      docker_network: {{.DockerNetwork}}
{{end}}{{if .MountPath}}      mount_path: {{.MountPath}}
{{end}}{{if .EnableGPU}}      enable_gpu: true
{{end}}{{if .LogDriver}}      log_driver: {{.LogDriver}}
{{end}}{{if .LogDriverOptions}}      log_driver_options: {{.LogDriverOptions}}
{{end}}{{end}}    secrets:
      JENKINS_URL: {{"{{"}} secrets.JENKINS_URL {{"}}"}}
      JENKINS_USER: {{"{{"}} secrets.JENKINS_USER {{"}}"}}
      JENKINS_TOKEN: {{"{{"}} secrets.JENKINS_TOKEN {{"}}"}}
      SMTP_PASSWORD: {{"{{"}} secrets.SMTP_PASSWORD {{"}}"}}
      AWS_CREDENTIALS: {{"{{"}} secrets.AWS_CREDENTIALS {{"}}"}}
`

	funcMap := template.FuncMap{
		"indent": func(spaces int, text string) string {
			indent := strings.Repeat(" ", spaces)
			lines := strings.Split(text, "\n")
			for i, line := range lines {
				if line != "" {
					lines[i] = indent + line
				}
			}
			return strings.Join(lines, "\n")
		},
	}

	t, err := template.New("ec2-workflow").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, request); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// generateKubernetesWorkflow generates Kubernetes deployment workflow YAML
func (s *WorkflowService) generateKubernetesWorkflow(request *WorkflowRequest) (string, error) {
	tmpl := `name: Build & Publish Image (Kubernetes)

on:
  push:
    tags:
      - v[0-9]+.[0-9]+.[0-9]+-rc[0-9]+
      - v[0-9]+.[0-9]+.[0-9]+

jobs:
  build-and-push-dockerimages:
    strategy:
      fail-fast: false
      matrix:
        project: [{{range $i, $p := .Projects}}{{if $i}}, {{end}}{{$p.Name}}{{end}}]
    permissions:
      contents: read
      packages: write
    secrets:
      IMAGE_REGISTRY_PASSWORD: {{"{{"}} secrets.IMAGE_REGISTRY_PASSWORD {{"}}"}}

    uses: Calance-US/calance-workflows/.github/workflows/build.yml@{{.KubernetesCommonFields.ReleaseTag}}
    with:
      image_name: {{.Owner}}/{{.Repository}}-{{"{{"}} matrix.project {{"}}"}}
      image_registry: {{"{{"}} vars.IMAGE_REGISTRY {{"}}"}}
      image_registry_username: {{"{{"}} vars.IMAGE_REGISTRY_USERNAME {{"}}"}}
      docker_context_path: {{"{{"}} matrix.project {{"}}"}}
      dockerfile_path: ./{{"{{"}} matrix.project {{"}}"}}/Dockerfile
{{range .Projects}}      dot_env_file_testing: |
{{indent 8 .DotEnvTesting}}
{{end}}
  deploy-to-kubernetes:
    needs: build-and-push-dockerimages
    strategy:
      fail-fast: false
      matrix:
        project: [{{range $i, $p := .KubernetesProjects}}{{if $i}}, {{end}}{{$p.Name}}{{end}}]
    permissions:
      contents: read
      packages: write

    uses: Calance-US/calance-workflows/.github/workflows/deploy.yml@{{.KubernetesCommonFields.ReleaseTag}}
    with:
      repository_name: {{"{{"}} github.event.repository.name {{"}}"}}
      image_name: {{.Owner}}/{{.Repository}}-{{"{{"}} matrix.project {{"}}"}}
      release_name: {{.Repository}}-{{"{{"}} matrix.project {{"}}"}}
      image_registry: {{"{{"}} vars.IMAGE_REGISTRY {{"}}"}}
      version: {{"{{"}} needs.build-and-push-dockerimages.outputs.version {{"}}"}}
      cluster_environment: {{"{{"}} needs.build-and-push-dockerimages.outputs.cluster_environment {{"}}"}}
      commit_id: {{"{{"}} needs.build-and-push-dockerimages.outputs.commit_id {{"}}"}}
      jenkins_job_name: {{.KubernetesCommonFields.JenkinsJobName}}
      workflows_release: {{.KubernetesCommonFields.ReleaseTag}}
      helm_values_repository: {{.KubernetesCommonFields.HelmValuesRepository}}
      codeowners_email_ids: {{.KubernetesCommonFields.CodeownersEmailIds}}
      devops_stakeholders_email_ids: {{.KubernetesCommonFields.DevopsStakeholdersEmailIds}}
    secrets:
      JENKINS_URL: {{"{{"}} secrets.JENKINS_URL {{"}}"}}
      JENKINS_USER: {{"{{"}} secrets.JENKINS_USER {{"}}"}}
      JENKINS_TOKEN: {{"{{"}} secrets.JENKINS_TOKEN {{"}}"}}
      SMTP_PASSWORD: {{"{{"}} secrets.SMTP_PASSWORD {{"}}"}}
`

	funcMap := template.FuncMap{
		"indent": func(spaces int, text string) string {
			indent := strings.Repeat(" ", spaces)
			lines := strings.Split(text, "\n")
			for i, line := range lines {
				if line != "" {
					lines[i] = indent + line
				}
			}
			return strings.Join(lines, "\n")
		},
	}

	t, err := template.New("k8s-workflow").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, request); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// CreateWorkflowFile creates a workflow file in the GitHub repository
func (s *WorkflowService) CreateWorkflowFile(ctx context.Context, token, owner, repo, workflowName, content string) (*WorkflowResponse, error) {
	// Construct the file path
	filePath := fmt.Sprintf(".github/workflows/%s.yml", workflowName)

	// Check if file already exists
	exists, existingSHA, err := s.checkFileExists(ctx, token, owner, repo, filePath)
	if err != nil {
		logger.Error().
			Err(err).
			Str("path", filePath).
			Msg("Failed to check if workflow exists")
		return nil, fmt.Errorf("%w: %v", ErrGitHubAPIFailed, err)
	}

	// Prepare the request body
	requestBody := map[string]interface{}{
		"message": fmt.Sprintf("Create workflow: %s", workflowName),
		"content": encodeBase64(content),
	}

	// If file exists, include the SHA for update
	if exists {
		requestBody["message"] = fmt.Sprintf("Update workflow: %s", workflowName)
		requestBody["sha"] = existingSHA
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Create the GitHub API request
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, filePath)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Execute the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		logger.Error().
			Err(err).
			Str("url", url).
			Msg("Failed to create workflow file")
		return nil, fmt.Errorf("%w: %v", ErrGitHubAPIFailed, err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, _ := io.ReadAll(resp.Body)

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger.Error().
			Int("status", resp.StatusCode).
			Str("body", string(respBody)).
			Msg("GitHub API returned error")

		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
			return nil, ErrInsufficientPermissions
		}

		return nil, fmt.Errorf("%w: status %d, body: %s", ErrGitHubAPIFailed, resp.StatusCode, string(respBody))
	}

	// Parse response
	var apiResponse struct {
		Content struct {
			Name        string `json:"name"`
			Path        string `json:"path"`
			SHA         string `json:"sha"`
			HTMLURL     string `json:"html_url"`
			DownloadURL string `json:"download_url"`
		} `json:"content"`
		Commit struct {
			SHA string `json:"sha"`
		} `json:"commit"`
	}

	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		logger.Error().
			Err(err).
			Str("body", string(respBody)).
			Msg("Failed to parse GitHub API response")
		return nil, err
	}

	response := &WorkflowResponse{
		Owner:        owner,
		Repository:   repo,
		WorkflowName: workflowName,
		FilePath:     apiResponse.Content.Path,
		FileURL:      apiResponse.Content.HTMLURL,
		ContentSHA:   apiResponse.Content.SHA,
		Message:      fmt.Sprintf("Workflow '%s' %s successfully", workflowName, map[bool]string{true: "updated", false: "created"}[exists]),
	}

	logger.Info().
		Str("owner", owner).
		Str("repo", repo).
		Str("workflow", workflowName).
		Str("path", filePath).
		Str("action", map[bool]string{true: "updated", false: "created"}[exists]).
		Msg("Workflow file created/updated successfully")

	return response, nil
}

// checkFileExists checks if a file exists in the repository
func (s *WorkflowService) checkFileExists(ctx context.Context, token, owner, repo, path string) (bool, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, "", nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response to get SHA
	var fileInfo struct {
		SHA string `json:"sha"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &fileInfo); err != nil {
		return true, "", err
	}

	return true, fileInfo.SHA, nil
}

// GetWorkflows retrieves all workflow files from a repository
func (s *WorkflowService) GetWorkflows(ctx context.Context, token, owner, repo string) ([]WorkflowFile, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/.github/workflows", owner, repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logger.Error().
			Err(err).
			Str("url", url).
			Msg("Failed to fetch workflows")
		return nil, fmt.Errorf("%w: %v", ErrGitHubAPIFailed, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		// .github/workflows directory doesn't exist
		return []WorkflowFile{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error().
			Int("status", resp.StatusCode).
			Str("body", string(body)).
			Msg("GitHub API returned error when fetching workflows")
		return nil, fmt.Errorf("%w: status %d", ErrGitHubAPIFailed, resp.StatusCode)
	}

	// Parse response
	var files []struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		SHA         string `json:"sha"`
		Size        int    `json:"size"`
		URL         string `json:"url"`
		HTMLURL     string `json:"html_url"`
		GitURL      string `json:"git_url"`
		DownloadURL string `json:"download_url"`
		Type        string `json:"type"`
	}

	if err := json.Unmarshal(body, &files); err != nil {
		logger.Error().
			Err(err).
			Str("body", string(body)).
			Msg("Failed to parse workflows response")
		return nil, err
	}

	// Filter only .yml and .yaml files
	workflows := make([]WorkflowFile, 0)
	for _, file := range files {
		if file.Type == "file" && (strings.HasSuffix(file.Name, ".yml") || strings.HasSuffix(file.Name, ".yaml")) {
			workflows = append(workflows, WorkflowFile{
				Name:        file.Name,
				Path:        file.Path,
				SHA:         file.SHA,
				Size:        file.Size,
				URL:         file.HTMLURL,
				DownloadURL: file.DownloadURL,
			})
		}
	}

	logger.Info().
		Str("owner", owner).
		Str("repo", repo).
		Int("count", len(workflows)).
		Msg("Successfully fetched workflows")

	return workflows, nil
}

// PreviewWorkflowYAML generates and returns the workflow YAML without creating it
func (s *WorkflowService) PreviewWorkflowYAML(request *WorkflowRequest) (string, error) {
	return s.GenerateWorkflowYAML(request)
}

// Helper functions

// isValidWorkflowName validates the workflow name format
func isValidWorkflowName(name string) bool {
	// Workflow name should only contain alphanumeric characters, hyphens, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, name)
	return matched && len(name) > 0 && len(name) <= 255
}

// encodeBase64 encodes a string to base64
func encodeBase64(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}
