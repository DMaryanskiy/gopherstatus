package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DMaryanskiy/gopherstatus/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func TestHandleLogin_Success(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := storage.User{ID: 1, PasswordHash: string(hashed)}
	server := &Server{db: &MockDB{User: mockUser}}

	body, _ := json.Marshal(map[string]string{
		"email":    "valid@example.com",
		"password": "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.handleLogin(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestHandleLogin_InvalidPassword(t *testing.T) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	mockUser := storage.User{ID: 1, PasswordHash: string(hashed)}
	server := &Server{db: &MockDB{User: mockUser}}

	body, _ := json.Marshal(map[string]string{
		"email":    "valid@example.com",
		"password": "wrong-password",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.handleLogin(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}

func TestHandleLogin_UnsupportedContentType(t *testing.T) {
	server := &Server{}

	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.Header.Set("Content-Type", "text/plain")

	rec := httptest.NewRecorder()
	server.handleLogin(rec, req)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415 Unsupported Media Type, got %d", rec.Code)
	}
}
