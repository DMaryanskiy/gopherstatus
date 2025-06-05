package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/DMaryanskiy/gopherstatus/internal/auth"
)

func (s *Server) handleDashboardAPI(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to get current user", http.StatusInternalServerError)
		return
	}

	services, err := s.db.LatestResultsByID(5, userID)
	if err != nil {
		http.Error(w, "db error: %v", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(services); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Server) handleDashboard(w http.ResponseWriter, _ *http.Request) {
	tmplPath := filepath.Join("web", "templates", "dashboard.html")
	headerPath := filepath.Join("web", "templates", "header.html")
	tmpl, err := template.ParseFiles(tmplPath, headerPath)
	if err != nil {
		http.Error(w, "failed to get template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}
