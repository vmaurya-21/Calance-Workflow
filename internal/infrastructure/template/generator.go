package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Generator defines the interface for workflow template generation
type Generator interface {
	Generate(data interface{}) (string, error)
}

// BaseGenerator provides common template functionality
type BaseGenerator struct {
	funcMap template.FuncMap
}

// NewBaseGenerator creates a new base generator with common functions
func NewBaseGenerator() *BaseGenerator {
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

	return &BaseGenerator{
		funcMap: funcMap,
	}
}

// Execute executes a template with the given data
func (bg *BaseGenerator) Execute(name, tmpl string, data interface{}) (string, error) {
	t, err := template.New(name).Funcs(bg.funcMap).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
