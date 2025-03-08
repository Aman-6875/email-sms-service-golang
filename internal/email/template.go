package email

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"
)

func ParseTemplate(templateName string, data map[string]interface{}) (string, error) {
	templatePath := filepath.Join("templates", templateName)

	tmpl, err := template.New(filepath.Base(templatePath)).ParseFiles(templatePath)

	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var buf bytes.Buffer 
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	return buf.String(), nil
}