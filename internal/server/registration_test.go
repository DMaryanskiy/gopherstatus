package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandleRegistration_JSONSuccess(t *testing.T) {
	db := &MockDB{}
	s := &Server{db: db}

	body := `{
		"email": "test@example.com",
		"password": "securepass",
		"name": "Test User",
		"telegram_username": "gopher"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.handleRegistration(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rec.Code)
	}

	if db.CreatedUser == nil || db.CreatedUser.Email != "test@example.com" {
		t.Errorf("expected user to be created, got %+v", db.CreatedUser)
	}
}

func TestHandleRegistration_FormSuccess(t *testing.T) {
	db := &MockDB{}
	s := &Server{db: db}

	form := url.Values{}
	form.Set("email", "form@example.com")
	form.Set("password", "formpass")
	form.Set("name", "Form User")
	form.Set("telegram_username", "formbot")

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rec := httptest.NewRecorder()
	s.handleRegistration(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 See Other, got %d", rec.Code)
	}
}

func TestHandleRegistration_MissingFields(t *testing.T) {
	db := &MockDB{}
	s := &Server{db: db}

	body := `{"email": "no-name@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.handleRegistration(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}
