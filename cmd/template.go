package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

// TemplateData holds the variables that can be injected into components
type TemplateData struct {
	ProjectName string
	ModuleName  string
}

// processTemplate processes the given content with Go's text/template engine.
// If the content is not a valid template, it returns the original content gracefully.
func processTemplate(content []byte) []byte {
	cwd, err := os.Getwd()
	if err != nil {
		return content
	}

	projectName := filepath.Base(cwd)
	data := TemplateData{
		ProjectName: projectName,
		ModuleName:  projectName,
	}

	// Support custom delimiters or just standard {{ }}.
	// To avoid conflicts with JSON/JS objects, users should ideally use {{ .ProjectName }}
	// If parsing fails (e.g., due to invalid Go template syntax from existing code),
	// we fall back to simple string replacement to be safe.

	// Pre-replace legacy curly brace syntax for backward compatibility
	content = bytes.ReplaceAll(content, []byte("{{.ProjectName}}"), []byte(projectName))
	content = bytes.ReplaceAll(content, []byte("{{.ModuleName}}"), []byte(projectName))

	// Try compiling as template using safe delimiters <% %>
	tmpl, err := template.New("file").Delims("<%", "%>").Option("missingkey=zero").Parse(string(content))
	if err != nil {
		// Fallback to simple replace if template logic fails
		content = bytes.ReplaceAll(content, []byte("<%.ProjectName%>"), []byte(projectName))
		content = bytes.ReplaceAll(content, []byte("<%.ModuleName%>"), []byte(projectName))
		return content
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return content // execution failed, return raw
	}
	return buf.Bytes()
}
