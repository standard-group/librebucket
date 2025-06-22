package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates *template.Template

func init() {
	// Initialize templates map
	templates = template.New("")
	templateFiles, err := filepath.Glob("cmd/web/templates/page/*.tmpl")
	if err != nil {
		log.Fatalf("Failed to find template files: %v", err)
	}

	layoutFiles, err := filepath.Glob("cmd/web/templates/layout/*.tmpl")
	if err != nil {
		log.Fatalf("Failed to find layout files: %v", err)
	}

	allFiles := append(templateFiles, layoutFiles...)

	if len(allFiles) == 0 {
		log.Println("No template files found.")
		return
	}

	// Parse all template files
	templates, err = template.ParseFiles(allFiles...)
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	log.Println("Successfully parsed all template files.")
}

// RenderTemplate executes the specified HTML template with the provided data and writes the result to the HTTP response.
// If template execution fails, it sends an HTTP 500 Internal Server Error with the error message.
func RenderTemplate(name string, data any, w http.ResponseWriter) {
	// Ensure the template exists
	tmpl := templates.Lookup(name)
	if tmpl == nil {
		http.Error(w, fmt.Sprintf("The template %s does not exist.", name), http.StatusInternalServerError)
		log.Printf("Template not found: %s", name)
		return
	}

	// Execute the template with the provided data
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Error executing template %s: %v", name, err)
	}
}
