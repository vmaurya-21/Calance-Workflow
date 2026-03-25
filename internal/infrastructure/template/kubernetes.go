package template

// KubernetesGenerator generates Kubernetes deployment workflow templates
type KubernetesGenerator struct {
	*BaseGenerator
}

// NewKubernetesGenerator creates a new Kubernetes template generator
func NewKubernetesGenerator() *KubernetesGenerator {
	return &KubernetesGenerator{
		BaseGenerator: NewBaseGenerator(),
	}
}

// Generate generates a Kubernetes workflow YAML
func (kg *KubernetesGenerator) Generate(data interface{}) (string, error) {
	tmpl := `{{- $hasTesting := false -}}
{{- $hasProduction := false -}}
{{- range .Projects -}}
  {{- if .DotEnvTesting -}}{{- $hasTesting = true -}}{{- end -}}
  {{- if .DotEnvProduction -}}{{- $hasProduction = true -}}{{- end -}}
{{- end -}}
name: {{.WorkflowName}}

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
        project: [{{range $i, $p := .Projects}}{{if $i}},{{end}}{{$p.Name}}{{end}}]
        include:
{{range .Projects}}          - project: {{.Name}}
            docker_context_path: {{.DockerContextPath}}
            dockerfile_path: {{.DockerfilePath}}
{{if .DotEnvTesting}}            dot_env_file_testing: |
{{indent 14 .DotEnvTesting}}
{{end}}{{if .DotEnvProduction}}            dot_env_file_production: |
{{indent 14 .DotEnvProduction}}
{{end}}
{{end}}    permissions:
      contents: read
      packages: write
    secrets:
      IMAGE_REGISTRY_PASSWORD: ${{"{{"}} secrets.IMAGE_REGISTRY_PASSWORD {{"}}"}}

    uses: Calance-US/calance-workflows/.github/workflows/{{.WorkflowFileName}}@{{.KubernetesCommonFields.ReleaseTag}}
    with:
      image_name: {{.Owner}}/{{.Repository}}-${{"{{"}} matrix.project {{"}}"}}
      image_registry: ${{"{{"}} vars.IMAGE_REGISTRY {{"}}"}}
      image_registry_username: ${{"{{"}} vars.IMAGE_REGISTRY_USERNAME {{"}}"}}
      docker_context_path: ${{"{{"}} matrix.docker_context_path {{"}}"}}
      dockerfile_path: ${{"{{"}} matrix.dockerfile_path {{"}}"}}
{{if $hasTesting}}      dot_env_file_testing: ${{"{{"}} matrix.dot_env_file_testing {{"}}"}}
{{end}}{{if $hasProduction}}      dot_env_file_production: ${{"{{"}} matrix.dot_env_file_production {{"}}"}}
{{end}}

  deploy-to-kubernetes:
    needs: build-and-push-dockerimages
    strategy:
      fail-fast: false
      matrix:
        project: [{{range $i, $p := .KubernetesProjects}}{{if $i}},{{end}}{{$p.Name}}{{end}}]
        include:
{{range .KubernetesProjects}}          - project: {{.Name}}
{{end}}    permissions:
      contents: read
      packages: write

    uses: Calance-US/calance-workflows/.github/workflows/deploy.yml@{{.KubernetesCommonFields.ReleaseTag}}
    with:
      repository_name: ${{"{{"}} github.event.repository.name {{"}}"}}
      image_name: {{.Owner}}/{{.Repository}}-${{"{{"}} matrix.project {{"}}"}}
      release_name: {{.Repository}}-${{"{{"}} matrix.project {{"}}"}}
      image_registry: ${{"{{"}} vars.IMAGE_REGISTRY {{"}}"}}
      version: ${{"{{"}} needs.build-and-push-dockerimages.outputs.version {{"}}"}}
      cluster_environment: ${{"{{"}} needs.build-and-push-dockerimages.outputs.cluster_environment {{"}}"}}
      commit_id: ${{"{{"}} needs.build-and-push-dockerimages.outputs.commit_id {{"}}"}}
      jenkins_job_name: {{.KubernetesCommonFields.JenkinsJobName}}
      workflows_release: {{.KubernetesCommonFields.ReleaseTag}}
      helm_values_repository: {{.KubernetesCommonFields.HelmValuesRepository}}
      codeowners_email_ids: {{.KubernetesCommonFields.CodeownersEmailIds}}
      devops_stakeholders_email_ids: {{.KubernetesCommonFields.DevopsStakeholdersEmailIds}}
    secrets:
      JENKINS_URL: ${{"{{"}} secrets.JENKINS_URL {{"}}"}}
      JENKINS_USER: ${{"{{"}} secrets.JENKINS_USER {{"}}"}}
      JENKINS_TOKEN: ${{"{{"}} secrets.JENKINS_TOKEN {{"}}"}}
      SMTP_PASSWORD: ${{"{{"}} secrets.SMTP_PASSWORD {{"}}"}}
`

	return kg.Execute("k8s-workflow", tmpl, data)
}
