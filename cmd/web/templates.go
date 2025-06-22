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

func LoadTemplates() {
	once.Do(func() {
		baseDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}

		pattern := filepath.Join(baseDir, "cmd", "web", "templates", "page", "*.tmpl")
		templates, err = template.ParseGlob(pattern)
		if err != nil {
			log.Fatalf("Failed to parse templates: %v", err)
		}
	})
}

func RenderTemplate(name string, data any, w http.ResponseWriter) {
	LoadTemplates()
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
