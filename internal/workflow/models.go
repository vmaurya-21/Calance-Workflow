package workflow

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeploymentType represents the type of deployment (EC2 or Kubernetes)
type DeploymentType string

const (
	DeploymentTypeEC2        DeploymentType = "ec2"
	DeploymentTypeKubernetes DeploymentType = "kubernetes"
)

// Project represents common project configuration
type Project struct {
	ID                string `json:"id" binding:"required"`
	Name              string `json:"name" binding:"required"`
	DockerContextPath string `json:"dockerContextPath" binding:"required"`
	DockerfilePath    string `json:"dockerfilePath" binding:"required"`
	DotEnvTesting     string `json:"dotEnvTesting"`
	DotEnvProduction  string `json:"dotEnvProduction"`
}

// EC2CommonFields represents shared EC2 configuration
type EC2CommonFields struct {
	CredentialID             string `json:"credentialId" binding:"required"`
	AWSRegion                string `json:"awsRegion" binding:"required"`
	JenkinsJobs              string `json:"jenkinsJobs" binding:"required"`
	ReleaseTag               string `json:"releaseTag" binding:"required"`
	CodeownersEmails         string `json:"codeownersEmails" binding:"required"`
	DevopsStakeholdersEmails string `json:"devopsStakeholdersEmails" binding:"required"`
}

// EC2Project represents EC2-specific project configuration
type EC2Project struct {
	ID               string `json:"id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	Command          string `json:"command" binding:"required"`
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
	DevopsStakeholdersEmailIds string `json:"devopsStakeholdersEmailIds" binding:"required"`
}

// KubernetesProject represents Kubernetes-specific project configuration
type KubernetesProject struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// WorkflowRequest represents the request to create a workflow
type WorkflowRequest struct {
	Owner                  string                  `json:"owner" binding:"required"`
	Repository             string                  `json:"repository" binding:"required"`
	WorkflowName           string                  `json:"workflowName" binding:"required"`
	DeploymentType         DeploymentType          `json:"deploymentType" binding:"required,oneof=ec2 kubernetes"`
	Projects               []Project               `json:"projects" binding:"required,min=1,dive"`
	EC2CommonFields        *EC2CommonFields        `json:"ec2CommonFields"`
	EC2Projects            []EC2Project            `json:"ec2Projects"`
	KubernetesCommonFields *KubernetesCommonFields `json:"kubernetesCommonFields"`
	KubernetesProjects     []KubernetesProject     `json:"kubernetesProjects"`
}

// Validate validates the workflow request based on deployment type
func (w *WorkflowRequest) Validate() error {
	if w.DeploymentType == DeploymentTypeEC2 {
		if w.EC2CommonFields == nil {
			return ErrEC2CommonFieldsRequired
		}
		if len(w.EC2Projects) == 0 {
			return ErrEC2ProjectsRequired
		}
	} else if w.DeploymentType == DeploymentTypeKubernetes {
		if w.KubernetesCommonFields == nil {
			return ErrKubernetesCommonFieldsRequired
		}
		if len(w.KubernetesProjects) == 0 {
			return ErrKubernetesProjectsRequired
		}
	}
	return nil
}

// WorkflowResponse represents the response after creating a workflow
type WorkflowResponse struct {
	Owner        string    `json:"owner"`
	Repository   string    `json:"repository"`
	WorkflowName string    `json:"workflowName"`
	FilePath     string    `json:"filePath"`
	FileURL      string    `json:"fileUrl"`
	ContentSHA   string    `json:"contentSha"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"createdAt"`
}

// WorkflowFile represents a workflow file in a repository
type WorkflowFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	DownloadURL string `json:"downloadUrl"`
}

// WorkflowHistory represents a record of workflow creation (optional - for future use)
type WorkflowHistory struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Owner          string         `gorm:"not null" json:"owner"`
	Repository     string         `gorm:"not null" json:"repository"`
	WorkflowName   string         `gorm:"not null" json:"workflow_name"`
	DeploymentType DeploymentType `gorm:"type:varchar(20);not null" json:"deployment_type"`
	FilePath       string         `json:"file_path"`
	ContentSHA     string         `json:"content_sha"`
	Status         string         `gorm:"type:varchar(20);not null" json:"status"` // created, updated, failed
	ErrorMessage   *string        `json:"error_message,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to generate UUID before creating workflow history
func (w *WorkflowHistory) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}
