package template

// EC2Generator generates EC2 deployment workflow templates
type EC2Generator struct {
	*BaseGenerator
}

// NewEC2Generator creates a new EC2 template generator
func NewEC2Generator() *EC2Generator {
	return &EC2Generator{
		BaseGenerator: NewBaseGenerator(),
	}
}

// Generate generates an EC2 workflow YAML
func (eg *EC2Generator) Generate(data interface{}) (string, error) {
	tmpl := `name: {{.WorkflowName}}

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
        include:
{{range .Projects}}          - project: {{.Name}}
            docker_context_path: {{.DockerContextPath}}
            dockerfile_path: {{.DockerfilePath}}
            dot_env_file_testing: |
{{indent 14 .DotEnvTesting}}
            dot_env_file_production: |
{{indent 14 .DotEnvProduction}}
{{end}}    permissions:
      contents: read
      packages: write
    secrets:
      IMAGE_REGISTRY_PASSWORD: {{"{{"}} secrets.IMAGE_REGISTRY_PASSWORD {{"}}"}}

    uses: Calance-US/calance-workflows/.github/workflows/{{.WorkflowFileName}}@{{.EC2CommonFields.ReleaseTag}}
    with:
      image_name: {{.Owner}}/{{.Repository}}-{{"{{"}} matrix.project {{"}}"}}
      image_registry: {{"{{"}} vars.IMAGE_REGISTRY {{"}}"}}
      image_registry_username: {{"{{"}} vars.IMAGE_REGISTRY_USERNAME {{"}}"}}
      docker_context_path: {{"{{"}} matrix.docker_context_path {{"}}"}}
      dockerfile_path: {{"{{"}} matrix.dockerfile_path {{"}}"}}
      dot_env_file_testing: {{"{{"}} matrix.dot_env_file_testing {{"}}"}}
      dot_env_file_production: {{"{{"}} matrix.dot_env_file_production {{"}}"}}

  deploy-to-ec2:
    needs: build-and-push-dockerimages
    strategy:
      fail-fast: false
      matrix:
        include:
{{range .EC2Projects}}          - project: {{.Name}}
{{if .Command}}            command: {{.Command}}
{{end}}            port: {{.Port}}
{{if .DockerNetwork}}            docker_network: {{.DockerNetwork}}
{{end}}{{if .MountPath}}            mount_path: {{.MountPath}}
{{end}}{{if .EnableGPU}}            enable_gpu: true
{{else}}            enable_gpu: false
{{end}}{{if .LogDriver}}            log_driver: {{.LogDriver}}
{{end}}{{if .LogDriverOptions}}            log_driver_options: {{.LogDriverOptions}}
{{end}}{{end}}    permissions:
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
      devops_stakeholders_email_ids: {{default " " .EC2CommonFields.DevopsStakeholdersEmails}}
      command: {{"{{"}} matrix.command | default "" {{"}}"}}
      port: {{"{{"}} matrix.port {{"}}"}}
      docker_network: {{"{{"}} matrix.docker_network | default "" {{"}}"}}
      mount_path: {{"{{"}} matrix.mount_path | default "" {{"}}"}}
      enable_gpu: {{"{{"}} matrix.enable_gpu | default false {{"}}"}}
      log_driver: {{"{{"}} matrix.log_driver | default "" {{"}}"}}
      log_driver_options: {{"{{"}} matrix.log_driver_options | default "" {{"}}"}}
    secrets:
      JENKINS_URL: {{"{{"}} secrets.JENKINS_URL {{"}}"}}
      JENKINS_USER: {{"{{"}} secrets.JENKINS_USER {{"}}"}}
      JENKINS_TOKEN: {{"{{"}} secrets.JENKINS_TOKEN {{"}}"}}
      SMTP_PASSWORD: {{"{{"}} secrets.SMTP_PASSWORD {{"}}"}}
      AWS_CREDENTIALS: {{"{{"}} secrets.AWS_CREDENTIALS {{"}}"}}
`

	return eg.Execute("ec2-workflow", tmpl, data)
}
