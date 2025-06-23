package web

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates/page/*.tmpl templates/layout/*.tmpl
var tmplFS embed.FS

//go:embed static/**/*
var staticFS embed.FS

//go:embed i18n/langs/**/*.yaml
var I18nFS embed.FS

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(tmplFS, "templates/page/*.tmpl", "templates/layout/*.tmpl")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	log.Println("Successfully parsed all embedded template files.")
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
