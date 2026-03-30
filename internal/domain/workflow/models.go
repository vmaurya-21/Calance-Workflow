package workflow

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeploymentType represents the type of deployment
type DeploymentType string

const (
	DeploymentTypeEC2        DeploymentType = "ec2"
	DeploymentTypeKubernetes DeploymentType = "kubernetes"
)

// Request represents a workflow creation request
type Request struct {
	Owner                     string                     `json:"owner" binding:"required"`
	Repository                string                     `json:"repository" binding:"required"`
	WorkflowName              string                     `json:"workflowName" binding:"required"`
	WorkflowFileName          string                     `json:"workflowFileName" binding:"required"`
	DeploymentType            DeploymentType             `json:"deploymentType" binding:"required,oneof=ec2 kubernetes"`
	Projects                  []Project                  `json:"projects" binding:"required,min=1,dive"`
	EC2CommonFields           *EC2CommonFields           `json:"ec2CommonFields"`
	EC2Projects               []EC2Project               `json:"ec2Projects"`
	InfrastructureCredentials []InfrastructureCredential `json:"infrastructureCredentials"`
	KubernetesCommonFields    *KubernetesCommonFields    `json:"kubernetesCommonFields"`
	KubernetesProjects        []KubernetesProject        `json:"kubernetesProjects"`
}

// Validate validates the workflow request
func (r *Request) Validate() error {
	if r.DeploymentType == DeploymentTypeEC2 {
		if r.EC2CommonFields == nil {
			return ErrEC2CommonFieldsRequired
		}
		if len(r.EC2Projects) == 0 {
			return ErrEC2ProjectsRequired
		}
	} else if r.DeploymentType == DeploymentTypeKubernetes {
		if r.KubernetesCommonFields == nil {
			return ErrKubernetesCommonFieldsRequired
		}
		if len(r.KubernetesProjects) == 0 {
			return ErrKubernetesProjectsRequired
		}
	}
	return nil
}

// Project represents common project configuration
type Project struct {
	ID                string `json:"id"`
	Name              string `json:"name" binding:"required"`
	DockerContextPath string `json:"dockerContextPath"`
	DockerfilePath    string `json:"dockerfilePath"`
	DotEnvTesting     string `json:"dotEnvTesting"`
	DotEnvProduction  string `json:"dotEnvProduction"`
}

// InfrastructureCredential represents EC2 and Kubernetes infrastructure credentials
type InfrastructureCredential struct {
	JenkinsNodeName    string `json:"jenkinsNodeName"`
	ProjectName        string `json:"projectName"`
	Environment        string `json:"environment"`
	AWSAccessKeyID     string `json:"awsAccessKeyId"`
	AWSSecretAccessKey string `json:"awsSecretAccessKey"`
	Environments       string `json:"environments"`
	Namespace          string `json:"namespace"` // Used for Kubernetes
}

// EC2CommonFields represents shared EC2 configuration
type EC2CommonFields struct {
	CredentialID             string `json:"credentialId" binding:"required"`
	AWSRegion                string `json:"awsRegion" binding:"required"`
	JenkinsJobs              string `json:"jenkinsJobs" binding:"required"`
	ReleaseTag               string `json:"releaseTag" binding:"required"`
	CodeownersEmails         string `json:"codeownersEmails" binding:"required"`
	DevopsStakeholdersEmails string `json:"devopsStakeholdersEmails"`
}

// EC2Project represents EC2-specific project configuration
type EC2Project struct {
	ID               string `json:"id"`
	Name             string `json:"name" binding:"required"`
	Command          string `json:"command"`
	Port             string `json:"port" binding:"required"`
	DockerNetwork    string `json:"dockerNetwork"`
	MountPath        string `json:"mountPath"`
	EnableGPU        bool   `json:"enableGpu"`
	LogDriver        string `json:"logDriver"`
	LogDriverOptions string `json:"logDriverOptions"`
}

// KubernetesCommonFields represents shared Kubernetes configuration
type KubernetesCommonFields struct {
	JenkinsJobName             string `json:"jenkinsJobName" binding:"required"`
	ReleaseTag                 string `json:"releaseTag" binding:"required"`
	HelmValuesRepository       string `json:"helmValuesRepository" binding:"required"`
	CodeownersEmailIds         string `json:"codeownersEmailIds" binding:"required"`
	DevopsStakeholdersEmailIds string `json:"devopsStakeholdersEmailIds"`
}

// KubernetesProject represents Kubernetes-specific project configuration
type KubernetesProject struct {
	ID   string `json:"id"`
	Name string `json:"name" binding:"required"`
}

// Response represents a workflow creation response
type Response struct {
	Owner        string    `json:"owner"`
	Repository   string    `json:"repository"`
	WorkflowName string    `json:"workflowName"`
	FilePath     string    `json:"filePath"`
	FileURL      string    `json:"fileUrl"`
	ContentSHA   string    `json:"contentSha"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"createdAt"`
}

// PRResponse represents a pull request status response
type PRResponse struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	URL    string `json:"url"`
	Branch string `json:"branch"`
}

// File represents a workflow file
type File struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	DownloadURL string `json:"downloadUrl"`
}

// History represents workflow creation history
type History struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Owner          string         `gorm:"not null" json:"owner"`
	Repository     string         `gorm:"not null" json:"repository"`
	WorkflowName   string         `gorm:"not null" json:"workflow_name"`
	DeploymentType DeploymentType `gorm:"type:varchar(20);not null" json:"deployment_type"`
	FilePath       string         `json:"file_path"`
	ContentSHA     string         `json:"content_sha"`
	Status         string         `gorm:"type:varchar(20);not null" json:"status"`
	PRNumber       int            `json:"pr_number"`
	PRURL          string         `json:"pr_url"`
	ErrorMessage   *string        `json:"error_message,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook
func (h *History) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// FileContentResponse represents workflow file content
type FileContentResponse struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	SHA     string `json:"sha"`
	Size    int    `json:"size"`
	Content string `json:"content"`
}

// UpdateWorkflowRequest represents workflow update request
type UpdateWorkflowRequest struct {
	Owner         string `json:"owner" binding:"required"`
	Repository    string `json:"repository" binding:"required"`
	FilePath      string `json:"filePath" binding:"required"`
	Content       string `json:"content" binding:"required"`
	SHA           string `json:"sha" binding:"required"`
	CommitMessage string `json:"commitMessage"`
}

// WorkflowTemplate represents a dynamic form schema for workflows
type WorkflowTemplate struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	TemplateID  string          `gorm:"type:varchar(100);uniqueIndex;not null" json:"templateId"`
	Name        string          `gorm:"type:varchar(255);not null" json:"name"`
	Version     string          `gorm:"type:varchar(50);not null" json:"version"`
	Description string          `gorm:"type:text" json:"description"`
	Schema      json.RawMessage `gorm:"type:jsonb" json:"schema"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// TemplateSummary represents a concise view of a template
type TemplateSummary struct {
	TemplateID  string `json:"templateId"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// BeforeCreate hook for WorkflowTemplate
func (t *WorkflowTemplate) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
