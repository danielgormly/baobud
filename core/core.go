package core

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/consul-template/template"
)

func DefaultTemplate(content string) (*template.Template, error) {
	tmplInput := &template.NewTemplateInput{
		Destination:   "",
		Contents:      content,
		LeftDelim:     "{{",
		RightDelim:    "}}",
		ErrMissingKey: true,
		ErrFatal:      true,
	}
	tmpl, err := template.NewTemplate(tmplInput)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// Analyze performs static analysis on template content
func Analyze(content string) ([]string, error) {
	// Parse template using consul-template's parser
	tmpl, err := DefaultTemplate(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error constructing template: %v", err)
		os.Exit(1)
	}

	executeInput := template.ExecuteInput{
		Brain: template.NewBrain(),
	}
	executed, err := tmpl.Execute(&executeInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v", err)
		os.Exit(1)
	}

	vaultPaths := make(map[string]bool)
	for _, dep := range executed.Used.List() {
		if dep.Type() == 1 {
			if matched := regexp.MustCompile(`(?:vault\.read\()([^)]+)(?:\))`).FindStringSubmatch(dep.String()); matched != nil {
				vaultPaths[matched[1]] = true
			}
		}
	}
	paths := make([]string, 0, len(vaultPaths))
	for path := range vaultPaths {
		paths = append(paths, path)
	}
	return paths, nil
}

func CreateVaultPolicy(paths []string) string {
	var policyBuilder strings.Builder

	for _, path := range paths {
		// Add the path with read capabilities
		policyBuilder.WriteString(fmt.Sprintf(`path "%s" {
    capabilities = ["read"]
}
`, path))
	}
	return policyBuilder.String()
}
