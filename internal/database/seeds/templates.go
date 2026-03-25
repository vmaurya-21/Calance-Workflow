package seeds

import (
	"encoding/json"

	"github.com/vmaurya-21/Calance-Workflow/internal/domain/workflow"
	"gorm.io/gorm"
)

// SeedTemplates populates the database with initial workflow templates
func SeedTemplates(db *gorm.DB) error {
	templateSchema := `{
  "templateId": "build-publish-v3",
  "name": "Build & Publish Image",
  "version": "3.0.1",
  "description": "Standard build and publish workflow for EC2 and Kubernetes",
  "steps": [
    {
      "id": "step1",
      "title": "Basic Information",
      "sections": [
        {
          "id": "base_info",
          "layout": "grid-2",
          "fields": [
            {
              "id": "workflowName",
              "label": "Workflow Name",
              "type": "text",
              "required": true,
              "placeholder": "e.g., Production Deployment",
              "helpText": "Display name for this workflow in the dashboard",
              "defaultValue": "Build & Publish Image"
            },
            {
              "id": "workflowFileName",
              "label": "File Name",
              "type": "text",
              "required": true,
              "placeholder": "build.yml",
              "suffix": ".yml",
              "helpText": "Actual YAML filename in .github/workflows/"
            }
          ]
        },
        {
          "id": "project_config",
          "title": "Project Configuration",
          "isList": true,
          "maxItems": 3,
          "itemLabel": "Project",
          "fields": [
            {
              "id": "name",
              "label": "Project Name",
              "type": "text",
              "required": true,
              "placeholder": "e.g., API Gateway"
            },
            {
              "id": "dockerContextPath",
              "label": "Context Path",
              "type": "text",
              "defaultValue": "./"
            },
            {
              "id": "dockerfilePath",
              "label": "Dockerfile Path",
              "type": "text",
              "defaultValue": "./Dockerfile"
            },
            {
              "id": "dotEnvTesting",
              "label": "Testing Configuration",
              "type": "textarea",
              "placeholder": "PORT=3000"
            },
            {
              "id": "dotEnvProduction",
              "label": "Production Configuration",
              "type": "textarea",
              "placeholder": "PORT=80"
            }
          ]
        }
      ]
    },
    {
      "id": "step2",
      "title": "Deployment Configuration",
      "sections": [
        {
          "id": "deployment_selection",
          "fields": [
            {
              "id": "deploymentType",
              "label": "Deployment Type",
              "type": "select",
              "required": true,
              "options": [
                { "label": "Deploy on EC2", "value": "ec2" },
                { "label": "Deploy on Kubernetes", "value": "kubernetes" }
              ]
            }
          ]
        },
        {
          "id": "dynamic_fields",
          "conditionalSections": [
            {
              "when": { "field": "deploymentType", "equals": "ec2" },
              "sections": [
                {
                  "id": "ec2_common",
                  "title": "Common Configuration",
                  "layout": "grid-2",
                  "fields": [
                    { "id": "credentialId", "label": "Credential ID", "type": "text", "required": true },
                    { "id": "awsRegion", "label": "AWS Region", "type": "text", "required": true },
                    { "id": "jenkinsJobs", "label": "Jenkins Jobs", "type": "text", "required": true },
                    { "id": "releaseTag", "label": "Release Tag", "type": "text", "required": true },
                    { "id": "codeownersEmails", "label": "Codeowners Emails", "type": "text", "required": true },
                    { "id": "devopsStakeholdersEmails", "label": "DevOps Stakeholders Emails", "type": "text", "required": true }
                  ]
                },
                {
                  "id": "ec2_projects",
                  "title": "EC2 Projects Configuration",
                  "isList": true,
                  "linkedToList": "project_config",
                  "fields": [
                    { "id": "name", "label": "Project Name", "type": "text", "disabled": true, "bindTo": "project_config.name" },
                    { "id": "command", "label": "Command (optional)", "type": "text" },
                    { "id": "port", "label": "Port", "type": "text", "required": true },
                    { "id": "dockerNetwork", "label": "Docker Network", "type": "text" },
                    { "id": "mountPath", "label": "Mount Path", "type": "text" },
                    { "id": "enableGpu", "label": "Enable GPU Support", "type": "boolean", "defaultValue": false },
                    { "id": "logDriver", "label": "Log Driver", "type": "select", "options": [{ "label": "json-file", "value": "json-file" }, { "label": "awslogs", "value": "awslogs" }] },
                    { "id": "logDriverOptions", "label": "Log Driver Options", "type": "textarea" }
                  ]
                }
              ]
            },
            {
              "when": { "field": "deploymentType", "equals": "kubernetes" },
              "sections": [
                {
                  "id": "kubernetes_common",
                  "title": "Common Configuration",
                  "layout": "grid-2",
                  "fields": [
                    { "id": "jenkinsJobName", "label": "Jenkins Job Name", "type": "text", "required": true },
                    { "id": "releaseTag", "label": "Release Tag", "type": "text", "required": true },
                    { "id": "helmValuesRepository", "label": "Helm Values Repository", "type": "text", "required": true, "defaultValue": "calance-services-helm-values" },
                    { "id": "codeownersEmailIds", "label": "Codeowners Emails", "type": "text", "required": true },
                    { "id": "devopsStakeholdersEmailIds", "label": "DevOps Stakeholders Emails", "type": "text", "required": true }
                  ]
                },
                {
                  "id": "kubernetes_projects",
                  "title": "Kubernetes Projects Configuration",
                  "isList": true,
                  "linkedToList": "project_config",
                  "fields": [
                    { "id": "name", "label": "Project Name (K8s Service Name)", "type": "text", "required": true, "bindTo": "project_config.name" }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}`

	var initialTemplate workflow.WorkflowTemplate
	err := db.Where("template_id = ?", "build-publish-v3").First(&initialTemplate).Error

	if err == gorm.ErrRecordNotFound {
		initialTemplate = workflow.WorkflowTemplate{
			TemplateID:  "build-publish-v3",
			Name:        "Build & Publish Image",
			Version:     "3.0.1",
			Description: "Standard build and publish workflow for EC2 and Kubernetes",
			Schema:      json.RawMessage(templateSchema),
		}
		return db.Create(&initialTemplate).Error
	}

	if err != nil {
		return err
	}

	// Always update the schema to the latest provided in code
	initialTemplate.Schema = json.RawMessage(templateSchema)
	initialTemplate.Version = "3.0.1"
	return db.Save(&initialTemplate).Error
}
