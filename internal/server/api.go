package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleAPIStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	results, err := s.db.LatestResults(100)
	if err != nil {
		http.Error(w, "db error: %v", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}

func (s *Server) handleHistoryForService(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	if service == "" {
		http.Error(w, "service is missing in query", http.StatusBadRequest)
		return
	}
	results, err := s.db.HistoryForService(50, service)
	if err != nil {
		http.Error(w, "db error: %v", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}
