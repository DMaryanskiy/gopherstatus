package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DMaryanskiy/gopherstatus/internal/auth"
	"github.com/DMaryanskiy/gopherstatus/internal/storage"
)

func (s *Server) handleCreateService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := auth.GetUserIDFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to get current user", http.StatusInternalServerError)
		return
	}

	var req struct {
		Name     string            `json:"name"`
		URL      string            `json:"url"`
		Method   string            `json:"method"`
		Interval int               `json:"interval"`
		Body     string            `json:"body"`
		Headers  map[string]string `json:"headers"`
	}

	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" ||
		strings.HasPrefix(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
	} else if contentType == "application/x-www-form-urlencoded" ||
		strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		req.Name = r.FormValue("name")
		req.URL = r.FormValue("url")
		req.Method = r.FormValue("method")

		intervalStr := r.FormValue("interval")
		interval, err := strconv.ParseInt(intervalStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid interval", http.StatusBadRequest)
			return
		}
		req.Interval = int(interval)
		req.Body = r.FormValue("body")

		headers := make(map[string]string)

		headersPrefix := "headers["
		for key, values := range r.Form {
			if strings.HasPrefix(key, headersPrefix) && strings.HasSuffix(key, "]") {
				headerName := key[len(headersPrefix) : len(key)-1]
				headers[headerName] = values[0]
			}
		}
		req.Headers = headers
	} else {
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	service := &storage.Service{
		Name:     req.Name,
		URL:      req.URL,
		Method:   req.Method,
		Interval: req.Interval,
		Body:     req.Body,
		UserId:   userID,
	}
	for k, v := range req.Headers {
		service.Headers = append(service.Headers, storage.Header{
			Key:   k,
			Value: v,
		})
	}

	if err := s.db.CreateService(service); err != nil {
		http.Error(w, "Failed to create service", http.StatusInternalServerError)
		return
	}

	if strings.HasPrefix(contentType, "application/json") {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"Service created successfully"}`))
	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (s *Server) handleServiceCreationForm(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("web", "templates", "service_create.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}
