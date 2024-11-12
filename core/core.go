package core

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/consul-template/dependency"
	"github.com/hashicorp/consul-template/logging"
	"github.com/hashicorp/consul-template/template"
)

type BaobudConfig struct {
	BaoAddress string
	BaoToken   string
}

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
func Analyze(content string, config BaobudConfig) ([]string, error) {
	logging.Setup(&logging.Config{
		Level: "WARN",
	})
	clients := dependency.NewClientSet()
	clients.CreateVaultClient((&dependency.CreateVaultClientInput{
		Address: config.BaoAddress,
		Token:   config.BaoToken,
	}))
	opts := &dependency.QueryOptions{
		// Set any required options
		WaitIndex:  0,
		WaitTime:   time.Second * 10,
		AllowStale: false,
	}

	brain := template.NewBrain()
	executeInput := template.ExecuteInput{
		Brain: brain,
	}
	vaultPaths := make(map[string]bool)

	// Loop through execution of template after each
	// set of Vault paths are discovered
	for {
		// Parse template using consul-template's parser
		tmpl, err := DefaultTemplate(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error constructing template: %v", err)
			os.Exit(1)
		}
		executed, err := tmpl.Execute(&executeInput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing template: %v", err)
			os.Exit(1)
		}

		for _, dep := range executed.Missing.List() {
			if dep.Type() == 1 {
				result, _, err := dep.Fetch(clients, opts)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching: %v\n", err)
					os.Exit(1)
				}
				brain.Remember(dep, result)
				if matched := regexp.MustCompile(`(?:vault\.read\()([^)]+)(?:\))`).FindStringSubmatch(dep.String()); matched != nil {
					vaultPaths[matched[1]] = true
				}
			}
		}
		if executed.Missing.Len() == 0 {
			break
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
