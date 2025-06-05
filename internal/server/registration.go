package server

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/DMaryanskiy/gopherstatus/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		Name             string `json:"name"`
		TelegramUsername string `json:"telegram_username,omitempty"`
	}

	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" ||
		strings.HasPrefix(contentType, "application/json") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else if contentType == "application/x-www-form-urlencoded" ||
		strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		req.Name = r.FormValue("name")
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
		req.TelegramUsername = r.FormValue("telegram_username")
	} else {
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	if req.Email == "" || req.Name == "" || req.Password == "" {
		http.Error(w, "Missing field", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user := &storage.User{
		Email:            req.Email,
		PasswordHash:     string(hashed),
		Name:             req.Name,
		TelegramUsername: req.TelegramUsername,
	}

	if err := s.db.CreateUser(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	if strings.HasPrefix(contentType, "application/json") {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"User registered successfully"}`))
	} else {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (s *Server) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("web", "templates", "register.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "failed to load template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}
