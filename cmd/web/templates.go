package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	templates *template.Template
	once      sync.Once
)

// LoadTemplates parses all HTML template files in the static/components/page directory and initializes the templates collection once for the application lifetime.
func LoadTemplates() {
	once.Do(func() {
		baseDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}

		pattern := filepath.Join(baseDir, "static", "components", "page", "*.tmpl")
		templates, err = template.ParseGlob(pattern)
		if err != nil {
			log.Fatalf("Failed to parse templates: %v", err)
		}
	})
}

// RenderTemplate executes the specified HTML template with the provided data and writes the result to the HTTP response.
// If template execution fails, it sends an HTTP 500 Internal Server Error with the error message.
func RenderTemplate(name string, data any, w http.ResponseWriter) {
	LoadTemplates()
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
