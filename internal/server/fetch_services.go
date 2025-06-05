package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/DMaryanskiy/gopherstatus/internal/auth"
)

func (s *Server) handleFetchServicesAPI(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to get current user", http.StatusInternalServerError)
		return
	}

	services, err := s.db.FetchServicesByUser(userID)
	if err != nil {
		http.Error(w, "db error: %v", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(services); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Server) handleFetchServices(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("web", "templates", "service_fetch.html")
	headerPath := filepath.Join("web", "templates", "header.html")

	tmpl, err := template.ParseFiles(tmplPath, headerPath)
	if err != nil {
		http.Error(w, "failed to get template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}
