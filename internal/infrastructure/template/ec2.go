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

	return eg.Execute("ec2-workflow", tmpl, data)
}
