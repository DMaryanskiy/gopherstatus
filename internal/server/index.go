package server

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("web", "templates", "index.html")
	headerPath := filepath.Join("web", "templates", "header.html")
	tmpl, err := template.ParseFiles(tmplPath, headerPath)
	if err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}
